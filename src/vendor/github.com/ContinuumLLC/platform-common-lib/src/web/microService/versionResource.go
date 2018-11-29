package microservice

import (
	"encoding/json"
	"net/http"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
	cweb "github.com/ContinuumLLC/platform-common-lib/src/web"
)

type versionResource struct {
	cweb.Post405
	cweb.Put405
	cweb.Delete405
	cweb.Others405
	log     logging.Logger
	version model.Version
	f       model.VersionDependencies
}

func (res versionResource) Get(rc cweb.RequestContext) {
	res.log.LogWithCorrelationf(logging.DEBUG, res.version.ServiceName, "Get Version response for version %v and request %v", res.version, rc)
	w := rc.GetResponse()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")
	service := res.f.GetVersionService()
	response := service.GetVersion(res.version)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		res.log.LogWithCorrelationf(logging.ERROR, res.version.ServiceName, "%v Error while Encoding Version response %v", err, response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

//CreateVersionRouteConfig is a function for creating route config for Version url
func CreateVersionRouteConfig(version model.Version, f model.VersionDependencies) *cweb.RouteConfig {
	return &cweb.RouteConfig{
		URLPathSuffix: "/version",
		URLVars:       []string{},
		Res: versionResource{
			version: version,
			f:       f,
			log:     logging.GetLoggerFactory().Get(),
		},
	}
}
