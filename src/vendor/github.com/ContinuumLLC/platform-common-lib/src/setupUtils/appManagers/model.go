package appManagers

//Manager interface which every package installer need to implement
type Manager interface {
	Install(config *Config) error
	Uninstall(config *Config) error
}

//GetManager gets installer based on the appname
func GetManager(appName string, appPath string, params []string) Manager {
	return junoAgentManager{
		AppPath: appPath,
		Params:  params,
	}
}

//Config Stores the configuration file values
type Config struct {
	RMMAgentBackupINIFilePath string
	RMMAgentINIFilePath       string
	AppManagerLogFilePath     string
	AgentCorePath             string
	AgentConfigFilePath       string
	AgentLogFilePath          string
	JunoAgentURL              string
	JunoAgentExeLocation      string
	UnsupportedOS             string
}
