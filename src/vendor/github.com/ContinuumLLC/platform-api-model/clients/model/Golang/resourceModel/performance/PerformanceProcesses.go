package performance

import "time"

//PerformanceProcesses stores data related to processes
type PerformanceProcesses struct {
	CreateTimeUTC time.Time            `json:"createTimeUTC"`
	CreatedBy     string               `json:"createdBy"`
	Name          string               `json:"name"`
	Type          string               `json:"type"`
	EndpointID    string               `json:"endpointID"`
	PartnerID     string               `json:"partnerID"`
	ClientID      string               `json:"clientID"`
	SiteID        string               `json:"siteID"`
	Processes     []PerformanceProcess `json:"processes"`
}
