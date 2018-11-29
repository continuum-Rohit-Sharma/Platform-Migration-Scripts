package microservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"reflect"

	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/services"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
	"github.com/ContinuumLLC/platform-common-lib/src/web"
	cwmock "github.com/ContinuumLLC/platform-common-lib/src/web/mock"
	"github.com/golang/mock/gomock"
)

type MockVersionDependenciesImpl struct {
	services.VersionFactoryImpl
}

func TestVersionGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listerMock := MockVersionDependenciesImpl{}
	mockRsp := httptest.NewRecorder()
	mockRc := cwmock.NewMockRequestContext(ctrl)
	mockRc.EXPECT().GetResponse().Return(mockRsp).AnyTimes()

	versionResource{
		f:   listerMock,
		log: logging.GetLoggerFactory().Get(),
		version: model.Version{
			SolutionName:    "SolutionName",
			ServiceName:     "ServiceName",
			ServiceProvider: "ContinuumLLC",
			Major:           "1",
			Minor:           "1",
			Patch:           "11",
		},
	}.Get(mockRc)

	if mockRsp.Code != http.StatusOK {
		t.Errorf("Unexpected error code : %v", mockRsp.Code)
	}
}

func TestCreateVersionRouteConfig(t *testing.T) {
	version := model.Version{
		SolutionName:    "SolutionName",
		ServiceName:     "ServiceName",
		ServiceProvider: "ContinuumLLC",
		Major:           "1",
		Minor:           "1",
		Patch:           "11",
	}

	listerMock := MockVersionDependenciesImpl{}
	route := web.RouteConfig{
		URLPathSuffix: "/version",
		URLVars:       []string{},
		Res: versionResource{
			version: version,
			f:       listerMock,
			log:     logging.GetLoggerFactory().Get(),
		},
	}

	rout := CreateVersionRouteConfig(version, listerMock)

	if reflect.DeepEqual(route, rout) {
		t.Errorf("Expected same but got Different %v : %v", route, rout)
	}
}
