package sync

import (
	"github.com/ContinuumLLC/platform-common-lib/src/web/rest"
)

//Service : sync service is a service which wiil be used for synchronization between Microservices
//Service : sync service is a centralized service for maintaining configuration information, naming,
//providing distributed synchronization, and providing group services.
//All of these kinds of services are used in some form or another by distributed applications.
type Service interface {
	Send(path string, data string) error
	Listen(path string, c chan Response) error
	Health() rest.Statuser
}

//Config : Config is a struct to keep all the configuration about the Sync Service, needed while connecting to Servers
type Config struct {
	Servers                []string
	SessionTimeoutInSecond int
}

//Response : is a struct returned on a channel from Listen API; This contains a response or error
type Response struct {
	Data  string
	Error error
}
