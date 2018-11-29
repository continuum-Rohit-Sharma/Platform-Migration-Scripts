package microservice

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"reflect"

	aModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/healthCheck"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/services/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
	"github.com/ContinuumLLC/platform-common-lib/src/web"
	cwmock "github.com/ContinuumLLC/platform-common-lib/src/web/mock"
	"github.com/golang/mock/gomock"
)

func TestHealtCheckGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listerMock := mock.NewMockHealthCheckDependencies(ctrl)
	service := mock.NewMockHealthCheckService(ctrl)
	listerMock.EXPECT().GetHealthCheckService(gomock.Any()).Return(service)
	service.EXPECT().GetHealthCheck(gomock.Any()).Return(aModel.HealthCheck{}, nil)
	mockRsp := httptest.NewRecorder()
	mockRc := cwmock.NewMockRequestContext(ctrl)
	mockRc.EXPECT().GetResponse().Return(mockRsp).AnyTimes()

	healthCheckResource{
		f:   listerMock,
		log: logging.GetLoggerFactory().Get(),
		healthCheck: model.HealthCheck{
			Version: model.Version{
				SolutionName:    "SolutionName",
				ServiceName:     "ServiceName",
				ServiceProvider: "ContinuumLLC",
				Major:           "1",
				Minor:           "1",
				Patch:           "11",
			},
			ListenPort: ":8081",
		},
	}.Get(mockRc)

	if mockRsp.Code != http.StatusOK {
		t.Errorf("Unexpected error code : %v", mockRsp.Code)
	}
}

func TestHealtCheckGetErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listerMock := mock.NewMockHealthCheckDependencies(ctrl)
	service := mock.NewMockHealthCheckService(ctrl)
	listerMock.EXPECT().GetHealthCheckService(gomock.Any()).Return(service)
	service.EXPECT().GetHealthCheck(gomock.Any()).Return(aModel.HealthCheck{}, errors.New("Error"))
	mockRsp := httptest.NewRecorder()
	mockRc := cwmock.NewMockRequestContext(ctrl)
	mockRc.EXPECT().GetResponse().Return(mockRsp).AnyTimes()

	healthCheckResource{
		f:   listerMock,
		log: logging.GetLoggerFactory().Get(),
		healthCheck: model.HealthCheck{
			Version: model.Version{
				SolutionName:    "SolutionName",
				ServiceName:     "ServiceName",
				ServiceProvider: "ContinuumLLC",
				Major:           "1",
				Minor:           "1",
				Patch:           "11",
			},
			ListenPort: ":8081",
		},
	}.Get(mockRc)

	if mockRsp.Code == http.StatusOK {
		t.Errorf("Unexpected error code : %v", mockRsp.Code)
	}
}

func TestCreateHealthCheckRouteConfig(t *testing.T) {
	healthCheck := model.HealthCheck{
		Version: model.Version{
			SolutionName:    "SolutionName",
			ServiceName:     "ServiceName",
			ServiceProvider: "ContinuumLLC",
			Major:           "1",
			Minor:           "1",
			Patch:           "11",
		},
		ListenPort: ":8081",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	listerMock := mock.NewMockHealthCheckDependencies(ctrl)
	route := web.RouteConfig{
		URLPathSuffix: "/healthCheck",
		URLVars:       []string{},
		Res: healthCheckResource{
			healthCheck: healthCheck,
			f:           listerMock,
			log:         logging.GetLoggerFactory().Get(),
		},
	}

	rout := CreateHealthCheckRouteConfig(healthCheck, listerMock)

	if reflect.DeepEqual(route, rout) {
		t.Errorf("Expected same but got Different %v : %v", route, rout)
	}
}
