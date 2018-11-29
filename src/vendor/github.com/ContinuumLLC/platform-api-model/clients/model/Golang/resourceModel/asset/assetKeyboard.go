package asset

//Keyboard is the struct definition of /resources/asset/assetKeyboard
type Keyboard struct {
	DeviceID    string `json:"deviceID,omitempty" cql:"device_id"`
	Name        string `json:"name,omitempty" cql:"name"`
	Description string `json:"description,omitempty" cql:"description"`
}
