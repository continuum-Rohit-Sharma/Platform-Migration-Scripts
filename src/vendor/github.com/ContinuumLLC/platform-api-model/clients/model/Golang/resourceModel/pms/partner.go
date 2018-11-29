package pms

// Partner represent partner response structure for PMS
type Partner struct {
	ResType        int         `json:"ResType,omitempty"`
	ResDescription string      `json:"ResDescription,omitempty"`
	Data           PartnerData `json:"Data,omitempty"`
}

// PartnerData represent partner response structure for PMS
// each attribute map to mstmember table column of database
type PartnerData struct {
	ActivedOn          string `json:"ActivedOn,omitempty"`
	Address            string `json:"Address,omitempty"`
	City               string `json:"City,omitempty"`
	Country            int    `json:"Country,omitempty"`
	DisabledOn         string `json:"DisabledOn,omitempty"`
	DtRegistr          string `json:"DtRegistr,omitempty"`
	EmailID            string `json:"EmailId,omitempty"`
	FreezeEffectedOn   string `json:"FreezeEffectedOn,omitempty"`
	Freezereason       string `json:"Freezereason,omitempty"`
	HDRegion           string `json:"HD_Region,omitempty"`
	ISFreezed          bool   `json:"ISFreezed,omitempty"`
	IsActive           bool   `json:"IsActive,omitempty"`
	Logo               string `json:"Logo,omitempty"`
	MemberCode         string `json:"MemberCode,omitempty"`
	MemberID           int    `json:"MemberId,omitempty"`
	MemberName         string `json:"MemberName,omitempty"`
	MobileNo           string `json:"MobileNo,omitempty"`
	SalesForceRefAccID string `json:"SalesForceRefAccID,omitempty"`
	State              string `json:"State,omitempty"`
	Status             int    `json:"Status,omitempty"`
	TelNo              string `json:"TelNo,omitempty"`
	TimeStamp          string `json:"TimeStamp,omitempty"`
	UnFreezeEffectedOn string `json:"UnFreezeEffectedOn,omitempty"`
	ZipCode            string `json:"ZipCode,omitempty"`
	Operation          int    `json:"Operation,omitempty"`
}
