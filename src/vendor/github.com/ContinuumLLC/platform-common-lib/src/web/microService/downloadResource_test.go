package microservice

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
	"github.com/ContinuumLLC/platform-common-lib/src/web"
	cwmock "github.com/ContinuumLLC/platform-common-lib/src/web/mock"
	"github.com/golang/mock/gomock"
)

func TestDownloadGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRsp := httptest.NewRecorder()
	mockRc := cwmock.NewMockRequestContext(ctrl)
	mockRc.EXPECT().GetResponse().Return(mockRsp).AnyTimes()

	downloadResource{
		log: logging.GetLoggerFactory().Get(),
		fileInfo: model.DownloadFileInfo{
			Name:        "Swagger.yaml",
			Path:        "abcd",
			ContentType: "application/octet-stream",
		},
	}.Get(mockRc)

	if mockRsp.Code == http.StatusOK {
		t.Errorf("Unexpected error code : %v", mockRsp.Code)
	}
}

func TestCreateDownloadRouteConfig(t *testing.T) {
	fileInfo := model.DownloadFileInfo{
		Name:        "Swagger.yaml",
		Path:        "abcd",
		ContentType: "application/octet-stream",
	}

	route := web.RouteConfig{
		URLPathSuffix: "/download",
		URLVars:       []string{},
		Res: downloadResource{
			fileInfo: fileInfo,
			log:      logging.GetLoggerFactory().Get(),
		},
	}

	rout := CreateDownloadRouteConfig(fileInfo, "/download")

	if reflect.DeepEqual(route, rout) {
		t.Errorf("Expected same but got Different %v : %v", route, rout)
	}
}
