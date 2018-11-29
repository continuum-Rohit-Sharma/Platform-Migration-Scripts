package asset

// PhysicalMemory is the struct definition of /resources/asset/assetPhysicalMemory
type PhysicalMemory struct {
	Manufacturer string `json:"manufacturer,omitempty" cql:"manufacturer"`
	SerialNumber string `json:"serialNumber,omitempty" cql:"serial_number"`
	SizeBytes    uint64 `json:"sizeBytes" cql:"size_bytes"`
}
