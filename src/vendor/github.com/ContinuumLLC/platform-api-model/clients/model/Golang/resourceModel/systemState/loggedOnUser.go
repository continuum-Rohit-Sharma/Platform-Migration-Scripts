package systemState

import (
	"time"
)

//LoggedOnUser is the struct definition of /resources/systemState/loggedOnUser
type LoggedOnUser struct {
	Username    string    `json:"username" cql:"username"`
	SessionID   string    `json:"sessionID" cql:"session_id"`
	SessionName string    `json:"sessionName" cql:"session_name"`
	Status      string    `json:"status" cql:"status"`
	Client      string    `json:"client" cql:"client"`
	IsAdmin     bool      `json:"isAdmin" cql:"is_admin"`
	LogonTime   time.Time `json:"logonTime" cql:"logon_time"`
}
