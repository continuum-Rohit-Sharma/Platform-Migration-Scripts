package systemState

import "time"

//StartupStatus is the struct definition of /resources/systemState/startupStatus
type StartupStatus struct {
	LastBootUpTimeUTC time.Time `json:"lastBootUpTimeUTC" cql:"last_bootuptime_utc"`
}
