// Package clar is - command line argument reader.
package clar

//ServiceInit interface is used for service initilization
type ServiceInit interface {
	//GetConfigPath returns the config path
	GetConfigPath() string
	//GetLogFilePath returns the log file path
	GetLogFilePath() string
	//SetupOsArgs setup the config & log file path based on command line argument passed. This function is expected to be called only once in main.
	SetupOsArgs(defaultConfig, defaultLog string, args []string, configIdex, logIndex int)
}

//ServiceInitFactory interface gives the instance of the ServiceInit
type ServiceInitFactory interface {
	//GetServiceInit returns the implementation of ServiceInit interface
	GetServiceInit() ServiceInit
}
