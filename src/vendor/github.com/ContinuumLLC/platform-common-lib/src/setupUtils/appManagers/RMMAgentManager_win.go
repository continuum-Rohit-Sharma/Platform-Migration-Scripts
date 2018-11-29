//go:generate goversioninfo
// +build windows

package appManagers

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"

	ini "gopkg.in/ini.v1"
)

type rmmAgentManager struct {
	AppPath string
	Params  []string
}

func (r rmmAgentManager) Install(config *Config) error {
	log.Printf("Cleaning up old INI File")
	err := cleanupOldINI(config.RMMAgentINIFilePath)
	if err != nil {
		log.Printf("Error cleaning up old INI file, error: %v \n", err)
		return err
	}

	err = makeIniFile(config)
	if err != nil {
		log.Printf("Error found while creating ini, error: %v \n", err)
		return err
	}
	log.Printf("Updating INI File")

	err = updateINI(config.RMMAgentINIFilePath, r.Params[0], r.Params[1], r.Params[2])
	if err != nil {
		log.Printf("Error found while updating file %s , with error %v \n", os.Args[2], err)
		return err
	}

	log.Println("INI file updated successfully")
	log.Println("Installing RMM Agent")
	err = installRMMAgent(r.AppPath)

	if err != nil {
		log.Printf("Error found while executing file, error: %v \n", err)
		log.Println("RMM Agent Installation Failed")
		return err
	}
	err = runAgentCore(config)
	if err != nil {
		log.Printf("Error found while executing AgentCore, error: %v \n", err)
		log.Println("Agent Registration Failed")
		return err
	}

	log.Println("RMM Agent Installed successfully")
	return nil
}

func (rmmAgentManager) Uninstall() error {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\SAAZOD", registry.QUERY_VALUE|registry.WOW64_32KEY)
	if err != nil {
		log.Printf("Error in getting registry key for unistall, Error: %v \n", err)
		return nil
	}
	defer k.Close() //nolint

	val, _, err := k.GetStringValue("UninstallString")
	if err != nil {
		return err
	}
	uninstallCommand := strings.Split(val, " ")
	err = runExe(uninstallCommand[0], uninstallCommand[1], "", "")
	if err != nil {
		return err
	}
	return nil
}

func installRMMAgent(exePath string) error {
	err := runExe(exePath, "/s", "", "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func makeIniFile(config *Config) error {
	from, err := os.OpenFile(config.RMMAgentBackupINIFilePath, os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}

	to, err := os.OpenFile(config.RMMAgentINIFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func updateINI(iniPath string, partnerID string, siteID string, agentType string) error {
	cfg, err := ini.InsensitiveLoad(iniPath)
	if err != nil {
		fmt.Println(err)
	}

	secCtrlInfo, err := cfg.GetSection("controller_info")
	if err != nil {
		fmt.Println(err)
	}
	memKey, err := secCtrlInfo.GetKey("memberid")
	if err != nil {
		fmt.Println(err)
	}
	memKey.SetValue(partnerID)

	siteKey, err := secCtrlInfo.GetKey("siteid")
	if err != nil {
		fmt.Println(err)
	}
	siteKey.SetValue(siteID)

	secmachInfo, err := cfg.GetSection("machine_info")
	if err != nil {
		fmt.Println(err)
	}

	agentKey, err := secmachInfo.GetKey("agentType")
	if err != nil {
		fmt.Println(err)
	}
	agentKey.SetValue(agentType)

	err = cfg.SaveTo(iniPath)
	return err
}

func runAgentCore(config *Config) error {
	err := runExe(config.AgentCorePath, config.AgentConfigFilePath, config.AppManagerLogFilePath, "registration")
	return err
}

func cleanupOldINI(path string) error {
	os.Remove(path)
	return nil
}
