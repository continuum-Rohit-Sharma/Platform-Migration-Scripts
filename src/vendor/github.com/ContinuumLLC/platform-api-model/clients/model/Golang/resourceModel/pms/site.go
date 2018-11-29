package pms

// Site represent site response structure for PMS
type Site struct {
	ResType        int      `json:"ResType,omitempty"`
	ResDescription string   `json:"ResDescription,omitempty"`
	Data           SiteData `json:"Data,omitempty"`
}

// SiteData represent site response structure for PMS
// each attribute map to mstsite table column of database
type SiteData struct {
	CallerGreeting    string `json:"CallerGreeting,omitempty"`
	ClientID          string `json:"ClientID,omitempty"`
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
	SiteID            int    `json:"SiteId,omitempty"`
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
