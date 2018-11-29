package performance

import "time"

//PerformanceNetworkInterface is the struct definition of /resources/performance/PerformanceNetworkInterface
type PerformanceNetworkInterface struct {
	CreateTimeUTC          time.Time `json:"createTimeUTC"`
	CreatedBy              string    `json:"createdBy"`
	Index                  int       `json:"index"`
	Name                   string    `json:"name" cql:"interface_name"`
	Type                   string    `json:"type"`
	ReceivedBytes          int64     `json:"receivedBytes" cql:"received_bytes"`
	TransmittedBytes       int64     `json:"transmittedBytes" cql:"transmitted_bytes"`
	TXQueueLength          int64     `json:"txQueueLength" cql:"txqueue_length"`
	RXQueueLength          int64     `json:"rxQueueLength" cql:"rxqueue_length"`
	ReceivedBytesPerSec    int64     `json:"receivedBytesPerSec" cql:"received_bytes_per_sec"`
	TransmittedBytesPerSec int64     `json:"transmittedBytesPerSec" cql:"transmitted_bytes_per_sec"`
	TotalBytesPerSec       int64     `json:"totalBytesPerSec" cql:"total_bytes_per_sec"`
}
