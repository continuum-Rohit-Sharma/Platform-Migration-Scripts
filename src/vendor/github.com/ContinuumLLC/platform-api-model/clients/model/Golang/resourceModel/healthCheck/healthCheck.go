package healthCheck

import "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/version"
import "time"

//HealthCheck is an API model to hold health check information
type HealthCheck struct {
	version.Version
	NetworkInterfaces []string  `json:"networkInterfaces"`
	Status            string    `json:"status"`
	StartTime         time.Time `json:"startTime"`
	CPUPercentage     float64   `json:"cpuPercentage"`
	NumOfOSThreads    int       `json:"numOfOSThreads"`
	MemoryPercentage  float64   `json:"memoryPercentage"`
}
