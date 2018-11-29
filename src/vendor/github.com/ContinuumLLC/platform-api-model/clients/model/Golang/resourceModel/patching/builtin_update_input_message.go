package patching

import "os"

// Built in update type
const (
	BuiltInUpdateRegistry         = "registry"
	BuiltInUpdateConfigFile       = "config_file"
	BuiltInUpdateFolderPermission = "folder_permission"
)

// BuiltInUpdateInputMessage input message
type BuiltInUpdateInputMessage struct {
	ExecutionID string `json:"executionID"`
	// Type options:
	// - registry
	// - config_file
	// - folder_permission
	Type string `json:"type"`
	// It's a JSON field, witch can be:
	// - RegistryParams
	// - ConfigFileParams
	// - FolderPermissionParams
	Params []byte `json:"params"`
}

// BuiltInUpdateInputMessageList is a list of BuiltInUpdateInputMessage
type BuiltInUpdateInputMessageList []BuiltInUpdateInputMessage

// RegistryParams registry params
type RegistryParams struct {
	File     string            `json:"file"`
	KeyValue map[string]string `json:"keyValue"`
}

// ConfigFileParams config file params
type ConfigFileParams struct {
	Files   []string `json:"files"`
	Find    string   `json:"find"`
	Replace string   `json:"replace"`
}

// FolderPermissionParams folder permission params
type FolderPermissionParams struct {
	Dir string `json:"dir"`
	// Mode options:
	// Unix chmod numeric notation (0755, 0400)
	Mode os.FileMode `json:"mode"`
}
