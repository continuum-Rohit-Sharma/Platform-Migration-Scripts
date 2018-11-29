package clar

var serviceInitSingleton ServiceInit

// ServiceInitFactoryImpl is the serviceInit Factory
type ServiceInitFactoryImpl struct{}

// GetServiceInit returns an instance of ServiceInit
// This code is not thread safe. the creation of ServiceInit is unlikely to
// happen across multiple goroutines, hence mutex is not used.
func (ServiceInitFactoryImpl) GetServiceInit() ServiceInit {
	if serviceInitSingleton == nil {
		serviceInitSingleton = newServiceInit()
	}
	return serviceInitSingleton
}

func newServiceInit() ServiceInit {
	return &serviceInit{}
}

type serviceInit struct {
	configFilePath string
	configIndex    int
	logFilePath    string
	logIndex       int
}

func (s *serviceInit) GetConfigPath() string {
	return s.configFilePath
}

func (s *serviceInit) GetLogFilePath() string {
	return s.logFilePath
}

func (s *serviceInit) SetupOsArgs(defaultConfig, defaultLog string, args []string, configIdex, logIndex int) {
	s.configFilePath = defaultConfig
	s.logFilePath = defaultLog
	s.configIndex = configIdex
	s.logIndex = logIndex
	s.setupConfigFile(args)
	s.setupLogFile(args)
}

func (s *serviceInit) setupConfigFile(args []string) {
	if len(args) > s.configIndex {
		value := args[s.configIndex]
		if value != "" {
			s.configFilePath = value
		}
	}
}

// irrespective of whether log file parameter is present or not,
// the log writer should be setup final value of the logFilePath
func (s *serviceInit) setupLogFile(args []string) {
	if len(args) > s.logIndex {
		value := args[s.logIndex]
		if value != "" {
			s.logFilePath = value
		}
	}
}
