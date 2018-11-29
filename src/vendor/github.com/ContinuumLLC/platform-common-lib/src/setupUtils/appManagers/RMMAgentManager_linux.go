package appManagers

type rmmAgentManager struct {
	AppPath string
	Params  []string
}

func (rmmAgentManager) Install(config *Config) error {
	return nil
}

func (rmmAgentManager) Uninstall() error {
	return nil
}
