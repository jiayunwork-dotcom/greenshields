package api

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() *http.ServeMux {
	srv := NewServer("example")
	mux := http.NewServeMux()
	srv.Register(mux, "web")
	return mux
}

// TestQKEndpoint checks the /api/qk success path and the invalid-input error
// path. (Medium)
func TestQKEndpoint(t *testing.T) {
	mux := newTestServer()

	// Success: k=45 on vf=120,kj=180 -> v=90, q=4050, free side.
	req := httptest.NewRequest(http.MethodPost, "/api/qk", strings.NewReader(`{"vf":120,"kj":180,"k":45}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/api/qk status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp qkResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.V != 90 {
		t.Errorf("v = %g, want 90", resp.V)
	}
	if resp.Q != 4050 {
		t.Errorf("q = %g, want 4050", resp.Q)
	}
	if resp.Side != "free" {
		t.Errorf("side = %q, want free", resp.Side)
	}

	// Failure: non-positive vf -> 400 with error JSON.
	req2 := httptest.NewRequest(http.MethodPost, "/api/qk", strings.NewReader(`{"vf":-1,"kj":180,"k":45}`))
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusBadRequest {
		t.Errorf("/api/qk bad status = %d, want 400", rec2.Code)
	}
	if !strings.Contains(rec2.Body.String(), "error") {
		t.Errorf("bad response missing error field: %s", rec2.Body.String())
	}
}

// TestCurveEndpoint checks the /api/curve success path returns the capacity
// point and a point list. (Medium)
func TestCurveEndpoint(t *testing.T) {
	mux := newTestServer()
	req := httptest.NewRequest(http.MethodPost, "/api/curve", strings.NewReader(`{"vf":120,"kj":180,"steps":21}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/api/curve status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp curveResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Qm != 120*180/4 {
		t.Errorf("qm = %g, want %g", resp.Qm, float64(120*180/4))
	}
	if resp.Km != 90 {
		t.Errorf("km = %g, want 90", resp.Km)
	}
	if len(resp.Points) != 21 {
		t.Errorf("points = %d, want 21", len(resp.Points))
	}
}

// TestWaveEndpoint checks the /api/wave success path (model-derived flows)
// and the equal-density failure path. (Medium)
func TestWaveEndpoint(t *testing.T) {
	mux := newTestServer()

	req := httptest.NewRequest(http.MethodPost, "/api/wave", strings.NewReader(`{"vf":120,"kj":180,"k1":45,"k2":162}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("/api/wave status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp waveResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// k1=45 -> q1=4050; k2=162 -> q2=120*162*(1-162/180)=120*162*0.1=1944.
	if math.Abs(resp.Q1-4050) > 1e-6 {
		t.Errorf("q1 = %g, want 4050", resp.Q1)
	}
	if math.Abs(resp.Q2-1944) > 1e-6 {
		t.Errorf("q2 = %g, want 1944", resp.Q2)
	}
	// w = (1944-4050)/(162-45) = -2106/117 = -18.
	if math.Abs(resp.W-(-18)) > 1e-6 {
		t.Errorf("w = %g, want -18", resp.W)
	}
	if resp.Direction != "upstream" {
		t.Errorf("direction = %q, want upstream", resp.Direction)
	}

	// Equal densities -> error.
	req2 := httptest.NewRequest(http.MethodPost, "/api/wave", strings.NewReader(`{"vf":120,"kj":180,"k1":90,"k2":90}`))
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusBadRequest {
		t.Errorf("/api/wave equal-density status = %d, want 400", rec2.Code)
	}
}
