package asset

//AssetDrive is the struct definition of /resources/asset/assetDrive
type AssetDrive struct {
	Product            string           `json:"product,omitempty" cql:"product"`
	Manufacturer       string           `json:"manufacturer,omitempty" cql:"manufacturer"`
	MediaType          string           `json:"mediaType,omitempty" cql:"media_type"`
	InterfaceType      string           `json:"interfaceType,omitempty" cql:"interface_type"`
	LogicalName        string           `json:"logicalName,omitempty" cql:"logical_name"`
	SerialNumber       string           `json:"serialNumber,omitempty" cql:"serial_number"`
	Partitions         []string         `json:"partitions,omitempty" cql:"partitions"`
	SizeBytes          int64            `json:"sizeBytes" cql:"size_bytes"`
	NumberOfPartitions int              `json:"numberOfPartitions" cql:"number_of_partitions"`
	PartitionData      []DrivePartition `json:"partitionData,omitempty" cql:"partition_data"`
}
