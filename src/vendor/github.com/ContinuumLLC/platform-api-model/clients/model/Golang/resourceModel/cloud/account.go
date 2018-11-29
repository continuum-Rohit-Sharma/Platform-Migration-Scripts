package cloud

//Account represents a Cloud Account
type Account struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	ClientID string `json:"clientid"`
	Status string `json:"status"`
}

//Accounts represents multiple cloud Accounts
type Accounts struct {
	Accounts []Account `json:"accounts"`
}
