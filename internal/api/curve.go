package api

import (
	"encoding/json"
	"net/http"

	"greenshields/internal/core"
)

type curveRequest struct {
	Vf    float64 `json:"vf"`
	Kj    float64 `json:"kj"`
	Steps int     `json:"steps"`
}

type pointDTO struct {
	K float64 `json:"k"`
	V float64 `json:"v"`
	Q float64 `json:"q"`
}

type curveResponse struct {
	Vf     float64    `json:"vf"`
	Kj     float64    `json:"kj"`
	Km     float64    `json:"km"`
	Qm     float64    `json:"qm"`
	Vm     float64    `json:"vm"`
	Points []pointDTO `json:"points"`
}

// handleCurve returns the sampled q(k) and v(k) curves and the capacity point.
func (s *Server) handleCurve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}
	var req curveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}
	m, err := core.New(req.Vf, req.Kj)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	pts := m.Curve(req.Steps)
	qm, km := m.Capacity()
	dto := make([]pointDTO, 0, len(pts))
	for _, p := range pts {
		dto = append(dto, pointDTO{K: p.K, V: p.V, Q: p.Q})
	}
	writeJSON(w, http.StatusOK, curveResponse{
		Vf:     req.Vf,
		Kj:     req.Kj,
		Km:     km,
		Qm:     qm,
		Vm:     m.SpeedAtCapacity(),
		Points: dto,
	})
}
