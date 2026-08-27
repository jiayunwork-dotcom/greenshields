package godunov

import "greenshields/internal/core"

func (s *Solver) Demand(k float64) (float64, error) {
	if err := s.Model.ValidateDensity(k); err != nil {
		return 0, err
	}
	if k <= s.Km()+core.Epsilon*s.Model.Kj {
		return s.Q(k)
	}
	return s.Qm(), nil
}

func (s *Solver) Supply(k float64) (float64, error) {
	if err := s.Model.ValidateDensity(k); err != nil {
		return 0, err
	}
	if k <= s.Km()+core.Epsilon*s.Model.Kj {
		return s.Qm(), nil
	}
	return s.Q(k)
}

func (s *Solver) Flux(kL, kR float64) (float64, error) {
	d, err := s.Demand(kL)
	if err != nil {
		return 0, err
	}
	sup, err := s.Supply(kR)
	if err != nil {
		return 0, err
	}
	if d < sup {
		return d, nil
	}
	return sup, nil
}

func (s *Solver) SendingBound(k float64) (float64, error) {
	d, err := s.Demand(k)
	if err != nil {
		return 0, err
	}
	q, err := s.Q(k)
	if err != nil {
		return 0, err
	}
	if d < q {
		return d, nil
	}
	if d > s.Qm()+core.Epsilon*s.Qm() {
		return s.Qm(), nil
	}
	return d, nil
}

func (s *Solver) ReceivingBound(k float64) (float64, error) {
	sup, err := s.Supply(k)
	if err != nil {
		return 0, err
	}
	if sup > s.Qm()+core.Epsilon*s.Qm() {
		return s.Qm(), nil
	}
	return sup, nil
}

func (s *Solver) AtCapacity(k float64) bool {
	km := s.Km()
	return relClose(k, km, core.Epsilon)
}
