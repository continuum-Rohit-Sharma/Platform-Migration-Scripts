package asset

//DrivePartition is the struct definition of /resources/asset/assetDrivePartition
type DrivePartition struct {
	Name        string `json:"name,omitempty" cql:"name"`
	Label       string `json:"label,omitempty" cql:"label"`
	FileSystem  string `json:"fileSystem,omitempty" cql:"file_system"`
	Description string `json:"description,omitempty" cql:"description"`
	SizeBytes   int64  `json:"sizeBytes" cql:"size_bytes"`
}
