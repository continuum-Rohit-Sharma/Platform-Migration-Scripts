package agent

import "time"

// EndpointHeartbeat struct processed by agent service
type EndpointHeartbeat struct {
	EndpointID       string
	DcDateTimeUTC    time.Time
	AgentDateTimeUTC time.Time
	HeartbeatCounter int64
	PublicIPAddress  string
	Availability     bool
}
