package server

import (
	"net/http"
	"path/filepath"
)

func CustomFileServer(dataDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		filePath := filepath.Join(dataDir, r.URL.Path)

		if filepath.Ext(filePath) == ".ndjson" {
			// Set the correct MIME type for NDJSON with fhir data
			w.Header().Set("Content-Type", "application/fhir+ndjson")
		}

		http.ServeFile(w, r, filePath)
	})
}
