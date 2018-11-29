package sysEvent

// RegEntry is the struct definition of registry entry
// that will be send by sysEvent plugin to registry plugin
// to recieve actual state of correspond registry entity
type RegEntry struct {
	Root   string   `json:"root"`
	Key    string   `json:"key,omitempty"`
	Values []string `json:"values,omitempty"`
}
