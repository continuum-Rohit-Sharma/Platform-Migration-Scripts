package pms

// Client represent client response structure for PMS
type Client struct {
	ResType        int        `json:"ResType,omitempty"`
	ResDescription string     `json:"ResDescription,omitempty"`
	Data           ClientData `json:"Data,omitempty"`
}

// ClientData represent client response structure for PMS
// each attribute map to mstsite table as in legacy
// system client concept does not exist column of database
type ClientData struct {
	ClientID          int    `json:"ClientID,omitempty"`
	CallerGreeting    string `json:"CallerGreeting,omitempty"`
	ClientMainPhoneNo string `json:"ClientMainPhoneNo,omitempty"`
	CreationDt        string `json:"CreationDt,omitempty"`
	DisabledON        string `json:"DisabledON,omitempty"`
	IsEnabled         bool   `json:"IsEnabled,omitempty"`
	MemberID          int    `json:"MemberId,omitempty"`
	Proxy             bool   `json:"Proxy,omitempty"`
	SiteAddress       string `json:"SiteAddress,omitempty"`
	SiteAddress2      string `json:"SiteAddress2,omitempty"`
	SiteCity          string `json:"SiteCity,omitempty"`
	SiteCountry       int    `json:"SiteCountry,omitempty"`
	SiteName          string `json:"ActivedOn,omitempty"`
	SitePostalCode    string `json:"SitePostalCode,omitempty"`
	SiteState         string `json:"SiteState,omitempty"`
	Sitecode          string `json:"Sitecode,omitempty"`
	Status            int    `json:"Status,omitempty"`
	TimeStamp         string `json:"TimeStamp,omitempty"`
	TimeZone          string `json:"TimeZone,omitempty"`
	UsrName           string `json:"Usr_Name,omitempty"`
	Operation         int    `json:"Operation,omitempty"`
}
