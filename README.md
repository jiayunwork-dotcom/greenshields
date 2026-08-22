# greenshields

A small, self-contained calculator for the **Greenshields traffic-flow
fundamental diagram**. Given a free-flow speed `vf` and a jam density `kj`, it
computes the linear speed–density relationship, the parabolic flow–density
curve, the capacity point, and the kinematic-wave (shock) speed between two
traffic states.

This is a traffic-flow *kernel*, not a signal-timing or queueing product: it
models the continuous relationship between speed, density, and flow on a single
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
upstream.

## Inputs, outputs, and boundaries

- `vf` and `kj` must both be strictly positive.
- Density `k` must lie in `[0, kj]`.
- A flow `q` larger than the capacity `qm` cannot be produced by any density;
  solving `q(k) = q` for `k` then returns no real roots and an error.
- Given a flow `q < qm`, solving for density yields **two** roots: a free-flow
  branch (`k < km`) and a congested branch (`k > km`), which always satisfy
  `k_free + k_congested = kj`.

All of these conditions are checked up front; invalid input produces a clear
error instead of a wrong number.

## Usage

Start the web server and JSON API:

```bash
go run . -http :8080
```

Then open <http://localhost:8080/>. The page can **load the example**
(`example/freeway.json`, a highway with `vf = 120`, `kj = 180`), draw the
`q–k` parabola from the backend point list, and compute capacity, single-point
speed/flow, and wave speed.

You can also query the API directly. A reproducible command using the bundled
example first prints the capacity point without starting the server:

```bash
go run . -print
# example "urban-freeway"
#   vf = 120  kj = 180
#   capacity: qm = 5400  at km = 90  (vm = 60)
```

And a couple of API calls:

```bash
curl -s -X POST localhost:8080/api/curve -d '{"vf":120,"kj":180}' | head -c 200; echo
curl -s -X POST localhost:8080/api/qk    -d '{"vf":120,"kj":180,"k":45}'
curl -s -X POST localhost:8080/api/wave  -d '{"vf":120,"kj":180,"k1":45,"k2":162}'
```

Invalid input is reported as a JSON error body with a non-2xx status, for
example:

```bash
curl -s -X POST localhost:8080/api/qk -d '{"vf":-1,"kj":180,"k":45}'
# {"error":"greenshields: free-flow speed vf must be greater than 0: got vf=-1"}
```

## API endpoints

| Method | Path           | Body                          | Returns                                   |
|--------|----------------|-------------------------------|-------------------------------------------|
| POST   | `/api/qk`      | `{vf,kj,k}`                   | `{v,q,side,congested}`                    |
| POST   | `/api/curve`   | `{vf,kj,steps?}`              | `{km,qm,vm,points:[{k,v,q}...]}`          |
| POST   | `/api/wave`    | `{vf?,kj?,k1,q1?,k2,q2?}`     | `{k1,q1,k2,q2,w,direction}`               |
| GET    | `/api/example` | —                             | the bundled `freeway.json`                |
| GET    | `/api/curve.svg` | `{vf,kj}`                   | an SVG of the `q–k` parabola              |

For `/api/wave`, supply the flows `q1`, `q2` explicitly, or omit them and let the
server derive them from `vf`, `kj`, and the two densities.

## Build and test

```bash
go build ./...
go test ./...
```
