package asset

//AssetBios is the struct definition of /resources/asset/assetBios
type AssetBios struct {
	Product       string `json:"product,omitempty" cql:"product"`
	Manufacturer  string `json:"manufacturer,omitempty" cql:"manufacturer"`
	Version       string `json:"version,omitempty" cql:"version"`
	SerialNumber  string `json:"serialNumber,omitempty" cql:"serial_number"`
	SmbiosVersion string `json:"smbiosVersion,omitempty" cql:"smbios_version"`
}
