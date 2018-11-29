package asset

//PhysicalDrive is the struct definition of /resources/asset/assetPhysicalDrive
type PhysicalDrive struct {
	Type          string           `json:"type,omitempty" cql:"type"`
	PartitionData []DrivePartition `json:"partitionData,omitempty" cql:"partition_data"`
}
