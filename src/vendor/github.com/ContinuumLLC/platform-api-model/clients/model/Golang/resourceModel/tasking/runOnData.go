package tasking

// RunOnData describes device type and count of particular task
type RunOnData struct {
	TargetCount int    `json:"count"`
	TargetType  string `json:"type"`
}
