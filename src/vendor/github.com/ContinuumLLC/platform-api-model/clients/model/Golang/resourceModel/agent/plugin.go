package agent

import "time"

//Plugin is the struct definition of /resources/agent/plugin
type Plugin struct {
	TimeStampUTC  time.Time `json:"timeStampUTC"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	PluginVersion string    `json:"pluginVersion"`
}
