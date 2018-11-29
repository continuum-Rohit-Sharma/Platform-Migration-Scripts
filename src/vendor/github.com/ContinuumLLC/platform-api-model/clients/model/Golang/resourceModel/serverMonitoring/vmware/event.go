package vmware

import (
	"time"
)

// Event describes Host system event (critical and noncritical)
type Event struct {
	CreatedTimeUTC time.Time `json:"createdTimeUTC,omitempty"`
	Category       string    `json:"category,omitempty"`
	Message        string    `json:"message,omitempty"`
}
