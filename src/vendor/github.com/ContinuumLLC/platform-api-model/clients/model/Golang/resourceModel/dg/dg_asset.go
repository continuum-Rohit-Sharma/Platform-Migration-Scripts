package dg

//DGAsset is the definition of /resources/dynamicGroups/example/managed_endpoint_message_example_asset.json
type DGAsset struct {
	OriginDomain string `json:"originDomain"`
	ID           string `json:"id"`
	Client       string `json:"client"`
	Partner      string `json:"partner"`
	Site         string `json:"site"`
	AssetData    Asset  `json:"asset"`
}

type Asset struct {
	OS                    string                     `json:"os"`
	ServicePack           string                     `json:"service_pack"`
	OSVersion             string                     `json:"os_version"`
	BaseboardManufacturer string                     `json:"baseboard_manufacturer"`
	VirtualMachine        string                     `json:"virtual_machine"`
	InstalledSoftware     []InstalledSoftwareMessage `json:"installed_software"`
	RAM                   uint64                     `json:"ram"`
}
