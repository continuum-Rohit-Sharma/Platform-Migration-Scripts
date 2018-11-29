package main

import (
	"fmt"

	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/env"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	"github.com/ContinuumLLC/platform-common-lib/src/services"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
)

type healthCheckDependencyImpl struct {
	services.HealthCheckServiceFactoryImpl
	services.HealthCheckDalFactoryImpl
	services.VersionFactoryImpl
	env.FactoryEnvImpl
	procParser.ParserFactoryImpl
}

func main() {
	h := services.HealthCheckServiceFactoryImpl{}
	s := h.GetHealthCheckService(healthCheckDependencyImpl{})
	model.StrartTime = time.Now()
	health, _ := s.GetHealthCheck(model.HealthCheck{
		Version: model.Version{
			SolutionName:    "SolutionName",
			ServiceName:     "ServiceName",
			ServiceProvider: "ContinuumLLC",
			Major:           "1",
			Minor:           "1",
			Patch:           "11",
		},
		ListenPort: ":8081",
	})

	fmt.Println(health)
}
