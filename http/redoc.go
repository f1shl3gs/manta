package http

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed openapi.yaml
var OpenAPI embed.FS

const indexTemplate = `<!DOCTYPE html>
<html>
  <head>
    <title>Manta API</title>
    <!-- needed for adaptive design -->
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- ReDoc doesn't change outer page styles -->
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
<body>
  <redoc spec-url='%s' suppressWarnings=true></redoc>
  <script src="https://cdn.jsdelivr.net/npm/redoc/bundles/redoc.standalone.js"> </script>
</body>
</html>
`

func Redoc() http.HandlerFunc {
	const specPath = "/docs/openapi.yaml"

	var index = []byte(fmt.Sprintf(indexTemplate, specPath))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == specPath {
			w.Header().Set("Content-Type", "text/yaml")
			w.WriteHeader(http.StatusOK)

			data, err := OpenAPI.ReadFile("openapi.yaml")
			if err != nil {
				panic("read openapi.yaml failed")
			}

			_, _ = w.Write(data)

			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(index)
	}
}
