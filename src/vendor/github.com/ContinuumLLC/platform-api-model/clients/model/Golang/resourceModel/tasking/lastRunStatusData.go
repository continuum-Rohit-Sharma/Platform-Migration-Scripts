package tasking

// LastRunStatusData stores info about amount of devices and its succeeded or failed execution results for task
type LastRunStatusData struct {
	Status       string `json:"status"`
	DeviceCount  int    `json:"deviceCount"`  // Total count of the devices the task was run on
	SuccessCount int    `json:"successCount"` // Total count of succeeded results among all devices
	FailureCount int    `json:"failureCount"` // Total count of failed results among all devices
}
