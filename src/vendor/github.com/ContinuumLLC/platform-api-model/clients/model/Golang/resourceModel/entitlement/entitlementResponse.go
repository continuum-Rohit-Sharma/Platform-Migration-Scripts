package entitlement

import (
	"time"
)

// Feature is the struct definition of /resources/entitlement/entitlementMessege
type Feature struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
	CreatedBy    string                 `json:"createdBy"`
	UpdatedBy    string                 `json:"updatedBy"`
	Entitlements []EndpointRelationship `json:"entitlements"`
}

// EndpointRelationship is the struct definition of /resources/entitlement/entitlementDefinition
type EndpointRelationship struct {
	PartnerID  string `json:"partnerID"`
	ClientID   string `json:"clientID"`
	SiteID     string `json:"siteID"`
	EndpointID string `json:"endpointID"`
}
