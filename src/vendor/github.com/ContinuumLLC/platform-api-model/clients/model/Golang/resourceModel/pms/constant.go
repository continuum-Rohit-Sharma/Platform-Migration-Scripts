package pms

// operation performed on data
const (
	Create = 2
	Update = 4
)

// response type and description for cdc pms
const (
	PartnerTypeCDC = 101
	PartnerDescCDC = "Partner response for cdc"
	SiteTypeCDC    = 102
	SiteDescCDC    = "Site response for cdc"
	ClientTypeCDC  = 103
	ClientDescCDC  = "Client response for cdc"
)

// response type and description for one time migration pms
const (
	OneTimeMigrationPartnerType = 1001
	OneTimeMigrationPartnerDesc = "Partner response for one time migration"
	OneTimeMigrationSiteType    = 1002
	OneTimeMigrationSiteDesc    = "Site response for one time migration"
	OneTimeMigrationClientType  = 1003
	OneTimeMigrationClientDesc  = "Client response for one time migration"
)
