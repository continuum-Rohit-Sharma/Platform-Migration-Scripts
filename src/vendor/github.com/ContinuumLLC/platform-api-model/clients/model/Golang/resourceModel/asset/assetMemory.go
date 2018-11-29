package asset

//AssetMemory is the struct definition of /resources/asset/assetMemory
type AssetMemory struct {
	TotalPhysicalMemoryBytes     int64 `json:"totalPhysicalMemoryBytes" cql:"total_physical_memory_bytes"`
	TotalVirtualMemoryBytes      int64 `json:"totalVirtualMemoryBytes" cql:"total_virtual_memory_bytes"`
	AvailableVirtualMemoryBytes  int64 `json:"availableVirtualMemoryBytes" cql:"available_virtual_memory_bytes"`
	AvailablePhysicalMemoryBytes int64 `json:"availablePhysicalMemoryBytes" cql:"available_physical_memory_bytes"`
	TotalPageFileSpaceBytes      int64 `json:"totalPageFileSpaceBytes" cql:"total_page_file_space_bytes"`
	AvailablePageFileSpaceBytes  int64 `json:"availablePageFileSpaceBytes" cql:"available_page_file_space_bytes"`
}
