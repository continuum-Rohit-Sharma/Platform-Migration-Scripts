package asset

//Share is the struct definition of /resources/asset/Share
type Share struct {
	Name        string   `json:"name,omitempty" cql:"name"`
	Caption     string   `json:"caption,omitempty" cql:"caption"`
	Description string   `json:"description,omitempty" cql:"description"`
	Path        string   `json:"path,omitempty" cql:"path"`
	Access      string   `json:"access,omitempty" cql:"access"`
	UserAccess  []string `json:"userAccess,omitempty" cql:"user_access"`
	Type        []string `json:"type,omitempty" cql:"type"`
}
