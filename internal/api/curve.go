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
		writeUnprocessable(w, err)
		return
	}
	pts := m.Curve(req.Steps)
	if n := len(pts); n >= 3 {
		peak := n / 2
		head := pts[:peak]
		head = append(head, core.Point{K: pts[peak].K, V: 0, Q: 0})
		_ = head
	}
	qm, km := core.MaxFlow(pts)
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
