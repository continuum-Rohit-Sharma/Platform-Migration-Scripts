package microservice

import (
	"encoding/json"
	"net/http"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
	cweb "github.com/ContinuumLLC/platform-common-lib/src/web"
)

type healthCheckResource struct {
	cweb.Post405
	cweb.Put405
	cweb.Delete405
	cweb.Others405
	log         logging.Logger
	healthCheck model.HealthCheck
	f           model.HealthCheckDependencies
}

func (res healthCheckResource) Get(rc cweb.RequestContext) {
	res.log.LogWithCorrelationf(logging.DEBUG, res.healthCheck.Version.ServiceName, "Get HealthCheck response for HealthCheck %v and request %v", res.healthCheck, rc)
	w := rc.GetResponse()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	service := res.f.GetHealthCheckService(res.f)
	response, err := service.GetHealthCheck(res.healthCheck)
	if err != nil {
		res.log.LogWithCorrelationf(logging.ERROR, res.healthCheck.Version.ServiceName, "%v Error while getting HealthCheck response for %v", err, res.healthCheck)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		res.log.LogWithCorrelationf(logging.ERROR, res.healthCheck.Version.ServiceName, "%v Error while Encoding HealthCheck response %v", err, response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

//CreateHealthCheckRouteConfig is a function for creating route config for healthCheck url
func CreateHealthCheckRouteConfig(healthCheck model.HealthCheck, f model.HealthCheckDependencies) *cweb.RouteConfig {
	return &cweb.RouteConfig{
		URLPathSuffix: "/healthCheck",
		URLVars:       []string{},
		Res: healthCheckResource{
			healthCheck: healthCheck,
			f:           f,
			log:         logging.GetLoggerFactory().Get(),
		},
	}
}
