package webroot

import "time"

// WebrootPluginCommandResult represents result details returned by executed command
type WebrootPluginCommandResult struct {
	ExecutionID  string    `json:"executionID"  description:"ExecutionID generated by Webroot MicroService"`
	ResultCode   int       `json:"resultCode"   description:"Return status: Success/Failed"`
	ErrorMessage string    `json:"errorMessage" description:"Blank if result code is 0 or Registry change error otherwise"`
	TimestampUTC time.Time `json:"timestampUTC" description:"UTC time when the Script execution finished"`
	MessageType  string    `json:"messageType"  description:"Type of result message, can be status or command"`
}
