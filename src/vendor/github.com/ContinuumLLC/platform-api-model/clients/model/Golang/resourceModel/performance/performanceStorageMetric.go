package performance

type PerformanceStorageMetric struct {
	IdleTime                float64 `json:"idleTime" cql:"idle_time"`
	WriteCompleted          int64   `json:"writeCompleted" cql:"write_completed"`
	WriteTimeMs             int64   `json:"writeTimeMs" cql:"write_time_ms"`
	AvgDiskWriteQueueLength uint64  `json:"avgDiskWriteQueueLength" cql:"avg_disk_write_queue_length"`
	AvgDiskSecPerWrite      uint32  `json:"avgDiskSecPerWrite" cql:"avg_disk_sec_per_write"`
	ReadCompleted           int64   `json:"readCompleted" cql:"read_completed"`
	ReadTimeMs              int64   `json:"readTimeMs" cql:"read_time_ms"`
	AvgDiskReadQueueLength  uint64  `json:"avgDiskReadQueueLength" cql:"avg_disk_read_queue_length"`
	AvgDiskSecPerRead       uint32  `json:"avgDiskSecPerRead" cql:"avg_disk_sec_per_read"`
	FreeSpaceBytes          int64   `json:"freeSpaceBytes" cql:"free_space_bytes"`
	UsedSpaceBytes          int64   `json:"usedSpaceBytes" cql:"used_space_bytes"`
	TotalSpaceBytes         int64   `json:"totalSpaceBytes" cql:"total_space_bytes"`
	DiskTimeTotal           float64 `json:"diskTimeTotal" cql:"disk_time_total"`
}
