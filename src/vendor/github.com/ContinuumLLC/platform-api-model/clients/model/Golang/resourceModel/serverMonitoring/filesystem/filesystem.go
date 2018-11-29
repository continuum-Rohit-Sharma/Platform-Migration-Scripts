package filesystem

// Entry is the struct definition of file / directory state
type Entry struct {
	Path       string      `json:"path"`
	Exists     bool        `json:"exists"`
	Attributes *Attributes `json:"attributes,omitempty"`
}
