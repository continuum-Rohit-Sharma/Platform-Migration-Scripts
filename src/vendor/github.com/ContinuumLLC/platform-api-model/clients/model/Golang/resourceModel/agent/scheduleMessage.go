package agent

import "time"

//ScheduleTask is a struct defining the actual schedule task
type ScheduleTask struct {
	Task             string `json:"task"`
	TaskInput        string `json:"taskInput"`
	ExecuteNow       string `json:"executeNow"`
	ExecuteOnStartup bool   `json:"executeOnStartup"`
	Schedule         string `json:"schedule"`
	TimeoutInSeconds int    `json:"timeout"`
}

//ScheduleMessage is the struct definition of /resources/agent/scheduleMessage
type ScheduleMessage struct {
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Version      string    `json:"version"`
	TimestampUTC time.Time `json:"timestampUTC"`
	Path         string    `json:"path"`
	ScheduleTask
}
