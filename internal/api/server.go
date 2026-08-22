// Package api exposes the Greenshields model over HTTP. It implements the JSON
// endpoints consumed by the web UI and used by automated checks:
//
//	POST /api/qk      -> speed, flow and congestion side for a density
//	POST /api/curve   -> the q(k) and v(k) curves plus the capacity point
//	POST /api/wave    -> the shock speed between two traffic states
//	GET  /api/example -> the embedded highway example
//	GET  /api/curve.svg -> an SVG of the q(k) parabola (server rendered)
//
// All error responses are JSON objects of the form {"error": "..."} with a
// non-2xx status code, so the web UI can surface backend failures instead of
// silently showing empty tables.
package api

import (
	"encoding/json"
	"net/http"
)

// Server holds the dependencies needed by the HTTP handlers.
type Server struct {
	// ExampleDir is the directory containing freeway.json. When empty, only
	// the embedded example is served.
	ExampleDir string
}

// NewServer builds an API server.
func NewServer(exampleDir string) *Server {
	return &Server{ExampleDir: exampleDir}
}

// Register installs all API routes and, if webDir is non-empty, serves the
// static web assets from it at the root path. API routes take precedence over
// the static handler.
func (s *Server) Register(mux *http.ServeMux, webDir string) {
	mux.HandleFunc("/api/qk", s.handleQK)
	mux.HandleFunc("/api/curve", s.handleCurve)
	mux.HandleFunc("/api/wave", s.handleWave)
	mux.HandleFunc("/api/example", s.handleExample)
	mux.HandleFunc("/api/curve.svg", s.handleCurveSVG)
	if webDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(webDir)))
	}
}

// writeJSON writes a value as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error body with the given status code.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
