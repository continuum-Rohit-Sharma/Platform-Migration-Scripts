package performance

import "time"

// Generated from /resources/processor.json
type PerformanceProcessor struct {
	Cores                []PerformanceProcessorCore `json:"cores" cql:"cores"`
	CreateTimeUTC        time.Time                  `json:"createTimeUTC"`
	CreatedBy            string                     `json:"createdBy"`
	Index                int                        `json:"index" cql:"processor_index"`
	Name                 string                     `json:"name" cql:"processor_name"`
	NumOfCores           int                        `json:"numOfCores" cql:"no_of_cores"`
	Metric               PerformanceProcessorMetric `json:"metric" cql:"metric"`
	ProcessorQueueLength float64                    `json:"processorQueueLength" cql:"processor_queue_length"`
	Type                 string                     `json:"type"`
}
