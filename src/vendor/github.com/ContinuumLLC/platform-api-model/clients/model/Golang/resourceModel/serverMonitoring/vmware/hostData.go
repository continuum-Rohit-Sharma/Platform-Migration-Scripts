package vmware

import (
	"time"
)

// HostData is the struct definition of VMWare Host instance info
type HostData struct {
	MonitoringGroupType string              `json:"monitoringGroupType,omitempty"`
	CreateTimeUTC       *time.Time          `json:"createTimeUTC,omitempty"`
	HostName            string              `json:"hostName,omitempty"`
	IPAddress           string              `json:"ipAddress,omitempty"`
	ProbeServerName     string              `json:"probeServerName,omitempty"`
	LastBootTimeUTC     *time.Time          `json:"lastBootTimeUTC,omitempty"`
	Os                  *Os                 `json:"os,omitempty"`
	Bios                *BiosInfo           `json:"bios,omitempty"`
	MaintenanceMode     bool                `json:"maintenanceMode,omitempty"`
	PerformanceMetrics  *PerformanceMetrics `json:"performanceMetrics,omitempty"`
	Processors          []Processor         `json:"processors,omitempty"`
	DataStores          []DataStore         `json:"dataStores,omitempty"`
	Services            []Service           `json:"services,omitempty"`
	Sensors             []Sensor            `json:"sensors,omitempty"`
	GuestInstances      []GuestInstance     `json:"guestInstances,omitempty"`
	PCIDevices          []PCIDevice         `json:"pciDevices,omitempty"`
	Events              []Event             `json:"events,omitempty"`
}
