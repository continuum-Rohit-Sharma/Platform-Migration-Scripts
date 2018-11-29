package dg

import (
	"time"
)

//InstalledSoftwareMessage is the definition of /resources/dynamicGroups/installedSoftwareMessage.json
type InstalledSoftwareMessage struct {
	Name             string    `json:"name"`
	InstallationDate time.Time `json:"installation_date"`
}
