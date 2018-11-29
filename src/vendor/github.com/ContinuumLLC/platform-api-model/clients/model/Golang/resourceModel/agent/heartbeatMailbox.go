package agent

import "time"

//HeartbeatMailbox contains the slice of the messages
type HeartbeatMailbox struct {
	Name         string             `json:"name"`
	Type         string             `json:"type"`
	TimestampUTC time.Time          `json:"timestampUTC"`
	Message      []HeartbeatMessage `json:"message"`
}
