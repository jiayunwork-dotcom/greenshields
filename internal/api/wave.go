package api

import (
	"encoding/json"
	"net/http"

	"greenshields/internal/core"
	"greenshields/internal/wave"
)

type waveRequest struct {
	// Vf and Kj are only needed when the flows are not supplied explicitly.
	Vf *float64 `json:"vf"`
	Kj *float64 `json:"kj"`
	// K1, K2 are the densities of the two states.
	K1 float64 `json:"k1"`
	K2 float64 `json:"k2"`
	// Q1, Q2 are the explicit flows. When nil they are derived from (Vf, Kj).
	Q1 *float64 `json:"q1"`
	Q2 *float64 `json:"q2"`
}

type waveResponse struct {
	K1        float64 `json:"k1"`
	Q1        float64 `json:"q1"`
	K2        float64 `json:"k2"`
	Q2        float64 `json:"q2"`
	W         float64 `json:"w"`
	Direction string  `json:"direction"`
}

// handleWave computes the shock speed between two traffic states. Flows may be
// supplied explicitly (as part of the two states) or derived from a
// Greenshields model given vf and kj.
func (s *Server) handleWave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}
	var req waveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	var q1, q2 float64
	if req.Q1 != nil {
		q1 = *req.Q1
	} else {
		if req.Vf == nil || req.Kj == nil {
			writeError(w, http.StatusBadRequest, "missing q1 and no (vf,kj) to derive it")
			return
		}
		m, err := core.New(*req.Vf, *req.Kj)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := m.ValidateDensity(req.K1); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		q1, _ = m.Flow(req.K1)
	}
	if req.Q2 != nil {
		q2 = *req.Q2
	} else {
		if req.Vf == nil || req.Kj == nil {
			writeError(w, http.StatusBadRequest, "missing q2 and no (vf,kj) to derive it")
			return
		}
		m, err := core.New(*req.Vf, *req.Kj)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := m.ValidateDensity(req.K2); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		q2, _ = m.Flow(req.K2)
	}

	a := wave.State{K: req.K1, Q: q1}
	b := wave.State{K: req.K2, Q: q2}
	wSpeed, err := wave.ShockSpeed(a, b)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, waveResponse{
		K1:        req.K1,
		Q1:        q1,
		K2:        req.K2,
		Q2:        q2,
		W:         wSpeed,
		Direction: wave.Direction(wSpeed),
	})
}
