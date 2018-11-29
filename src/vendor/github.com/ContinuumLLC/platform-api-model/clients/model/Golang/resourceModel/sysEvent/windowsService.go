package sysEvent

type Service struct {
	Name                    string `json:"serviceName,omitempty" cql:"service_name"`
	DisplayName             string `json:"displayName,omitempty" cql:"display_name"`
	PathName                string `json:"executablePath,omitempty" cql:"executable_path"`
	StartMode               string `json:"startupType,omitempty" cql:"startup_type"`
	State                   string `json:"serviceStatus,omitempty" cql:"service_status"`
	LogOnAs                 string `json:"logOnAs,omitempty" cql:"logon_as"`
	StopEnableAction        bool   `json:"stopEnableAction" cql:"stop_enable_action"`
	DelayedAutoStart        bool   `json:"delayedAutoStart" cql:"delayed_auto_start"`
	ExitCode                uint32 `json:"win32ExitCode" cql:"win32_exit_code"`
	ServiceSpecificExitCode uint32 `json:"serviceSpecificExitCode" cql:"service_specific_exit_code"`
}
