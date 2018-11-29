package performance

import "time"

//PerformanceStorage stores data related to Disk
type PerformanceStorages struct {
	CreateTimeUTC time.Time            `json:"createTimeUTC"`
	CreatedBy     string               `json:"createdBy"`
	Name          string               `json:"name"`
	Type          string               `json:"type"`
	EndpointID    string               `json:"endpointID"`
	PartnerID     string               `json:"partnerID"`
	ClientID      string               `json:"clientID"`
	SiteID        string               `json:"siteID"`
	Storages      []PerformanceStorage `json:"storages" cql:"storages"`
}
