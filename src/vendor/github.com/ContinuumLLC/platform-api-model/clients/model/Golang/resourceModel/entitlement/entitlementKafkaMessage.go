package entitlement

// Action constants
const (
	CreateAction = "CREATE"
	DeleteAction = "DELETE"
)

// FeatureChange is the struct definition of /resources/entitlement/entitlementMessege
type FeatureChange struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Action       string                 `json:"action"`
	Entitlements []EndpointRelationship `json:"entitlements"`
}
