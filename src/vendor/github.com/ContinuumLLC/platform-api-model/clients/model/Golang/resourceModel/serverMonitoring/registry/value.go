package registry

// Value describes Registry value
type Value struct {
	Name  string `json:"name"`
	Data  string `json:"data,omitempty"`
	Exist bool   `json:"exist"`
}
