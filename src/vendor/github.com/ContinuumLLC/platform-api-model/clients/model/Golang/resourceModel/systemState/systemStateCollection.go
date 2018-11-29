package systemState

import "time"

//SystemStateCollection is the struct definition of /resources/systemState/systemStateCollection
type SystemStateCollection struct {
	CreateTimeUTC    time.Time        `json:"createTimeUTC"`
	CreatedBy        string           `json:"createdBy"`
	Name             string           `json:"name"`
	Type             string           `json:"type"`
	EndpointID       string           `json:"endpointID"`
	PartnerID        string           `json:"partnerID"`
	ClientID         string           `json:"clientID"`
	SiteID           string           `json:"siteID"`
	StartupStatus    StartupStatus    `json:"startupStatus"`
	LastLoggedOnUser LastLoggedOnUser `json:"lastLoggedOnUser"`
	LoggedOnUsers    []LoggedOnUser   `json:"loggedOnUsers"`
}
