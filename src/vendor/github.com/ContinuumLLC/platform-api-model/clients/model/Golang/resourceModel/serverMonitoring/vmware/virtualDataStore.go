package vmware

//VirtualDataStore describes Guest OS hard disk
type VirtualDataStore struct {
	ID          string `json:"id,omitempty"`
	CapacityMB  int64  `json:"capacity,omitempty"`
	FreeSpaceMB int64  `json:"freeSpace,omitempty"`
}
