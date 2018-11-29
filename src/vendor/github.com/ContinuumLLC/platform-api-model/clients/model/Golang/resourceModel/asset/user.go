package asset

//User is the struct definition of /resources/asset/assetUser
type User struct {
	UserID                    string `json:"userid,omitempty" cql:"userid"`
	Username                  string `json:"username,omitempty" cql:"username"`
	DomainName                string `json:"domainName,omitempty" cql:"domain_name"`
	UserType                  string `json:"userType,omitempty" cql:"user_type"`
	AccountType               string `json:"accountType,omitempty" cql:"account_type"`
	UserDisabled              bool   `json:"userDisabled" cql:"user_disabled"`
	UserLockout               bool   `json:"userLockout" cql:"user_lockout"`
	PasswordRequired          bool   `json:"passwordRequired" cql:"password_required"`
	PasswordChangeable        bool   `json:"passwordChangeable" cql:"password_changeable"`
	PasswordExpires           bool   `json:"passwordExpires" cql:"password_expires"`
	PasswordComplexityEnabled bool   `json:"passwordComplexityEnabled" cql:"password_complexity_enabled"`
	RemoteDesktopAllowed      bool   `json:"remoteDesktopAllowed" cql:"remote_desktop_allowed"`
	LockoutThreshold          uint32 `json:"lockoutThreshold" cql:"lockout_threshold"`
	LockoutDurationMins       uint32 `json:"lockoutDurationMins" cql:"lockout_duration_mins"`
	LastLogonTimestamp        uint32 `json:"lastLogonTimestamp" cql:"lastlogon_timestamp"`
	RemoteAccessAllowed       bool   `json:"remoteAccessAllowed" cql:"remote_access_allowed"`
	SystemRole                string `json:"systemRole,omitempty" cql:"system_role"`
}
