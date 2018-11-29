package asset

//AssetRaidController is the struct definition of /resources/asset/assetRaidController
type AssetRaidController struct {
	SoftwareRaid string `json:"softwareRaid,omitempty" cql:"software_raid"`
	HardwareRaid string `json:"hardwareRaid,omitempty" cql:"hardware_raid"`
	Vendor       string `json:"vendor,omitempty" cql:"vendor"`
}
