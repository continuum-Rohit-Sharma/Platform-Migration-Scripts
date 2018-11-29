package services

import (
	"os"
	"strconv"

	aModel "github.com/ContinuumLLC/platform-api-model/clients/model/Golang/resourceModel/healthCheck"
	"github.com/ContinuumLLC/platform-common-lib/src/logging"
	"github.com/ContinuumLLC/platform-common-lib/src/procParser"
	"github.com/ContinuumLLC/platform-common-lib/src/services/model"
)

//IfConfig constants
const (
	cIPAddressRowIndex    int = 1
	cIPAddressColumnIndex int = 2
)

const (
	cProcPath = "/proc/%d/status"
)

//HealthCheckDalFactoryImpl returns the concrete implementation of Factory
type HealthCheckDalFactoryImpl struct {
}

//GetHealthCheckDal : A factory function to create an instance of HealthCheck Dal
func (HealthCheckDalFactoryImpl) GetHealthCheckDal(f model.HealthCheckDependencies) model.HealthCheckDal {
	return healthCheckDalImpl{
		f:   f,
		log: logging.GetLoggerFactory().Get(),
	}
}

//healthCheckDalImpl returns the concrete implementation of HealthCheckDal
type healthCheckDalImpl struct {
	f   model.HealthCheckDependencies
	log logging.Logger
}

func (h healthCheckDalImpl) GetHealthCheck(healthCheck model.HealthCheck) (aModel.HealthCheck, error) {
	h.log.LogWithCorrelationf(logging.DEBUG, healthCheck.Version.ServiceName, "Retriving Health Information for %v", healthCheck)
	health := &aModel.HealthCheck{}
	var err error
	versionService := h.f.GetVersionService()
	health.Version = versionService.GetVersion(healthCheck.Version)
	health.Type = "HealthCheck"
	health.NetworkInterfaces, err = h.getNetworkInterface(healthCheck)
	if err != nil {
		h.log.LogWithCorrelationf(logging.ERROR, healthCheck.Version.ServiceName, "Error while finding network interface data %v", err)
		return *health, err
	}

	procID := os.Getpid()
	err = h.findProcessInformation(procID, health)
	if err != nil {
		h.log.LogWithCorrelationf(logging.ERROR, healthCheck.Version.ServiceName, "Error while finding Process Information %v", err)
		return *health, err
	}
	return *health, nil
}

func (h healthCheckDalImpl) findProcessInformation(procID int, health *aModel.HealthCheck) error {
	data, err := h.executeCommand(procParser.ModeTabular, "ps", "-f", "-p", strconv.Itoa(procID), "-eo", "pid,pcpu,nlwp,pmem,state")
	if err != nil {
		h.log.LogWithCorrelationf(logging.ERROR, health.Version.ServiceName, "Error while \"%s\" command for %d id : %v", "ps -f -p <ProcId> -eo pid,pcpu,nlwp,pmem,state", procID, err)
		return err
	}
	psResult := data.Lines[1]
	health.CPUPercentage, err = strconv.ParseFloat(psResult.Values[1], 64)
	if err != nil {
		health.CPUPercentage = 0.0
		h.log.LogWithCorrelationf(logging.ERROR, health.Version.ServiceName, "Error while converting CPU persent for %d : %v", procID, err)
	}
	health.NumOfOSThreads, err = strconv.Atoi(psResult.Values[2])
	if err != nil {
		h.log.LogWithCorrelationf(logging.ERROR, health.Version.ServiceName, "Error while converting Thread Count for %d : %v", procID, err)
		health.NumOfOSThreads = 0
	}
	health.MemoryPercentage, err = strconv.ParseFloat(psResult.Values[3], 64)
	if err != nil {
		h.log.LogWithCorrelationf(logging.ERROR, health.Version.ServiceName, "Error while converting Memory persent for %d : %v", procID, err)
		health.MemoryPercentage = 0.0
	}
	health.Status = h.getStatus(psResult.Values[4])
	health.StartTime = model.StrartTime
	return nil
}

func (h healthCheckDalImpl) getNetworkInterface(healthCheck model.HealthCheck) ([]string, error) {
	data, err := h.executeCommand(procParser.ModeKeyValue, "ifconfig", "-a")
	if err != nil {
		h.log.LogWithCorrelationf(logging.ERROR, healthCheck.Version.ServiceName, "Error while executing \"ifconfig -a\" command %v", err)
		return nil, err
	}
	ifConfigData := h.parseIfConfigData(data)

	networkInterface := make([]string, 0)
	for _, v := range ifConfigData {
		ipAddress := v[cIPAddressRowIndex].Values[cIPAddressColumnIndex]
		networkInterface = append(networkInterface, ipAddress+healthCheck.ListenPort)
	}
	return networkInterface, nil
}

func (h healthCheckDalImpl) parseIfConfigData(data *procParser.Data) map[string][]procParser.Line {
	result := make(map[string][]procParser.Line)
	key := ""
	lenLine := len(data.Lines)
	for i := 0; i < lenLine; i++ {
		line := data.Lines[i]
		if key == "" {
			key = line.Values[0]
		}
		if len(line.Values) == 0 {
			key = ""
		} else {
			lines := result[key]
			if lines == nil {
				lines = []procParser.Line{line}
				result[key] = lines
			} else {
				lines = append(lines, line)
				result[key] = lines
			}
		}
	}
	return result
}

func (h healthCheckDalImpl) executeCommand(mode procParser.Mode, command string, arg ...string) (*procParser.Data, error) {
	reader, err := h.f.GetEnv().GetCommandReader(command, arg...)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	cfg := procParser.Config{
		ParserMode: mode,
	}
	return h.f.GetParser().Parse(cfg, reader)
}

func (h healthCheckDalImpl) getStatus(status string) string {
	switch status {
	case "R":
		return "Running"
	case "S":
		return "Sleeping in an interruptible wait"
	case "D":
		return "Waiting in uninterruptible disk sleep"
	case "Z":
		return "Zombie"
	case "T":
		return "Stopped (on a signal)"
	case "t":
		return "Tracing stop"
	case "W":
		return "Paging"
	case "X":
		return "Dead"
	case "x":
		return "Dead"
	case "K":
		return "Wakekill"
	case "P":
		return "Parked"
	}
	return "N/A"
}
