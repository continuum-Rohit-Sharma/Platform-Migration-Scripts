package services

import (
	"testing"
	"time"

	"reflect"

	aModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/version"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
)

func TestGetVersionService(t *testing.T) {
	srv := VersionFactoryImpl{}.GetVersionService()
	_, ok := srv.(versionServiceImpl)
	if !ok {
		t.Error("Version Service is not IMPL")
	}
}

func TestGetVersion(t *testing.T) {
	srv := VersionFactoryImpl{}.GetVersionService()

	ver := model.Version{
		SolutionName:    "SolutionName",
		ServiceName:     "ServiceName",
		ServiceProvider: "ContinuumLLC",
		Major:           "1",
		Minor:           "1",
		Patch:           "11",
	}

	version := srv.GetVersion(ver)

	expectedVersion := aModel.Version{
		Name:            ver.SolutionName,
		Type:            "Version",
		TimeStampUTC:    time.Now(),
		ServiceName:     ver.ServiceName,
		ServiceProvider: ver.ServiceProvider,
		ServiceVersion:  ver.Major + "." + ver.Minor + "." + ver.Patch,
		BuildNumber:     model.BuildVersion,
	}

	if reflect.DeepEqual(expectedVersion, version) {
		t.Errorf("Expected same but got Different %v : %v", expectedVersion, version)
	}
}
