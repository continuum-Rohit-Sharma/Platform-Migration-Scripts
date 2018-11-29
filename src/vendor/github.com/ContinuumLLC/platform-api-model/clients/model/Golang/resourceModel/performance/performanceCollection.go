package performance

import "time"

//PerformanceCollection is the struct definition of /resources/performance/performanceCollection
type PerformanceCollection struct {
	CreateTimeUTC time.Time               `json:"createTimeUTC"`
	CreatedBy     string                  `json:"createdBy"`
	Name          string                  `json:"name"`
	Type          string                  `json:"type"`
	EndpointID    string                  `json:"endpointID"`
	PartnerID     string                  `json:"partnerID"`
	ClientID      string                  `json:"clientID"`
	SiteID        string                  `json:"siteID"`
	Processors    []PerformanceProcessors `json:"processors"`
	Memory        []PerformanceMemory     `json:"memory"`
	Storages      []PerformanceStorages   `json:"storages"`
	Network       []PerformanceNetwork    `json:"network"`
	Processes     []PerformanceProcesses  `json:"processes"`
}
