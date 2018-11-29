package asset

//AssetProcessor is the struct definition of /resources/asset/assetProcessor
type AssetProcessor struct {
	Product       string  `json:"product,omitempty" cql:"product"`
	NumberOfCores int     `json:"numberOfCores" cql:"number_of_cores"`
	ClockSpeedMhz float64 `json:"clockSpeedMhz" cql:"clock_speed_mhz"`
	Family        int     `json:"family" cql:"family"`
	Manufacturer  string  `json:"manufacturer,omitempty" cql:"manufacturer"`
	ProcessorType string  `json:"processorType,omitempty" cql:"processor_type"`
	SerialNumber  string  `json:"serialNumber,omitempty" cql:"serial_number"`
	Level         int     `json:"level" cql:"level"`
}
