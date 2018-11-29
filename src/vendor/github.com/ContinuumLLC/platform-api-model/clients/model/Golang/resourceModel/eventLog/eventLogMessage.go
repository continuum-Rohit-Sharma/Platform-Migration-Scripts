package eventlog

import (
	"time"
)

//EventLogMessage is the struct definition of /resources/eventLog/eventLogMessage
type EventLogMessage struct {
	Hostname     string    `json:"hostname"`
	Source       string    `json:"source"`
	Facility     int       `json:"facility"`
	Severity     int       `json:"severity"`
	EventID      int       `json:"eventID"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"createdAt"`
	Duplications int       `json:"duplications"`
}
