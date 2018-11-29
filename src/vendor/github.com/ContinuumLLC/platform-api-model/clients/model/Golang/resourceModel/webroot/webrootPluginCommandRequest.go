package webroot

import "time"

// WebrootPluginCommandRequest represents command execution request
type WebrootPluginCommandRequest struct {
	ExecutionID string    `json:"executionID"`
	Command     string    `json:"command"`
	Parameters  string    `json:"parameters"`
	ElapsedTime time.Time `json:"elapsedTime"`
}
