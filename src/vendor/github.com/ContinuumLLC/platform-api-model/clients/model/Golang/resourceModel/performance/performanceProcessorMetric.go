package performance

//PerformanceProcessorMetric ...
type PerformanceProcessorMetric struct {
	NumOfProcesses       int64   `json:"numOfProcesses" cql:"num_of_processes"`
	PercentIOTime        float64 `json:"percentIOTime" cql:"percent_io_time"`
	PercentIdleTime      float64 `json:"percentIdleTime" cql:"percent_idle_time"`
	PercentInterruptTime float64 `json:"percentInterruptTime" cql:"percent_interrupt_time"`
	PercentSystemTime    float64 `json:"percentSystemTime" cql:"percent_system_time"`
	PercentUserTime      float64 `json:"percentUserTime" cql:"percent_user_time"`
	PercentUtil          float64 `json:"percentUtil" cql:"percent_util"`
	InterruptsPerSec     int64   `json:"interruptsPerSec" cql:"interrupts_per_sec"`
}
