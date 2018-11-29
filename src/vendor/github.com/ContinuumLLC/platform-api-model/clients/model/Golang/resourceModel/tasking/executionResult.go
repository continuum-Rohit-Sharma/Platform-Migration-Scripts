package tasking

import (
	"github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/agent"
	"time"
)

// ExecutionResult represents result details returned by Origin MS to Tasking MS
// Tasking POST API route /partners/{partnerID}/task-execution-results/task-instances/{taskInstanceID}
type ExecutionResult struct {
	// Possible values: "Success", "Failed", "Some Failures"
	CompletionStatus string    `json:"completionStatus"  valid:"required"`
	EndpointID       string    `json:"endpointId"        valid:"required,uuid"`
	UpdateTime       time.Time `json:"updateTime"        valid:"required"`
	ErrorDetails     string    `json:"errorDetails"      valid:"-"`
	ResultDetails    string    `json:"resultDetails"     valid:"-"`
}

// ExecutionResultKafkaMessage structure represents the Kafka message with Script execution results
type ExecutionResultKafkaMessage struct {
	agent.BrokerEnvelope
	Message ScriptPluginReturnMessage `json:"message"`
}
