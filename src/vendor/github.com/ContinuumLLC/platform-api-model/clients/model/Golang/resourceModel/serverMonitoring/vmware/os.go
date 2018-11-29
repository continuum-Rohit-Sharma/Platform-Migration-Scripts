package vmware

import (
	"time"
)

// Os describes ESX Host OS info
type Os struct {
	Product      string     `json:"product,omitempty"`
	SerialNumber string     `json:"serialNumber,omitempty"`
	Version      string     `json:"version,omitempty"`
	InstallDate  *time.Time `json:"installDate,omitempty"`
	BuildNumber  string     `json:"buildNumber,omitempty"`
	ProductKey   string     `json:"productKey,omitempty"`
}
