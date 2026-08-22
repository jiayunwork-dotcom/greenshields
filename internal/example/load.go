package example

import (
	"encoding/json"
	"fmt"
	"os"
)

// Parameters is a decoded example file. It carries the two Greenshields
// parameters plus an optional human-readable name and unit labels.
type Parameters struct {
	Name string             `json:"name"`
	Vf   float64            `json:"vf"`
	Kj   float64            `json:"kj"`
	Unit map[string]string  `json:"unit"`
}

// LoadFile reads a Greenshields example JSON file from disk.
func LoadFile(path string) (*Parameters, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("example: cannot read %s: %w", path, err)
	}
	return Decode(data)
}

// Decode parses example JSON into Parameters and validates the parameters.
func Decode(data []byte) (*Parameters, error) {
	var p Parameters
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("example: invalid json: %w", err)
	}
	if p.Vf <= 0 {
		return nil, fmt.Errorf("example: vf must be positive, got %g", p.Vf)
	}
	if p.Kj <= 0 {
		return nil, fmt.Errorf("example: kj must be positive, got %g", p.Kj)
	}
	return &p, nil
}

// Embedded returns the decoded embedded highway example.
func Embedded() (*Parameters, error) {
	return Decode(FreewayJSON)
}
