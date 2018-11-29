package autoupdate

//PackageManifest is a struct defining the Package level manifest having operations list
type PackageManifest struct {
	Name                string      `json:"name,omitempty"`
	Type                string      `json:"type,omitempty"`
	Version             string      `json:"version,omitempty"`
	MinimumAgentVersion string      `json:"minimumAgentVersion,omitempty"`
	Operations          []Operation `json:"operations,omitempty"`
	UninstallOperations []Operation `json:"uninstallOperations,omitempty"`
	UnsupportedOS       []OS        `json:"unsupportedOS,omitempty"`
	Arch                []string    `json:"supportedArch,omitempty"`
	Backup              []string    `json:"backup,omitempty"`
}

//Operation is a struct defining the operation structure
type Operation struct {
	Type                          string   `json:"type,omitempty"`
	Action                        string   `json:"action,omitempty"`
	Name                          string   `json:"name,omitempty"`
	NewFileName                   string   `json:"newFileName,omitempty"`
	RestoreOnFailure              bool     `json:"restoreOnFailure,omitempty"`
	InstallationPath              string   `json:"installationPath,omitempty"`
	FileHash                      string   `json:"fileHash,omitempty"`
	URL                           string   `json:"url,omitempty"`
	Arguments                     []string `json:"arguments,omitempty"`
	ReadResult                    bool     `json:"readResult,omitempty"`
	Input                         string   `json:"input,omitempty"`
	FileHashType                  string   `json:"fileHashType,omitempty"`
	FileExecutionTimeoutinSeconds int      `json:"fileExecutionTimeoutinSeconds,omitempty"`
}
