package handlers

import (
	"embed"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday/v2"
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
		filePath := findFile("docs/openapi.yaml", []string{
			"./docs/openapi.yaml",
			"../docs/openapi.yaml",
			"../../docs/openapi.yaml",
		})

		if filePath == "" {
			http.Error(w, "OpenAPI spec file not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, filePath)
	}
}

// MarkdownDocHandler sirve archivos Markdown convertidos a HTML
func MarkdownDocHandler(docPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener el directorio de trabajo actual
		wd, err := os.Getwd()
		if err != nil {
			http.Error(w, "Error getting working directory", http.StatusInternalServerError)
			return
		}

		// Buscar el archivo en varias ubicaciones posibles
		// Los archivos están en:
		// - api/docs/API_DOCUMENTATION.md
		// - api/docs/USER_GUIDE.md
		// - docs/GRAPHQL_API_REFERENCE.md (en el workspace root)
		possiblePaths := []string{
			// Desde el directorio api/ (cuando se ejecuta desde api/)
			filepath.Join(wd, "docs", docPath),           // api/docs/API_DOCUMENTATION.md
			// Desde el workspace root (cuando se ejecuta desde workspace/)
			filepath.Join(wd, "api", "docs", docPath),    // api/docs/API_DOCUMENTATION.md
			filepath.Join(wd, "docs", docPath),           // docs/GRAPHQL_API_REFERENCE.md
			// Rutas relativas desde api/
			"./docs/" + docPath,
			"../docs/" + docPath,
			"../../docs/" + docPath,
			"./api/docs/" + docPath,
			"../api/docs/" + docPath,
		}

		// También buscar desde el directorio padre (workspace root)
		parentWd := filepath.Dir(wd)
		if parentWd != wd {
			possiblePaths = append(possiblePaths,
				filepath.Join(parentWd, "docs", docPath),
				filepath.Join(parentWd, "api", "docs", docPath),
			)
		}

		var filePath string
		for _, path := range possiblePaths {
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				filePath = path
				break
			}
		}

		if filePath == "" {
			http.Error(w, "Documentation file not found: "+docPath, http.StatusNotFound)
			return
		}

		// Leer el archivo Markdown
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		// Convertir Markdown a HTML
		html := blackfriday.Run(content)

		// Crear página HTML completa
		htmlPage := createMarkdownHTMLPage(string(html), getTitleFromPath(docPath))

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlPage))
	}
}

// findFile busca un archivo en múltiples rutas posibles
func findFile(filename string, possiblePaths []string) string {
	wd, _ := os.Getwd()
	possiblePaths = append(possiblePaths,
		filepath.Join(wd, filename),
		filepath.Join(wd, "api", filename),
	)

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// getTitleFromPath extrae un título legible del path del archivo
func getTitleFromPath(path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, ".md")
	base = strings.TrimSuffix(base, ".MD")
	
	// Reemplazar guiones y guiones bajos con espacios
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	
	// Capitalizar primera letra de cada palabra
	words := strings.Fields(base)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	
	return strings.Join(words, " ")
}

// createMarkdownHTMLPage crea una página HTML completa con el contenido Markdown
func createMarkdownHTMLPage(content, title string) string {
	return `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>` + template.HTMLEscapeString(title) + `</title>
  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
      line-height: 1.6;
      color: #333;
      background: #f5f5f5;
      padding: 20px;
    }
    .container {
      max-width: 1000px;
      margin: 0 auto;
      background: white;
      padding: 40px;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    }
    .header {
      margin-bottom: 30px;
      padding-bottom: 20px;
      border-bottom: 2px solid #667eea;
    }
    .header h1 {
      color: #667eea;
      font-size: 2.5em;
      margin-bottom: 10px;
    }
    .back-link {
      display: inline-block;
      margin-top: 20px;
      color: #667eea;
      text-decoration: none;
      font-weight: 500;
    }
    .back-link:hover {
      text-decoration: underline;
    }
    .content {
      line-height: 1.8;
    }
    .content h1, .content h2, .content h3, .content h4 {
      color: #333;
      margin-top: 30px;
      margin-bottom: 15px;
    }
    .content h1 {
      font-size: 2em;
      border-bottom: 2px solid #eee;
      padding-bottom: 10px;
    }
    .content h2 {
      font-size: 1.5em;
      color: #667eea;
    }
    .content h3 {
      font-size: 1.3em;
    }
    .content code {
      background: #f4f4f4;
      padding: 2px 6px;
      border-radius: 3px;
      font-family: 'Courier New', monospace;
      font-size: 0.9em;
      color: #d63384;
    }
    .content pre {
      background: #f4f4f4;
      padding: 15px;
      border-radius: 6px;
      overflow-x: auto;
      margin: 20px 0;
    }
    .content pre code {
      background: none;
      padding: 0;
      color: #333;
    }
    .content table {
      width: 100%;
      border-collapse: collapse;
      margin: 20px 0;
    }
    .content table th,
    .content table td {
      border: 1px solid #ddd;
      padding: 12px;
      text-align: left;
    }
    .content table th {
      background: #667eea;
      color: white;
      font-weight: 600;
    }
    .content table tr:nth-child(even) {
      background: #f9f9f9;
    }
    .content ul, .content ol {
      margin: 15px 0;
      padding-left: 30px;
    }
    .content li {
      margin: 8px 0;
    }
    .content blockquote {
      border-left: 4px solid #667eea;
      padding-left: 20px;
      margin: 20px 0;
      color: #666;
      font-style: italic;
    }
    .content a {
      color: #667eea;
      text-decoration: none;
    }
    .content a:hover {
      text-decoration: underline;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>` + template.HTMLEscapeString(title) + `</h1>
      <a href="/docs" class="back-link">← Volver a Documentación</a>
    </div>
    <div class="content">
      ` + string(content) + `
    </div>
    <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #eee;">
      <a href="/docs" class="back-link">← Volver a Documentación</a>
    </div>
  </div>
</body>
</html>`
}
