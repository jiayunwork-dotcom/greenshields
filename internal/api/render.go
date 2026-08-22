package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"greenshields/internal/core"
	"greenshields/internal/example"
)

const (
	svgWidth  = 520
	svgHeight = 360
	svgPad    = 44
	svgSteps  = 80
)

// SVGCurve builds an SVG string of the q(k) parabola. Every coordinate is
// derived from the model's Curve output; nothing is hard-coded, so the picture
// always reflects the supplied parameters. The capacity point is marked with a
// red circle.
func SVGCurve(m *core.Model) string {
	qm, km := m.Capacity()
	if qm <= 0 {
		qm = 1
	}
	pts := m.Curve(svgSteps)
	plotW := float64(svgWidth - 2*svgPad)
	plotH := float64(svgHeight - 2*svgPad)

	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`,
		svgWidth, svgHeight, svgWidth, svgHeight))
	b.WriteString(`<rect width="100%" height="100%" fill="#ffffff"/>`)

	// Axes.
	x0 := svgPad
	y0 := svgHeight - svgPad
	x1 := svgWidth - svgPad
	y1 := svgPad
	b.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#444" stroke-width="1"/>`, x0, y0, x1, y0))
	b.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#444" stroke-width="1"/>`, x0, y0, x0, y1))

	// Axis labels.
	b.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="12" fill="#444">k</text>`, x1-10, y0+18))
	b.WriteString(fmt.Sprintf(`<text x="%d" y="%d" font-size="12" fill="#444">q</text>`, x0-30, y1+4))

	// Parabola polyline.
	coords := make([]string, 0, len(pts))
	for _, p := range pts {
		x := float64(svgPad) + (p.K/m.Kj)*plotW
		y := float64(svgHeight-svgPad) - (p.Q/qm)*plotH
		coords = append(coords, fmt.Sprintf("%.1f,%.1f", x, y))
	}
	b.WriteString(fmt.Sprintf(`<polyline fill="none" stroke="#1f77b4" stroke-width="2" points="%s"/>`, strings.Join(coords, " ")))

	// Capacity marker.
	cx := float64(svgPad) + (km/m.Kj)*plotW
	cy := float64(svgHeight-svgPad) - (qm/qm)*plotH
	b.WriteString(fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="4" fill="#d62728"/>`, cx, cy))
	b.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="11" fill="#d62728">qm=%.2g</text>`, cx+6, cy-6, qm))

	b.WriteString(`</svg>`)
	return b.String()
}

// handleCurveSVG serves a server-rendered SVG of the q(k) parabola. It uses
// the same model math as /api/curve.
func (s *Server) handleCurveSVG(w http.ResponseWriter, r *http.Request) {
	var req curveRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Vf <= 0 || req.Kj <= 0 {
		p, err := example.Embedded()
		if err != nil || p.Vf <= 0 {
			writeError(w, http.StatusBadRequest, "missing or invalid vf/kj")
			return
		}
		req.Vf, req.Kj = p.Vf, p.Kj
	}
	m, err := core.New(req.Vf, req.Kj)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(SVGCurve(m)))
}
