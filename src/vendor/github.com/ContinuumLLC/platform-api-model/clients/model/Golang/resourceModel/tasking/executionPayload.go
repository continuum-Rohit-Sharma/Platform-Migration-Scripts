package tasking

import (
	"time"
	"github.com/gocql/gocql"
)

// ExecutionPayload would be sent from Tasking MS to Origin microservice corresponds to the type of the Task
// The Origin MS should support API POST route /partners/{partnerID}/executions
type ExecutionPayload struct {
	ExecutionID      string            `json:"executionId"          valid:"required,uuid"`
	OriginID         string            `json:"originId"             valid:"required,uuid"`
	ManagedEndpoints []ManagedEndpoint `json:"managedEndpoints"     valid:"required"`
	WebhookURL       string            `json:"webHookURL"           valid:"required,url"`
	Parameters       string            `json:"parameters"           valid:"-"`
}

// ManagedEndpoint struct is used to display id and nextRunTime of each ManagedEndpoint in ExecutionPayload
type ManagedEndpoint struct {
	ID          string    `json:"id"             valid:"required"`
	NextRunTime time.Time `json:"nextRunTime"    valid:"-"`
}

// ExecutionPayloadResult would be sent from Scripting MS to Tasking MS in response
type ExecutionPayloadResult struct {
	ExecutionID string `json:"executionId"`
	StatusCode  int    `json:"statusCode"`
	StatusText  string `json:"statusText"`
}

// ExpiredExecution describes expired execution data for particular task instance ID which would be sent from Tasking MS to Scripting MS
type ExpiredExecution struct {
	TaskInstanceID     gocql.UUID   `json:"taskInstanceId"`
	ManagedEndpointIDs []gocql.UUID `json:"managedEndpointIds"`
}
