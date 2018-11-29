package registry

import "time"

// Data is the struct definition of registry state
type Data struct {
	Path      string     `json:"path"`
	Exist     bool       `json:"exist"`
	Timestamp *time.Time `json:"timestamp"`
	Values    []Value    `json:"values,omitempty"`
}
