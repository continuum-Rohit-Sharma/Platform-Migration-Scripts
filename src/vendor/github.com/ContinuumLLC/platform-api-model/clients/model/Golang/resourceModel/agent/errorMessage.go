package agent

import "time"

//ErrorMessage is the struct definition of Error message structure
type ErrorMessage struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Version      string            `json:"version"`
	TimeUUID     string            `json:"timeUUID"`
	TimestampUTC time.Time         `json:"timestampUTC"`
	Path         string            `json:"path"`
	ErrorTrace   string            `json:"errorTrace"`
	StatusCode   int               `json:"statusCode"`
	ErrorCode    string            `json:"errorCode"`
	ErrorData    map[string]string `json:"errorData"`
}
