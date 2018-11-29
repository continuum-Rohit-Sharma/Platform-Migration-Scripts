package model

import aModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/version"

//Version is a common struct to create an common API struct instance for any Service
type Version struct {
	SolutionName    string
	ServiceName     string
	ServiceProvider string
	Major           string
	Minor           string
	Patch           string
}

//BuildVersion for any Service
var BuildVersion = "v1"

//VersionFactory : A factory to create an instance of Version Service
type VersionFactory interface {
	GetVersionService() VersionService
}

//VersionService : A service to create API model version object, so that service can return this
type VersionService interface {
	GetVersion(version Version) aModel.Version
}

//VersionDependencies : A dependencies for Version service and factory
type VersionDependencies interface {
	VersionFactory
}
