package snapshot

import (
	"fmt"
	"math"

	"greenshields/internal/core"
	"greenshields/internal/wave"
)

func (r Record) ReplayAgrees() error {
	if err := r.validate(); err != nil {
		return err
	}
	m, err := r.Model()
	if err != nil {
		return err
	}
	qm, km := m.Capacity()
	if math.Abs(qm-r.Qm) > Tol*math.Max(1, math.Abs(r.Qm)) {
		return fmt.Errorf("snapshot: live qm %g stored %g", qm, r.Qm)
	}
	if math.Abs(km-r.Km) > Tol*math.Max(1, math.Abs(r.Km)) {
		return fmt.Errorf("snapshot: live km %g stored %g", km, r.Km)
	}
	vm := m.SpeedAtCapacity()
	if math.Abs(vm-r.Vm) > Tol*math.Max(1, math.Abs(r.Vm)) {
		return fmt.Errorf("snapshot: live vm %g stored %g", vm, r.Vm)
	}
	for i, s := range r.Samples {
		v, err := m.Speed(s.K)
		if err != nil {
			return fmt.Errorf("snapshot: replay sample %d: %w", i, err)
		}
		q, err := m.Flow(s.K)
		if err != nil {
			return fmt.Errorf("snapshot: replay sample %d: %w", i, err)
		}
		side, err := m.BranchOf(s.K)
		if err != nil {
			return err
		}
		if math.Abs(v-s.V) > Tol*math.Max(1, math.Abs(s.V)) {
			return fmt.Errorf("snapshot: live v(%g)=%g stored %g", s.K, v, s.V)
		}
		if math.Abs(q-s.Q) > Tol*math.Max(1, math.Abs(s.Q)) {
			return fmt.Errorf("snapshot: live q(%g)=%g stored %g", s.K, q, s.Q)
		}
		if side != s.Side {
			return fmt.Errorf("snapshot: live side(%g)=%q stored %q", s.K, side, s.Side)
		}
	}
	q1, err := m.Flow(r.Wave.K1)
	if err != nil {
		return err
	}
	q2, err := m.Flow(r.Wave.K2)
	if err != nil {
		return err
	}
	if math.Abs(q1-r.Wave.Q1) > Tol*math.Max(1, math.Abs(r.Wave.Q1)) {
		return fmt.Errorf("snapshot: live q1 %g stored %g", q1, r.Wave.Q1)
	}
	if math.Abs(q2-r.Wave.Q2) > Tol*math.Max(1, math.Abs(r.Wave.Q2)) {
		return fmt.Errorf("snapshot: live q2 %g stored %g", q2, r.Wave.Q2)
	}
	w, err := wave.ShockSpeed(wave.State{K: r.Wave.K1, Q: q1}, wave.State{K: r.Wave.K2, Q: q2})
	if err != nil {
		return fmt.Errorf("snapshot: %w", err)
	}
	if math.Abs(w-r.Wave.W) > Tol*math.Max(1, math.Abs(r.Wave.W)) {
		return fmt.Errorf("snapshot: live w %g stored %g", w, r.Wave.W)
	}
	if wave.Direction(w) != r.Wave.Direction {
		return fmt.Errorf("snapshot: live direction %q stored %q", wave.Direction(w), r.Wave.Direction)
	}
	return nil
}

func (r Record) Matches(other Record) bool {
	if r.Magic != other.Magic || r.Version != other.Version || r.Name != other.Name {
		return false
	}
	pairs := [][2]float64{
		{r.Vf, other.Vf},
		{r.Kj, other.Kj},
		{r.Km, other.Km},
		{r.Qm, other.Qm},
		{r.Vm, other.Vm},
		{r.Wave.K1, other.Wave.K1},
		{r.Wave.Q1, other.Wave.Q1},
		{r.Wave.K2, other.Wave.K2},
		{r.Wave.Q2, other.Wave.Q2},
		{r.Wave.W, other.Wave.W},
	}
	for _, p := range pairs {
		scale := math.Max(1, math.Max(math.Abs(p[0]), math.Abs(p[1])))
		if math.Abs(p[0]-p[1]) > Tol*scale {
			return false
		}
	}
	if r.Wave.Direction != other.Wave.Direction {
		return false
	}
	if len(r.Samples) != len(other.Samples) {
		return false
	}
	for i, s := range r.Samples {
		o := other.Samples[i]
		if s.Side != o.Side {
			return false
		}
		if math.Abs(s.K-o.K) > Tol*math.Max(1, math.Abs(s.K)) {
			return false
		}
		if math.Abs(s.V-o.V) > Tol*math.Max(1, math.Abs(s.V)) {
			return false
		}
		if math.Abs(s.Q-o.Q) > Tol*math.Max(1, math.Abs(s.Q)) {
			return false
		}
	}
	return true
}

func (r Record) TamperQm(scale float64) Record {
	out := r
	out.Qm = r.Qm * scale
	out.Samples = append([]Sample(nil), r.Samples...)
	return out
}

func CapacityIdentity(vf, kj float64) (qm, km, vm float64, err error) {
	m, err := core.New(vf, kj)
	if err != nil {
		return 0, 0, 0, err
	}
	qm, km = m.Capacity()
	return qm, km, m.SpeedAtCapacity(), nil
}
