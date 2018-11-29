package asset

import "time"

//AssetOs is the struct definition of /resources/asset/assetOs
type AssetOs struct {
	Product          string    `json:"product" cql:"product"`
	Manufacturer     string    `json:"manufacturer,omitempty" cql:"manufacturer"`
	OsLanguage       string    `json:"osLanguage,omitempty" cql:"os_language"`
	SerialNumber     string    `json:"serialNumber,omitempty" cql:"serial_number"`
	Version          string    `json:"version" cql:"version"`
	InstallDate      time.Time `json:"installDate,omitempty" cql:"installdate"`
	Type             string    `json:"type,omitempty" cql:"type"`
	Arch             string    `json:"arch,omitempty" cql:"arch"`
	ServicePack      string    `json:"servicePack" cql:"service_pack"`
	OsInstalledDrive string    `json:"osInstalledDrive" cql:"os_installed_drive"`
	BuildNumber      string    `json:"buildNumber" cql:"build_number"`
	ProductID        string    `json:"productID,omitempty" cql:"product_id"`
	ProductKey       string    `json:"productKey,omitempty" cql:"product_key"`
}
