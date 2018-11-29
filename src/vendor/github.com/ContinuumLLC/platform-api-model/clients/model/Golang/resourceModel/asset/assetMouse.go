package asset

//Mouse is the definition of /resources/asset/assetMouse.json
type Mouse struct {
	Manufacturer    string `json:"manufacturer,omitempty" cql:"manufacturer"`
	Name            string `json:"name,omitempty" cql:"name"`
	DeviceID        string `json:"deviceId,omitempty" cql:"device_id"`
	DeviceInterface int    `json:"deviceInterface" cql:"device_interface"`
	PointingType    int    `json:"pointingType" cql:"pointing_type"`
	Buttons         int    `json:"buttons" cql:"buttons"`
}
