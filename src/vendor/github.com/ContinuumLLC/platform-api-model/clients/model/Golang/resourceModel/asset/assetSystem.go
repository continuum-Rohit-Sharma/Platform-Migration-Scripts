package asset

//AssetSystem is the struct definition of /resources/asset/assetSystem
type AssetSystem struct {
	Product             string `json:"product,omitempty" cql:"product"`
	Model               string `json:"model,omitempty" cql:"model"`
	TimeZone            string `json:"timeZone,omitempty" cql:"time_zone"`
	TimeZoneDescription string `json:"timeZoneDescription,omitempty" cql:"time_zone_description"`
	SerialNumber        string `json:"serialNumber,omitempty" cql:"serial_number"`
	SystemName          string `json:"systemName" cql:"system_name"`
	Category            string `json:"category,omitempty" cql:"category"`
	Domain              string `json:"domain,omitempty" cql:"domain"`
	DomainRole          int    `json:"domainRole" cql:"domain_role"`
}
