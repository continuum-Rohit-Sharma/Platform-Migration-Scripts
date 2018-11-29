package microservice

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
	cweb "github.com/ContinuumLLC/platform-common-lib/src/web"
)

type downloadResource struct {
	cweb.Post405
	cweb.Put405
	cweb.Delete405
	cweb.Others405
	log      logging.Logger
	fileInfo model.DownloadFileInfo
}

func (res downloadResource) Get(rc cweb.RequestContext) {
	res.log.LogWithCorrelationf(logging.DEBUG, res.fileInfo.Name, "Get File response for Info %v and request %v", res.fileInfo, rc)
	w := rc.GetResponse()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Disposition", "attachment; filename="+res.fileInfo.Name)
	w.Header().Set("Content-Type", res.fileInfo.ContentType)
	f, err := os.Open(res.fileInfo.Path)
	if err != nil {
		res.log.LogWithCorrelationf(logging.ERROR, res.fileInfo.Name, "%v Error while reading File at %s Path", err, res.fileInfo.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	count, err := io.Copy(w, f)
	if err != nil {
		res.log.LogWithCorrelationf(logging.ERROR, res.fileInfo.Name, "%v Error while writing File at %s Path", err, res.fileInfo.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Length", strconv.FormatInt(count, 10))
}

//CreateDownloadRouteConfig is a function for creating route config for Download url
func CreateDownloadRouteConfig(fileInfo model.DownloadFileInfo, pathSuffix string) *cweb.RouteConfig {
	return &cweb.RouteConfig{
		URLPathSuffix: pathSuffix,
		URLVars:       []string{},
		Res: downloadResource{
			fileInfo: fileInfo,
			log:      logging.GetLoggerFactory().Get(),
		},
	}
}
