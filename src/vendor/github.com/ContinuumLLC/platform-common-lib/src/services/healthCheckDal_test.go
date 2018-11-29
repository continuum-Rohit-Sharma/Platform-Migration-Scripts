package services

import (
	"bytes"
	"io/ioutil"
	"testing"

	"errors"

	aModelHealth "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/healthCheck"
	aModelVersion "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/version"

	envMock "github.com/ContinuumLLC/platform-common-lib/src/env/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	ppMock "github.com/ContinuumLLC/platform-common-lib/src/procParser/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/services/mock"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
	"github.com/golang/mock/gomock"
)

var version = aModelVersion.Version{
	Name:            "SolutionName",
	ServiceName:     "ServiceName",
	ServiceProvider: "ContinuumLLC",
	ServiceVersion:  "1.1.11",
	BuildNumber:     model.BuildVersion,
}

var healthCheck = model.HealthCheck{
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

func getIfConfigResult() map[string]procParser.Line {
	link := procParser.Line{Values: []string{"enp0s3", "Link", "encap", "Ethernet", "HWaddr", "08", "00", "27", "21", "05", "e2"}}
	inet := procParser.Line{Values: []string{"inet", "addr", "192.168.10.135", "Bcast", "192.168.11.255", "Mask", "255.255.252.0"}}
	tcp := procParser.Line{Values: []string{"tcp", "ESTAB", "0", "0", "192.168.10.135", "34442", "192.30.253.124", "443"}}
	rxTX := procParser.Line{Values: []string{"RX", "bytes", "12338890", "(12.3", "MB)", "TX", "bytes", "1299879", "(1.2", "MB)"}}
	m := make(map[string]procParser.Line)
	m["enp0s3"] = procParser.Line{Values: link.Values}
	m["inet"] = procParser.Line{Values: inet.Values}
	m["tcp"] = procParser.Line{Values: tcp.Values}
	m["RX"] = procParser.Line{Values: rxTX.Values}
	return m
}

func getDataIfConfig() *procParser.Data {
	data := getIfConfigResult()
	return &procParser.Data{
		Map:   data,
		Lines: []procParser.Line{data["enp0s3"], data["inet"], data["tcp"], data["RX"]},
	}
}

func getPSResult() map[string]procParser.Line {
	proc1 := procParser.Line{Values: []string{"PID", "%CPU", "NLWP", "%MEM", "S", "START"}}
	proc2 := procParser.Line{Values: []string{"2043", "22.3", "2", "12.4", "S", "13:01"}}
	m := make(map[string]procParser.Line)
	m["PID"] = procParser.Line{Values: proc1.Values}
	m["2043"] = procParser.Line{Values: proc2.Values}
	return m
}

func getDataPS() *procParser.Data {
	data := getPSResult()
	return &procParser.Data{
		Map:   data,
		Lines: []procParser.Line{data["PID"], data["2043"]},
	}
}

func getPSErrorResult() map[string]procParser.Line {
	proc1 := procParser.Line{Values: []string{"PID", "%CPU", "NLWP", "%MEM", "S", "START"}}
	proc2 := procParser.Line{Values: []string{"2043", "22.3a", "2a", "12.4a", "S", "13:01"}}
	m := make(map[string]procParser.Line)
	m["PID"] = procParser.Line{Values: proc1.Values}
	m["2043"] = procParser.Line{Values: proc2.Values}
	return m
}

func getDataErrorPS() *procParser.Data {
	data := getPSErrorResult()
	return &procParser.Data{
		Map:   data,
		Lines: []procParser.Line{data["PID"], data["2043"]},
	}
}

func TestGetHealthCheckDal(t *testing.T) {
	srv := HealthCheckDalFactoryImpl{}.GetHealthCheckDal(nil)
	_, ok := srv.(healthCheckDalImpl)
	if !ok {
		t.Error("healthCheckDalImpl is not IMPL")
	}
}

func TestGetHealthCheckNetworkErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listerMock := mock.NewMockHealthCheckDependencies(ctrl)
	versionService := mock.NewMockVersionService(ctrl)
	listerMock.EXPECT().GetVersionService().Return(versionService)
	versionService.EXPECT().GetVersion(gomock.Any()).Return(version)
	mockEnv := envMock.NewMockEnv(ctrl)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any()).Return(nil, errors.New("Error"))
	listerMock.EXPECT().GetEnv().Return(mockEnv)
	dal := HealthCheckDalFactoryImpl{}.GetHealthCheckDal(listerMock)
	dal.(healthCheckDalImpl).log.SetLogLevel(logging.OFF)
	_, err := dal.GetHealthCheck(healthCheck)
	if err == nil {
		t.Errorf("Expected Err but Got Result")
	}
}

func TestGetHealthCheckProcessErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listerMock := mock.NewMockHealthCheckDependencies(ctrl)
	versionService := mock.NewMockVersionService(ctrl)
	listerMock.EXPECT().GetVersionService().Return(versionService)
	versionService.EXPECT().GetVersion(gomock.Any()).Return(version)
	mockEnv := envMock.NewMockEnv(ctrl)
	byteReader := bytes.NewReader([]byte(""))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any()).Return(reader, nil)

	mockParser := ppMock.NewMockParser(ctrl)
	listerMock.EXPECT().GetParser().Return(mockParser)
	mockParser.EXPECT().Parse(gomock.Any(), gomock.Any()).Return(getDataIfConfig(), nil)

	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("Error"))

	listerMock.EXPECT().GetEnv().Return(mockEnv).Times(2)
	dal := HealthCheckDalFactoryImpl{}.GetHealthCheckDal(listerMock)
	dal.(healthCheckDalImpl).log.SetLogLevel(logging.OFF)
	_, err := dal.GetHealthCheck(healthCheck)
	if err == nil {
		t.Errorf("Expected Err but Got Result")
	}
}

func TestFindProcessInformation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listerMock := mock.NewMockHealthCheckDependencies(ctrl)
	mockEnv := envMock.NewMockEnv(ctrl)
	listerMock.EXPECT().GetEnv().Return(mockEnv)
	byteReader := bytes.NewReader([]byte(""))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, nil)

	mockParser := ppMock.NewMockParser(ctrl)
	listerMock.EXPECT().GetParser().Return(mockParser)
	mockParser.EXPECT().Parse(gomock.Any(), gomock.Any()).Return(getDataPS(), nil)

	dal := healthCheckDalImpl{f: listerMock, log: logging.GetLoggerFactory().New("")}
	dal.log.SetLogLevel(logging.OFF)

	health := &aModelHealth.HealthCheck{}
	err := dal.findProcessInformation(2, health)
	if err != nil {
		t.Errorf("Expected Result but Got Error")
	}
}

func TestFindProcessInformationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	listerMock := mock.NewMockHealthCheckDependencies(ctrl)
	mockEnv := envMock.NewMockEnv(ctrl)
	listerMock.EXPECT().GetEnv().Return(mockEnv)
	byteReader := bytes.NewReader([]byte(""))
	reader := ioutil.NopCloser(byteReader)
	mockEnv.EXPECT().GetCommandReader(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(reader, nil)

	mockParser := ppMock.NewMockParser(ctrl)
	listerMock.EXPECT().GetParser().Return(mockParser)
	mockParser.EXPECT().Parse(gomock.Any(), gomock.Any()).Return(getDataErrorPS(), nil)

	dal := healthCheckDalImpl{f: listerMock, log: logging.GetLoggerFactory().New("")}
	dal.log.SetLogLevel(logging.OFF)

	health := &aModelHealth.HealthCheck{}
	err := dal.findProcessInformation(2, health)
	if err != nil {
		t.Errorf("Expected Err but Got Result")
	}
}

func TestGetStatus(t *testing.T) {
	dal := healthCheckDalImpl{log: logging.GetLoggerFactory().New("")}
	dal.log.SetLogLevel(logging.OFF)

	status := map[string]string{
		"R":   "Running",
		"S":   "Sleeping in an interruptible wait",
		"D":   "Waiting in uninterruptible disk sleep",
		"Z":   "Zombie",
		"T":   "Stopped (on a signal)",
		"t":   "Tracing stop",
		"W":   "Paging",
		"X":   "Dead",
		"x":   "Dead",
		"K":   "Wakekill",
		"P":   "Parked",
		"N/A": "N/A",
	}
	for key, value := range status {
		newStatus := dal.getStatus(key)
		if newStatus != value {
			t.Errorf("Expected %s but Got %s", value, newStatus)
		}
	}
}
