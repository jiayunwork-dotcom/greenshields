# greenshields

A Go HTTP service for the Greenshields traffic-flow fundamental diagram.
Given a free-flow speed `vf` and a jam density `kj`, it computes the linear
speed–density relationship, the parabolic flow–density curve, the capacity
point, and the kinematic-wave (shock) speed between two traffic states. The
same process also hosts a small static page under `web/` that calls the JSON
API. This is a traffic-flow kernel, not a signal-timing product: it models
the continuous relationship between speed, density, and flow on a single
homogeneous link.

## What it computes

Greenshields assumes a straight speed–density line:

```
v(k) = vf · (1 − k/kj)
q(k) = k · v(k)            # flow–density parabola
qm   = vf · kj / 4         # capacity flow, at km = kj/2
vm   = vf / 2              # speed at capacity
```

The wave speed between two states `(k₁, q₁)` and `(k₂, q₂)` is the slope of the
chord joining them on the flow–density plane:

```
w = (q₂ − q₁) / (k₂ − k₁)
```

A negative `w` means a disturbance (for example a downstream jam) propagates
upstream. The Godunov flux on the same diagram is `min(demand(kL), supply(kR))`
and must match the entropy flux of the concave parabola (interval min on a
shock, interval max — hence capacity — on a rarefaction that straddles `km`).

Illegal parameters (`vf ≤ 0`, `kj ≤ 0`, density outside `[0, kj]`, flow above
`qm`, equal densities for a shock) are rejected. HTTP callers receive **422**
with a JSON `{"error":"..."}` body. Malformed JSON is **400**.

## HTTP API

Start the server (default listen address `:8080`):

```bash
go run . -http :8080
```

Then open <http://localhost:8080/>. The page can load the example
(`example/freeway.json`, a highway with `vf = 120`, `kj = 180`), draw the
`q–k` parabola from the backend point list, and compute capacity, single-point
speed/flow, and wave speed.

A CLI print path is still available and exits after loading the example:

```bash
go run . -print
# example "urban-freeway"
#   vf = 120  kj = 180
#   capacity: qm = 5400  at km = 90  (vm = 60)
```

| Method | Path             | Body                          | Returns                                   |
|--------|------------------|-------------------------------|-------------------------------------------|
| GET    | `/api/health`    | —                             | `{"status":"ok"}`                         |
| POST   | `/api/qk`        | `{vf,kj,k}`                   | `{v,q,side,congested}`                    |
| POST   | `/api/curve`     | `{vf,kj,steps?}`              | `{km,qm,vm,points:[{k,v,q}...]}`          |
| POST   | `/api/wave`      | `{vf?,kj?,k1,q1?,k2,q2?}`     | `{k1,q1,k2,q2,w,direction}`               |
| GET    | `/api/example`   | —                             | the bundled `freeway.json`                |
| GET    | `/api/curve.svg` | `{vf,kj}`                     | an SVG of the `q–k` parabola              |

```bash
curl -s http://localhost:8080/api/health
curl -s -X POST localhost:8080/api/curve -d '{"vf":120,"kj":180}' | head -c 200; echo
curl -s -X POST localhost:8080/api/qk    -d '{"vf":120,"kj":180,"k":45}'
curl -s -X POST localhost:8080/api/wave  -d '{"vf":120,"kj":180,"k1":45,"k2":162}'
curl -s -X POST localhost:8080/api/qk    -d '{"vf":-1,"kj":180,"k":45}'
# HTTP 422  {"error":"greenshields: free-flow speed vf must be greater than 0: got vf=-1"}
```

For `/api/wave`, supply the flows `q1`, `q2` explicitly, or omit them and let the
server derive them from `vf`, `kj`, and the two densities.

A solved diagram (capacity, sample states, and one wave pair) can be written as
a versioned JSON snapshot and read back; empty or truncated files are refused.

## Build and test

```bash
go build ./...
go test ./...
```
