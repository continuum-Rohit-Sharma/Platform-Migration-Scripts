package webroot

import "time"

// WebrootPluginStatusResult represents status result details returned from collected status data
type WebrootPluginStatusResult struct {
	TimestampUTC time.Time     `json:"timestampUTC" description:"UTC time when the Script execution finished"`
	Status       WebrootStatus `json:"status"       desciption:"Status collected by WR Plugin"`
	MessageType  string        `json:"messageType"  description:"Type of result message, can be status or command"`
}

// WebrootStatus represents collected information about WR agent
type WebrootStatus struct {
	Installed         bool      `json:"installed"         description:"Indicates whether WR Agent installed on the endpoint"`
	Running           bool      `json:"running"           description:"Indicates whether WR Agent is running on the endpoint"`
	AttentionRequired bool      `json:"AttentionRequired" description:"Infected and requires attention"`
	AgentVersion      string    `json:"agentVersion"      description:"Version of the WR Agent"`
	AgentMachineID    string    `json:"agentMachineID"    description:"WR Machine ID of the endpoint"`
	ActiveThreats     uint64    `json:"activeThreats"     description:"Count of the active threats found by WR on the endpoint"`
	Threats           []Threat  `json:"Threats"           description:"Threat records found in registry Threats\history"`
	LastScan          time.Time `json:"LastScan"          description:"Date of last scan"`
	ActiveScans       int       `json:"ActiveScans"       description:"Number of scans happening at this moment"`
	CurrentlyCleaning bool      `json:"CurrentlyCleaning" description:"Defines if any cleanings is running at this moment"`
	Infected          bool      `json:"Infected"          description:"Is the system infected or not"`
	UpdateTime        time.Time `json:"UpdateTime"        description:"The timestamp of the latest status update, as seen by the local agent"`
}

// Threat represents information about webroot threat
type Threat struct {
	PathName      string    `json:"PathName"        description:"Path to the threat file"`
	FileName      string    `json:"FileName"        description:"Name of the file that contains threat"`
	InfectionName string    `json:"InfectionName"   description:"The name of the virus"`
	FirstSeen     time.Time `json:"FirstSeen"       description:"The date when threat was discovered"`
	FileSize      int64     `json:"FileSize"        description:"The size of threat file in bytes"`
	FileMD5       string    `json:"FileMD5"         description:"MD5 hash of threat file"`
}
