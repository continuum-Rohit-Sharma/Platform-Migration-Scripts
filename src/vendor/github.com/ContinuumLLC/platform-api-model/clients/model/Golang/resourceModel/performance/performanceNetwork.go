package performance

import "time"

//PerformanceNetwork is the struct definition of /resources/performance/performanceNetwork
type PerformanceNetwork struct {
	CreateTimeUTC time.Time                     `json:"createTimeUTC"`
	CreatedBy     string                        `json:"createdBy"`
	Index         int                           `json:"index"`
	Name          string                        `json:"name"`
	Type          string                        `json:"type"`
	EndpointID    string                        `json:"endpointID"`
	PartnerID     string                        `json:"partnerID"`
	ClientID      string                        `json:"clientID"`
	SiteID        string                        `json:"siteID"`
	Interface     []PerformanceNetworkInterface `json:"interface" cql:"interface_list"`
}
