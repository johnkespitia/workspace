package handlers

import (
	"encoding/json"
	"net/http"

	gql "github.com/graphql-go/graphql"
)

// GraphQLHandler maneja las peticiones GraphQL
type GraphQLHandler struct {
	schema gql.Schema
}

// NewGraphQLHandler crea un nuevo handler GraphQL
func NewGraphQLHandler(schema gql.Schema) *GraphQLHandler {
	return &GraphQLHandler{
		schema: schema,
	}
}

// ServeHTTP implementa http.Handler
func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Manejar preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Solo aceptar POST para queries GraphQL
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parsear request body
	var req struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Crear un contexto con timeout más largo para operaciones como syncStocks
	ctx := r.Context()
	
	// Ejecutar query
	result := gql.Do(gql.Params{
		Schema:         h.schema,
		RequestString:  req.Query,
		VariableValues: req.Variables,
		Context:       ctx,
	})

	// Verificar si el contexto fue cancelado
	if ctx.Err() != nil {
		// Si el contexto fue cancelado, intentar enviar un error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusRequestTimeout)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []map[string]interface{}{
				{
					"message": "Request timeout: operation took too long",
				},
			},
		})
		return
	}

	// Configurar headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // GraphQL siempre retorna 200, errores van en el body

	// Escribir respuesta
	if err := json.NewEncoder(w).Encode(result); err != nil {
		// Si falla al escribir, ya no podemos hacer mucho
		// Pero al menos intentamos
		return
	}
	
	// Asegurar que la respuesta se envíe
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// PlaygroundHandler maneja el GraphQL Playground (simple HTML)
func PlaygroundHandler(title, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(playgroundHTML(title, endpoint)))
	}
}

func playgroundHTML(title, endpoint string) string {
	return `<!DOCTYPE html>
<html>
<head>
	<title>` + title + `</title>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
	<link rel="shortcut icon" href="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/favicon.png" />
	<script src="https://cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
</head>
<body>
	<div id="root">
		<style>
			body {
				margin: 0;
				height: 100vh;
				overflow: hidden;
			}
			#root {
				width: 100vw;
				height: 100vh;
			}
		</style>
		<script>
			window.addEventListener('load', function (event) {
				GraphQLPlayground.init(document.getElementById('root'), {
					endpoint: '` + endpoint + `'
				})
			})
		</script>
	</div>
</body>
</html>`
}
