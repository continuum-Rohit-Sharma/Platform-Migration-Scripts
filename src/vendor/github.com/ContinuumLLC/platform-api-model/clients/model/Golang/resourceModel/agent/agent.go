package agent

import "time"

//Agent is the struct definition of /resources/agent/agent
type Agent struct {
	TimeStampUTC time.Time   `json:"timeStampUTC"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Components   []Component `json:"components"`
	Plugins      []Plugin    `json:"plugins"`
}
