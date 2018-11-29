package rest

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const versionType = "version"

var (
	versionMutex           = &sync.Mutex{}
	versionLastTimeSuccess int64
)

//GeneralInfo contains general info about service
type GeneralInfo struct {
	TimeStampUTC    time.Time `json:"timeStampUTC"`
	ServiceName     string    `json:"serviceName"`
	ServiceProvider string    `json:"serviceProvider"`
	ServiceVersion  string    `json:"serviceVersion"`
	Name            string    `json:"name"`
}

// Version represents message format for /version endpoint
type Version struct {
	GeneralInfo
	Type                 string   `json:"type"`
	BuildCommitSHA       string   `json:"buildCommitSHA"`
	Repository           string   `json:"repository"`
	SupportedAPIVersions []string `json:"supportedAPIVersions"`
	BuildNumber          string   `json:"buildNumber"`
}

// Update for custom fields update
func (v *Version) Update() error {
	return nil
}

// init initialise of Version structure
func (v *Version) version() {
	v.TimeStampUTC = time.Now().UTC()
	v.Type = versionType
}

// Versioner to ensure the integrity of the data
type Versioner interface {
	version()
	Update() error
}

// versionData instance for version data
var versionData Versioner

// RegistryVersion set versionData obj
func RegistryVersion(version Versioner) {
	versionData = version
}

// HandlerVersion used for returning version of service
func HandlerVersion(w http.ResponseWriter, _ *http.Request) {
	startTime := time.Now().Unix()

	versionMutex.Lock()
	defer versionMutex.Unlock()

	if startTime <= versionLastTimeSuccess {
		RenderJSON(w, versionData)
		return
	}

	versionData.version()
	if err := versionData.Update(); err != nil {
		SendInternalServerError(w, "Update error", fmt.Errorf("rest.HandlerVersion err: %s", err))
		return
	}

	versionLastTimeSuccess = time.Now().Unix()
	RenderJSON(w, versionData)
}
