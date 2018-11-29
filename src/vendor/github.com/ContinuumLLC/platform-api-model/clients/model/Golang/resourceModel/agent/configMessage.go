package agent

import "time"

//ConfigMessage is the struct definition of /resources/agent/configMessage
type ConfigMessage struct {
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Version          string    `json:"version"`
	TimestampUTC     time.Time `json:"timestampUTC"`
	Path             string    `json:"path"`
	ConfigPath       string    `json:"configPath"`
	UpdateNow        string    `json:"updateNow"`
	FutureUpdateTime string    `json:"futureUpdateTime"`
	ConfigJSONDelta  string    `json:"configJSONDelta"`
}
