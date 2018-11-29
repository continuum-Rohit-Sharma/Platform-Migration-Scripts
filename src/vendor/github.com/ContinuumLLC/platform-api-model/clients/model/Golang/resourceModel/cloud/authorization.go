package cloud

import "time"

//Authorization stores the values used for authorization of Partner
type Authorization struct {
	TimeStampUTC  time.Time
	PartnerID     string
	SiteID        string
	CloudPlatform string
	AccessToken   string
	TokenType     string
	Expiration    string
	IDToken       string
	RefreshToken  string
}
