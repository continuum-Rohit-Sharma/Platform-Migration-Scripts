package useraccount

import "time"

//Collection is the struct definition of /resources/userAccount/collection
type Collection struct {
	CreateTimeUTC time.Time `json:"createTimeUTC,omitempty"`
	CreatedBy     string    `json:"createdBy,omitempty"`
	Name          string    `json:"name,omitempty"`
	Type          string    `json:"type,omitempty"`
	Action        string    `json:"action,omitempty"`
	Topic         string    `json:"topic,omitempty"`
	EndpointID    string    `json:"endpointID,omitempty"`
	PartnerID     string    `json:"partnerID,omitempty"`
	ClientID      string    `json:"clientID,omitempty"`
	SiteID        string    `json:"siteID,omitempty"`
	RegID         string    `json:"regID,omitempty"`
	ActionType    string    `json:"actionType,omitempty"`
	DomainRole    string    `json:"domainRole,omitempty"`
	Users         []Info    `json:"users,omitempty"`
}
