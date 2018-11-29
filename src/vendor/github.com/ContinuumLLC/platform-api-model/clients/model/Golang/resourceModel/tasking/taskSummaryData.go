package tasking

import "time"

// TaskSummaryData represents task summary data in TasksAndSequences page
type TaskSummaryData struct {
	TaskID        string            `json:"taskID"`
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	CreatedAt     time.Time         `json:"createdAt"`
	RunOn         RunOnData         `json:"runOn"`
	Regularity    string            `json:"regularity"`
	InitiatedBy   string            `json:"initiatedBy"`
	Status        string            `json:"status"`
	LastRunTime   time.Time         `json:"lastRunTime"`
	LastRunStatus LastRunStatusData `json:"lastRunStatus"`
}
