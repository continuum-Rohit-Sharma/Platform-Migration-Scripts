package agent

import "time"

//AppManagement is a struct defining the actual app management task
type AppManagement struct {
	Action      string `json:"action"`
	PackageName string `json:"packageName"`
}

//AppMessage is the struct definition of /resources/agent/appMessage
type AppMessage struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	TimestampUTC time.Time `json:"timestampUTC"`
	Path         string    `json:"path"`
	AppManagement
	MessageID string
}
