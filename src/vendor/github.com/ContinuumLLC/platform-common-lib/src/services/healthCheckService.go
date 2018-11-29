package services

import (
	aModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/healthCheck"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
)

//HealthCheckServiceFactoryImpl returns the concrete implementation of Factory
type HealthCheckServiceFactoryImpl struct {
}

//GetHealthCheckService : A factory function to create an instance of HealthCheck Service
func (HealthCheckServiceFactoryImpl) GetHealthCheckService(f model.HealthCheckDependencies) model.HealthCheckService {
	return healthCheckServiceImpl{
		f:   f,
		log: logging.GetLoggerFactory().Get(),
	}
}

//healthCheckServiceImpl returns the concrete implementation of HealthCheckService
type healthCheckServiceImpl struct {
	f   model.HealthCheckDependencies
	log logging.Logger
}

func (h healthCheckServiceImpl) GetHealthCheck(healthCheck model.HealthCheck) (aModel.HealthCheck, error) {
	h.log.LogWithCorrelationf(logging.DEBUG, healthCheck.Version.ServiceName, "Retriving Health Information for %v", healthCheck)
	dal := h.f.GetHealthCheckDal(h.f)
	return dal.GetHealthCheck(healthCheck)
}
