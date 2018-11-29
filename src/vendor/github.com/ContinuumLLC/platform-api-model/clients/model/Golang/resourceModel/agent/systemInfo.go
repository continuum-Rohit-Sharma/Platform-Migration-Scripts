package agent

import "time"

const (
	//ResourceTypeDesktop represents desktop
	ResourceTypeDesktop string = "Desktop"
	//ResourceTypeServer represents server
	ResourceTypeServer string = "Server"
)

//SystemInfo is the struct definition of system info structure for agent provisioning data
type SystemInfo struct {
	TimestampUTC                  time.Time `json:"timestamp_utc" cql:"timestamp_utc"`
	OSType                        string    `json:"osType"  cql:"os_type"`
	OSName                        string    `json:"osName"  cql:"os_name"`
	OSVersion                     string    `json:"osVersion"  cql:"os_version"`
	OSSerialNumber                string    `json:"osSerialNumber"  cql:"os_serial_number"`
	HostName                      string    `json:"hostName"  cql:"host_name"`
	MACAddress                    string    `json:"macAddress"  cql:"mac_address"`
	ProcessorID                   string    `json:"processorid"  cql:"processor_id"`
	ProcessorType                 string    `json:"processorType"  cql:"processor_type"`
	HardDriveSerialNumber         string    `json:"hardDriveSerialNumber"  cql:"hard_drive_serial_number"`
	Memory                        string    `json:"memory"  cql:"memory"`
	MotherboardAdapter            string    `json:"motherboardAdapter"  cql:"motherboard_adapter"`
	CDROMSerial                   string    `json:"cdromSerial"  cql:"cdrom_serial"`
	LogicalDiskVolumeSerialNumber string    `json:"logicalDiskVolumeSerialNumber"  cql:"logicaldisk_volumeserialnumber"`
	BiosSerial                    string    `json:"biosSerial"  cql:"bios_serial"`
	VirtualMachineUUID            string    `json:"virtualMachineIdentity"  cql:"virtual_machine_identity"`
	SystemManufacturerRef         string    `json:"systemManufacturerReference"  cql:"system_manufacturer_reference"`
	Mode                          string    `json:"mode" cql:"mode"`
	EndpointType                  string    `json:"endpointType" cql:"endpoint_type"`
}
