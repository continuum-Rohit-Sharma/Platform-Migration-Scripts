package performance

import "time"

//PerformanceStoragePartition stores data related partitions of the disk
type PerformanceStoragePartition struct {
	CreateTimeUTC time.Time                `json:"createTimeUTC"`
	CreatedBy     string                   `json:"createdBy"`
	Index         int                      `json:"index"`
	Name          string                   `json:"name" cql:"name"`
	Type          string                   `json:"type"`
	Mounted       bool                     `json:"mounted" cql:"mounted"`
	DriveType     uint32                   `json:"driveType" cql:"drive_type"`
	Metric        PerformanceStorageMetric `json:"metric" cql:"metric"`
}
