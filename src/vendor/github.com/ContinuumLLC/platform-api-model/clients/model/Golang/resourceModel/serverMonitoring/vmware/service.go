package vmware

// Service describes Host system service
type Service struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	State        string `json:"state,omitempty"`
}
