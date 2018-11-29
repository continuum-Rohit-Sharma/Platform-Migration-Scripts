package performance

import "time"

//PerformanceStorage stores data related to Disk
type PerformanceStorage struct {
	CreateTimeUTC time.Time                     `json:"createTimeUTC"`
	CreatedBy     string                        `json:"createdBy"`
	Index         int                           `json:"index"`
	Name          string                        `json:"name" cql:"name"`
	Type          string                        `json:"type"`
	Metric        PerformanceStorageMetric      `json:"metric" cql:"metric"`
	Partitions    []PerformanceStoragePartition `json:"partitions" cql:"partitions"`
}
