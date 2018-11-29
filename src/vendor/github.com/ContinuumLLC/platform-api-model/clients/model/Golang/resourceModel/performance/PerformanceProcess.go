package performance

import "time"

//PerformanceProcess stores data related to process
type PerformanceProcess struct {
	CreateTimeUTC   time.Time `json:"createTimeUTC"`
	CreatedBy       string    `json:"createdBy"`
	Index           int       `json:"index"`
	Name            string    `json:"name"`
	ProcessID       int       `json:"processID"`
	Type            string    `json:"type"`
	PercentCPUUsage float64   `json:"percentCPUUsage"`
	HandleCount     int64     `json:"handleCount"`
	ThreadCount     int64     `json:"threadCount"`
	PrivateBytes    int64     `json:"privateBytes"`
	UserName        string    `json:"userName"`
	DiskReadBytes   int64     `json:"diskReadBytes"`
	DiskWriteBytes  int64     `json:"diskWriteBytes"`
	NetSendBytes    int64     `json:"netSendBytes"`
	NetReceiveBytes int64     `json:"netReceiveBytes"`
}
