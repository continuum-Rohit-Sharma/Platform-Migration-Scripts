package asset

import "time"

//AssetCollection is the struct definition of /resources/asset/assetCollection
type AssetCollection struct {
	CreateTimeUTC      time.Time           `json:"createTimeUTC,omitempty"`
	CreatedBy          string              `json:"createdBy,omitempty"`
	Name               string              `json:"name,omitempty"`
	Type               string              `json:"type,omitempty"`
	EndpointID         string              `json:"endpointID,omitempty"`
	PartnerID          string              `json:"partnerID,omitempty"`
	ClientID           string              `json:"clientID,omitempty"`
	SiteID             string              `json:"siteID,omitempty"`
	RegID              string              `json:"regID,omitempty"`
	FriendlyName       string              `json:"friendlyName"`
	RemoteAddress      string              `json:"remoteAddress,omitempty"`
	ResourceType       string              `json:"resourceType,omitempty"`
	EndpointType       string              `json:"endpointType,omitempty"`
	BaseBoard          AssetBaseBoard      `json:"baseBoard,omitempty"`
	Bios               AssetBios           `json:"bios,omitempty"`
	Drives             []AssetDrive        `json:"drives,omitempty"`
	Memory             []PhysicalMemory    `json:"physicalMemory,omitempty"`
	Networks           []AssetNetwork      `json:"networks,omitempty"`
	Os                 AssetOs             `json:"os,omitempty"`
	Processors         []AssetProcessor    `json:"processors,omitempty"`
	RaidController     AssetRaidController `json:"raidController,omitempty"`
	System             AssetSystem         `json:"system,omitempty"`
	InstalledSoftwares []InstalledSoftware `json:"installedSoftwares,omitempty"`
	Keyboards          []Keyboard          `json:"keyboards,omitempty"`
	Mouse              []Mouse             `json:"mouse,omitempty"`
	Monitors           []Monitor           `json:"monitors,omitempty"`
	PhysicalDrives     []PhysicalDrive     `json:"physicalDrives,omitempty"`
	Users              []User              `json:"users,omitempty"`
	Services           []Service           `json:"services,omitempty"`
	Shares             []Share             `json:"shares,omitempty"`
	SoftwareLicenses   []SoftwareLicenses  `json:"softwareLicenses,omitempty"`
}
