package performance

import "time"

type PerformanceProcessorCore struct {
	CreateTimeUTC time.Time                  `json:"createTimeUTC"`
	CreatedBy     string                     `json:"createdBy"`
	Index         int                        `json:"index" cql:"core_index"`
	Name          string                     `json:"name" cql:"core_name"`
	Metric        PerformanceProcessorMetric `json:"metric" cql:"metric"`
	Type          string                     `json:"type"`
}
