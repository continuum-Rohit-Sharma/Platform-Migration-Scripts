package vmware

//GuestInstance is the struct definition of Guest OS info
type GuestInstance struct {
	ID                string             `json:"id,omitempty"`
	Name              string             `json:"name,omitempty"`
	Os                Os                 `json:"os,omitempty"`
	State             string             `json:"state,omitempty"`
	RmmAgentInstalled bool               `json:"rmmAgent,omitempty"`
	IPAddress         string             `json:"ipAddress,omitempty"`
	Drives            []VirtualDataStore `json:"drives,omitempty"`
}
