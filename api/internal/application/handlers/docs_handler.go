package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed docs_templates/*
var docsTemplates embed.FS

// DocsHandler maneja las peticiones de documentación
func DocsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(docsTemplates, "docs_templates/*.html")
		if err != nil {
			http.Error(w, "Error loading templates", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		data := map[string]interface{}{
			"Title":               "API Documentation",
			"GraphQLPlaygroundURL": "/playground",
			"GraphQLEndpoint":      "/query",
			"HealthEndpoint":       "/health",
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
}

// SwaggerUIHandler sirve Swagger UI para el OpenAPI spec
func SwaggerUIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		swaggerHTML := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>API Documentation - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui.css" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "/docs/openapi.yaml",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
    };
  </script>
</body>
</html>`

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(swaggerHTML))
	}
}

// OpenAPISpecHandler sirve el archivo OpenAPI YAML
func OpenAPISpecHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Intentar encontrar el archivo openapi.yaml
		// Buscar en varias ubicaciones posibles
		possiblePaths := []string{
			"./docs/openapi.yaml",
			"../docs/openapi.yaml",
			"../../docs/openapi.yaml",
		}

		// También buscar desde el directorio de trabajo actual
		wd, _ := os.Getwd()
		possiblePaths = append(possiblePaths,
			filepath.Join(wd, "docs", "openapi.yaml"),
			filepath.Join(wd, "api", "docs", "openapi.yaml"),
		)

		var filePath string
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				filePath = path
				break
			}
		}

		if filePath == "" {
			http.Error(w, "OpenAPI spec file not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, filePath)
	}
}
