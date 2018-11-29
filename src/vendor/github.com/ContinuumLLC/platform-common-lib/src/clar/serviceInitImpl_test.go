package clar

import (
	"strings"
	"testing"
)

func TestGetServiceInit(t *testing.T) {
	service := ServiceInitFactoryImpl{}.GetServiceInit()
	_, ok := service.(*serviceInit)
	if !ok {
		t.Error("serviceInit is not serviceInit")
	}
}

func TestSetupOsArgs(t *testing.T) {
	service := ServiceInitFactoryImpl{}.GetServiceInit()
	args := []string{"performanceService", "config.json", ""} //Log file name passed as blank so that logger does not create a log file
	service.SetupOsArgs("", "", args, 1, 2)
	if strings.Compare(service.GetConfigPath(), "config.json") != 0 {
		t.Error("Config Path is not Equal to config.json")
	}
}
