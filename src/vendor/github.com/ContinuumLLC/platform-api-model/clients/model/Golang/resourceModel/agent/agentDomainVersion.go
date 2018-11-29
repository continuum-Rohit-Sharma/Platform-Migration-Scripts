package agent

import "time"

//AgentDomainVersion is the struct definition of /resources/agent/agentDomainVersion
type AgentDomainVersion struct {
	TimeStampUTC   time.Time `json:"timeStampUTC"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	ServiceName    string    `json:"serviceName"`
	ServiceVersion string    `json:"serviceVersion"`
	Agents         []Agent   `json:"agent"`
}
