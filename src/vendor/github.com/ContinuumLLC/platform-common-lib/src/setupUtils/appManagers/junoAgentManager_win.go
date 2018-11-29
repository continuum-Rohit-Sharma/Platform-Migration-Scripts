//go:generate goversioninfo
// +build windows

package appManagers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type junoAgentManager struct {
	AppPath string
	Params  []string
}

func (r junoAgentManager) Install(config *Config) error {
	err := r.downloadExe(config.JunoAgentURL, config.JunoAgentExeLocation)
	if err != nil {
		log.Printf("Error downloading Juno agent, error: %v \n", err)
		return err
	}

	err = runExe(r.AppPath, r.Params[0], r.Params[1], r.Params[2])
	if err != nil {
		log.Printf("Error installing Juno agent, error: %v \n", err)
		return err
	}
	return err
}

func (r junoAgentManager) Uninstall(config *Config) error {
	err := runExe(r.AppPath, r.Params[0], r.Params[1], r.Params[2])
	if err != nil {
		log.Printf("Error uninstalling Juno agent, error: %v \n", err)
		return err
	}
	r.cleanupJunoFiles(r.AppPath)
	return err
}
func (r junoAgentManager) downloadExe(fileurl string, outputFilename string) error {
	file, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(fileurl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func (r junoAgentManager) cleanupJunoFiles(path string) error {
	os.Remove(path)
	return nil
}

func (r junoAgentManager) checkOSVerion(unSupportedOS string) error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	osName, _, err := k.GetStringValue("ProductName")
	if err != nil {
		return err
	}
	osName = strings.ToLower(osName)
	uos := strings.Split(unSupportedOS, ",")
	for _, v := range uos {
		if strings.Contains(osName, strings.ToLower(v)) {
			log.Printf("Error installing Juno agent, OS Not Supported")
			return errors.New("OS Not Supported")
		}
	}
	return nil
}

func (r junoAgentManager) checkJunoVersion() {

}
