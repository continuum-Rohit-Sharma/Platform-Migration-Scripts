package autoupdate

//FeaturePackages is a struct defining association between features and their packages
type FeaturePackages struct {
	Name     string   `json:"name,omitempty" cql:"name"`
	Packages []string `json:"packages,omitempty" cql:"packages"`
}

//EndpointPackages is a struct defining association between endpoint and their packages
type EndpointPackages struct {
	EndpointID string   `json:"endpointID,omitempty" cql:"endpoint_id"`
	Packages   []string `json:"packages,omitempty" cql:"packages"`
}
