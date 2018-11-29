package vmware

//DataStore describes Host OS hard disk
type DataStore struct {
	Name               string `json:"name,omitempty"`
	CapacityMB         int64  `json:"capacityMB,omitempty"`
	FreeSpaceMB        int64  `json:"freeSpaceMB,omitempty"`
	Type               string `json:"type,omitempty"`
	NumberOfPartitions int    `json:"numberOfPartitions,omitempty"`
	MaintenanceMode    string `json:"maintenanceMode,omitempty"`
}
