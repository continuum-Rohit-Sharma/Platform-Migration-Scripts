package asset

//Service is the struct definition of /resources/asset/Service
type Service struct {
	Name    string `json:"serviceName,omitempty" cql:"service_name"`
	Details ServiceDetails
}

//ServiceDetails is the struct definition /resources/asset/Servicedetails
type ServiceDetails struct {
	DisplayName             string `json:"displayName,omitempty" cql:"display_name"`
	ExecutablePath          string `json:"executablePath,omitempty" cql:"executable_path"`
	StartupType             string `json:"startupType,omitempty" cql:"startup_type"`
	ServiceStatus           string `json:"serviceStatus,omitempty" cql:"service_status"`
	LogOnAs                 string `json:"logOnAs,omitempty" cql:"logon_as"`
	StopEnableAction        bool   `json:"stopEnableAction" cql:"stop_enable_action"`
	DelayedAutoStart        bool   `json:"delayedAutoStart" cql:"delayed_auto_start"`
	Win32ExitCode           uint32 `json:"win32ExitCode" cql:"win32_exit_code"`
	ServiceSpecificExitCode uint32 `json:"serviceSpecificExitCode" cql:"service_specific_exit_code"`
}
