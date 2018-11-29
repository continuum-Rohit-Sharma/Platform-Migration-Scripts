package performance

import "time"

//Plugin is the struct definition of /resources/performance/performanceDomainVersion
type performanceDomainVersion struct {
	TimeStampUTC   time.Time           `json:"timeStampUTC"`
	Name           string              `json:"name"`
	Type           string              `json:"type"`
	ServiceName    string              `json:"serviceName"`
	ServiceVersion string              `json:"serviceVersion"`
	Plugins        []PerformancePlugin `json:"plugins"`
}
