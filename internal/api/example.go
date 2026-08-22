package api

import (
	"net/http"
	"os"
	"path/filepath"

	"greenshields/internal/example"
)

// handleExample serves the highway example as JSON. It prefers the on-disk
// copy (so operators can override it) and falls back to the embedded copy.
func (s *Server) handleExample(w http.ResponseWriter, r *http.Request) {
	data := example.FreewayJSON
	if s.ExampleDir != "" {
		if f, err := os.ReadFile(filepath.Join(s.ExampleDir, "freeway.json")); err == nil {
			data = f
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
