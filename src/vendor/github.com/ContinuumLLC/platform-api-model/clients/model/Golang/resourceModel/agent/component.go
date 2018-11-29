package agent

import "time"

//Component is the struct definition of /resources/agent/component
type Component struct {
	TimeStampUTC     time.Time `json:"timeStampUTC"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	ComponentVersion string    `json:"componentVersion"`
}
