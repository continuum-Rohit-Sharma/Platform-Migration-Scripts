package config

//Configuration is a struct to hold values for config file Update
type Configuration struct {
	FilePath      string
	Content       string
	TransationID  string
	PartialUpdate bool
}

//UpdatedConfig is a struct to hold updated values for config file Update
type UpdatedConfig struct {
	Key      string
	Existing interface{}
	Updated  interface{}
}

//ConfigurationService is a service to Merge Configuration Recieved and update File
type ConfigurationService interface {
	Update(cfg Configuration) ([]UpdatedConfig, error)
}
