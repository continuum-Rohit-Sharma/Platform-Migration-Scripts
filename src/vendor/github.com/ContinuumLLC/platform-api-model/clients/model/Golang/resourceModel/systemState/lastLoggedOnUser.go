package systemState

import "time"

//LastLoggedOnUser is the struct definition of /resources/systemState/LastLoggedOnUser
type LastLoggedOnUser struct {
	Username  string    `json:"username" cql:"username"`
	LogonTime time.Time `json:"logonTime" cql:"logon_time"`
	Status    string    `json:"status" cql:"status"`
	IsAdmin   bool      `json:"isAdmin" cql:"is_admin"`
}
