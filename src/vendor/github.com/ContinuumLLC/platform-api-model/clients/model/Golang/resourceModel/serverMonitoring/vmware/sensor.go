package vmware

//Sensor is the struct definition of VMWare Hardware Sensor info
type Sensor struct {
	Name         string `json:"name,omitempty"`
	HealthStatus string `json:"healthStatus,omitempty"`
	Value        string `json:"value,omitempty"`
	Type         string `json:"type,omitempty"`
}
