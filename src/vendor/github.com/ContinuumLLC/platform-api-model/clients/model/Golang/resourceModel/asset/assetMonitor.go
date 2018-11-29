package asset

//Monitor is the struct definition of /resources/asset/assetMonitor
type Monitor struct {
	Name         string `json:"name,omitempty" cql:"name"`
	DeviceID     string `json:"deviceID,omitempty" cql:"device_id"`
	Manufacturer string `json:"manufacturer,omitempty" cql:"manufacturer"`
	ScreenHeight uint32 `json:"screenHeight" cql:"screen_height"`
	ScreenWidth  uint32 `json:"screenWidth" cql:"screen_width"`
}
