package main

import (
	"log"
	"os"
	"strings"

	"github.com/ContinuumLLC/platform-common-lib/src/json"
	"github.com/ContinuumLLC/platform-common-lib/src/setupUtils/appManagers"
)

func main() {
	config, err := GetConfig()
	f, err := os.OpenFile(config.AppManagerLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("Starting Utility")
	log.Printf("Flag Passed: %s \n", os.Args[1])

	flag := os.Args[1]
	appName := ""
	appPath := ""
	params := []string{}

	if err != nil {
		log.Printf("Error reading config file, error: %v \n", err)
		os.Exit(1)

	}
	switch strings.ToLower(flag) {
	case "install":
		appName = "JunoAgent"
		appPath = config.JunoAgentExeLocation
		params = []string{}
		manager := appManagers.GetManager(appName, appPath, params)
		err := manager.Install(config)
		if err != nil {
			log.Printf("Error found while creating ini, error: %v \n", err)
			os.Exit(1)
		}
	case "uninstall":
		manager := appManagers.GetManager(appName, appPath, params)
		err := manager.Uninstall(config)
		if err != nil {
			log.Printf("Error found while creating ini, error: %v \n", err)
			os.Exit(1)
		}
	}
}

//GetConfig gets the configuration for Appinstaller
func GetConfig() (*appManagers.Config, error) {
	config := &appManagers.Config{}
	err := json.FactoryJSONImpl{}.GetDeserializerJSON().ReadFile(config, "..\\appManager\\appManager_cfg.json")
	return config, err
}
