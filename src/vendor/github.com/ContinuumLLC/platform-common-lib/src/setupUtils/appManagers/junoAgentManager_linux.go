package appManagers

type junoAgentManager struct {
	AppPath string
	Params  []string
}

func (r junoAgentManager) Install(config *Config) error {
	return nil
}

func (r junoAgentManager) Uninstall(config *Config) error {
	return nil
}
