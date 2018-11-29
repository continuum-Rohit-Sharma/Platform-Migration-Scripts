package version

import "time"

//Version is the struct definition of /resources/service-name/version
type Version struct {
	TimeStampUTC    time.Time `json:"timeStampUTC"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	ServiceName     string    `json:"serviceName"`
	ServiceProvider string    `json:"serviceProvider"`
	ServiceVersion  string    `json:"serviceVersion"`
	BuildNumber     string    `json:"buildNumber"`
}
