package api

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	ExampleDir string
}

func NewServer(exampleDir string) *Server {
	return &Server{ExampleDir: exampleDir}
}

func (s *Server) Register(mux *http.ServeMux, webDir string) {
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/qk", s.handleQK)
	mux.HandleFunc("/api/curve", s.handleCurve)
	mux.HandleFunc("/api/wave", s.handleWave)
	mux.HandleFunc("/api/example", s.handleExample)
	mux.HandleFunc("/api/curve.svg", s.handleCurveSVG)
	if webDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(webDir)))
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeUnprocessable(w http.ResponseWriter, err error) {
	writeError(w, http.StatusUnprocessableEntity, err.Error())
}
