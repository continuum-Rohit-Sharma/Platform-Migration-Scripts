package vmware

//Processor is the struct definition of the Host system processor
type Processor struct {
	Model            string `json:"model,omitempty"`
	Vendor           string `json:"vendor,omitempty"`
	Level            int    `json:"level,omitempty"`
	NumberOfPackages int    `json:"numberOfPackages,omitempty"`
	NumberOfCores    int    `json:"numberOfCores,omitempty"`
	NumberOfThreads  int    `json:"numberOfThreads,omitempty"`
	ClockSpeedHz     int64  `json:"clockSpeedHz,omitempty"`
}
