package api

import (
	"encoding/json"
	"net/http"

	"greenshields/internal/core"
)

type qkRequest struct {
	Vf float64 `json:"vf"`
	Kj float64 `json:"kj"`
	K  float64 `json:"k"`
}

type qkResponse struct {
	Vf        float64 `json:"vf"`
	Kj        float64 `json:"kj"`
	K         float64 `json:"k"`
	V         float64 `json:"v"`
	Q         float64 `json:"q"`
	Side      string  `json:"side"`
	Congested bool    `json:"congested"`
}

// handleQK computes the speed and flow at a single density and reports which
// side of the capacity point the density lies on.
func (s *Server) handleQK(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}
	var req qkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}
	m, err := core.New(req.Vf, req.Kj)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := m.ValidateDensity(req.K); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rec := newQKRecorder()
	defer rec.Close()

	v, _ := m.Speed(req.K)
	q, _ := m.Flow(req.K)
	_, km := m.Capacity()
	side := "free"
	congested := false
	if req.K > km+core.Epsilon*m.Kj {
		side = "congested"
		congested = true
	}
	writeJSON(w, http.StatusOK, qkResponse{
		Vf:        req.Vf,
		Kj:        req.Kj,
		K:         req.K,
		V:         v,
		Q:         q,
		Side:      side,
		Congested: congested,
	})
}
