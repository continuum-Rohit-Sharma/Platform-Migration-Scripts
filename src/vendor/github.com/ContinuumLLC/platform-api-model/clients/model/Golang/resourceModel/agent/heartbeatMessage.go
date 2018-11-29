package agent

import (
	"time"
)

//HeartbeatMessage contains the messageID & URL to fetch further message data
type HeartbeatMessage struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	TimestampUTC time.Time `json:"timestampUTC"`
	Index        int       `json:"index"`
	MessageID    string    `json:"messageID"`
	MessageURL   string    `json:"messageURL"`
}

//EventMessage is pushed to managed_endpoints_change kafka topic on certain events
type EventMessage struct {
	MessageType      string
	EndpointMapping  EndpointMapping
	DcDateTimeUTC    time.Time
	AgentDateTimeUTC time.Time
}
