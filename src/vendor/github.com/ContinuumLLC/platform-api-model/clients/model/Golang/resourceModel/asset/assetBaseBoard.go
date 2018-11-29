package asset

import "time"

//AssetBaseBoard is the struct definition of /resources/asset/assetBaseBoard
type AssetBaseBoard struct {
	Product      string    `json:"product,omitempty" cql:"product"`
	Manufacturer string    `json:"manufacturer" cql:"manufacturer"`
	Model        string    `json:"model,omitempty" cql:"model"`
	SerialNumber string    `json:"serialNumber,omitempty" cql:"serial_number"`
	Version      string    `json:"version,omitempty" cql:"version"`
	Name         string    `json:"name,omitempty" cql:"name"`
	InstallDate  time.Time `json:"installDate,omitempty" cql:"install_date"`
}
