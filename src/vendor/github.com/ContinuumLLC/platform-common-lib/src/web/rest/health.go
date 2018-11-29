package rest

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	healthType           = "health"
	healthRunningStatus  = "Running"
	healthDegradedStatus = "Degraded"
	dbMessageType        = "OutboundConnectionStatus"
	dbNameSuffix         = "OutboundConnectionStatus"
	//ConnectionStatusActive correct connection
	ConnectionStatusActive = "Active"
	//ConnectionStatusUnavailable fail connection
	ConnectionStatusUnavailable = "Unavailable"
)

var (
	//TimeStampUTC is a variable holds service start time and gets updated at bootstrap
	TimeStampUTC          time.Time
	healthMutex           = &sync.Mutex{}
	healthLastTimeSuccess int64
)

func init() {
	TimeStampUTC = time.Now().UTC()
}

// Statuser interface for outbound connection status
type Statuser interface {
	Status(status OutboundConnectionStatus) *OutboundConnectionStatus
}

// Health represents message format for /health endpoint
type Health struct {
	GeneralInfo
	Type                     string                     `json:"type"`
	Status                   string                     `json:"status"`
	LastStartTimeUTC         time.Time                  `json:"lastStartTimeUTC"`
	NetworkInterfaces        []string                   `json:"networkInterfaces"`
	OutboundConnectionStatus []OutboundConnectionStatus `json:"outboundConnectionStatus"`
	ConnMethods              []Statuser                 `json:"-"`
	ListenURL                string                     `json:"-"`
}

// Update for custom fields update
func (h *Health) Update() error {
	return nil
}

// Healther to ensure the integrity of the data
type Healther interface {
	health() error
	Update() error
}

// healthData instance for health data
var healthData Healther

// RegistryHealth set healthData obj
func RegistryHealth(health Healther) {
	healthData = health
}

// init initialise of Health structure
func (h *Health) health() (err error) {
	h.Type = healthType
	h.TimeStampUTC = time.Now().UTC()
	h.LastStartTimeUTC = TimeStampUTC

	h.OutboundConnectionStatus = GetOutboundConnectionStatus(h.ConnMethods, h.ServiceName)
	h.Status = GetHealthStatus(h.OutboundConnectionStatus)

	h.NetworkInterfaces, err = GetNetworkInterfaces(h.ListenURL)
	if err != nil {
		Log.Error("Health.init: ", err)
		return
	}
	return
}

// GetHealthStatus used for getting status of common connection state
func GetHealthStatus(conns []OutboundConnectionStatus) string {
	for _, conn := range conns {
		if conn.ConnectionStatus != ConnectionStatusActive {
			return healthDegradedStatus
		}
	}
	return healthRunningStatus
}

// HandlerHealth used for returning health status of service
func HandlerHealth(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Unix()

	healthMutex.Lock()
	defer healthMutex.Unlock()

	if startTime <= healthLastTimeSuccess {
		RenderJSON(w, healthData)
		return
	}

	if err := healthData.health(); err != nil {
		SendInternalServerError(w, "Network error", fmt.Errorf("rest.HandlerHealth err: %s", err))
		return
	}
	if err := healthData.Update(); err != nil {
		SendInternalServerError(w, "Update error", fmt.Errorf("rest.HandlerHealth err: %s", err))
		return
	}

	healthLastTimeSuccess = time.Now().Unix()
	RenderJSON(w, healthData)
}

// OutboundConnectionStatus represents status of all connections into service
type OutboundConnectionStatus struct {
	TimeStampUTC     time.Time `json:"timeStampUTC"`
	Type             string    `json:"type"`
	Name             string    `json:"name"`
	ConnectionType   string    `json:"connectionType"`
	ConnectionURLs   []string  `json:"connectionURLs"`
	ConnectionStatus string    `json:"connectionStatus"`
}

// GetOutboundConnectionStatus used for getting status of all connections in service
func GetOutboundConnectionStatus(methods []Statuser, serviceName string) []OutboundConnectionStatus {
	connections := make([]OutboundConnectionStatus, 0, len(methods))
	baseConn := OutboundConnectionStatus{
		TimeStampUTC: time.Now().UTC(),
		Type:         dbMessageType,
		Name:         fmt.Sprintf("%s-%s", serviceName, dbNameSuffix),
	}

	for _, conn := range methods {
		connections = append(connections, *conn.Status(baseConn))
	}
	return connections
}
