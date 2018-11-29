package useraccount

import (
	"time"
)

//Info is the struct definition of /resources/useraccount/info
type Info struct {
	UserID      string                 `json:"userID,omitempty"`
	UserName    string                 `json:"userName,omitempty"`
	DomainName  string                 `json:"domain,omitempty"`
	AccountType string                 `json:"accountType,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt   time.Time              `json:"createdAt,omitempty"`
	UpdatedAt   time.Time              `json:"updatedAt,omitempty"`
}
