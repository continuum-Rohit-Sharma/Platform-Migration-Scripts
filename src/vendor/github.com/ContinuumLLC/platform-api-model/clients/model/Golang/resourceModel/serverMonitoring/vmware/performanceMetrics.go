package vmware

// PerformanceMetrics describes main Host performance metrics
type PerformanceMetrics struct {
	NetBytesRx                int64 `json:"netBytesRx"`
	NetBytesTx                int64 `json:"netBytesTx"`
	CPUUtilization            int64 `json:"cpuUtilization"`
	CPUTotalCapacity          int64 `json:"cpuTotalCapacity"`
	MemoryUtilizationMB       int64 `json:"memoryUtilizationMB"`
	MemoryActiveAverageMB     int64 `json:"memoryActiveAverageMB"`
	MemoryCompressedAverageMB int64 `json:"memoryCompressedAverageMB"`
	MemoryConsumedAverageMB   int64 `json:"memoryConsumedAverageMB"`
	MemoryOverheadAverageMB   int64 `json:"memoryOverheadAverageMB"`
	MemorySharedAverageMB     int64 `json:"memorySharedAverageMB"`
	MemorySwapInAverageMB     int64 `json:"memorySwapInAverageMB"`
	MemorySwapOutAverageMB    int64 `json:"memorySwapOutAverageMB"`
	MemorySwapUsedAverageMB   int64 `json:"memorySwapUsedAverageMB"`
	MemoryVmmemctlAverageMB   int64 `json:"memoryVmmemctlAverageMB"`
	MemoryTotalCapacityMB     int64 `json:"memoryTotalCapacityMB"`
}
