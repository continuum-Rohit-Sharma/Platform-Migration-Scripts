package performance

import "time"

type PerformanceProcessors struct {
	CreateTimeUTC time.Time                  `json:"createTimeUTC"`
	CreatedBy     string                     `json:"createdBy"`
	Index         int                        `json:"index"`
	Name          string                     `json:"name"`
	Type          string                     `json:"type"`
	EndpointID    string                     `json:"endpointID"`
	PartnerID     string                     `json:"partnerID"`
	ClientID      string                     `json:"clientID"`
	SiteID        string                     `json:"siteID"`
	Metric        PerformanceProcessorMetric `json:"metric" cql:"metric"`
	CPUs          []PerformanceProcessor     `json:"cpus" cql:"processors"`
}
