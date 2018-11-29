package agent

// ProvisioningData is the struct definition of agent provisioning data structure
type ProvisioningData struct {
	EndpointMapping
	SysInfo               SystemInfo `json:"sysInfo,omitempty"`
	PublicIPAddress       string     `json:"publicIPAddress,omitempty"`
	Token                 string     `json:"token,omitempty"`
	AgentInstalledVersion string     `json:"agentInstalledVersion,omitempty"`
}

//EndpointMapping struct provides mapping between endpoint and its partner,client,site,agent.
//This will be returned back to agent as response to successful registration.
type EndpointMapping struct {
	EndpointID  string `json:"endpointID,omitempty"`
	AgentID     string `json:"agentID,omitempty"`
	PartnerID   string `json:"partnerID,omitempty"`
	SiteID      string `json:"siteID,omitempty"`
	ClientID    string `json:"clientID,omitempty"`
	LegacyRegID string `json:"legacyRegID,omitempty"`
	Installed   bool   `json:"installed,omitempty"`
}

//MigrationData stores the old and new endpoint mapping values. Used in case of migration.
type MigrationData struct {
	Newmapping EndpointMapping `json:"newmapping,omitempty"`
	Oldmapping EndpointMapping `json:"oldmapping,omitempty"`
}
