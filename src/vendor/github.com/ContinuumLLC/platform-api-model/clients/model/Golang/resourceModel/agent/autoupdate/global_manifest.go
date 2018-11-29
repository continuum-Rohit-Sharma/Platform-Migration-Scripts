package autoupdate

const (
	//InstallOperation ...
	InstallOperation = "install"
	//UninstallOperation ...
	UninstallOperation = "uninstall"
	//UpdateOperation ...
	UpdateOperation = "update"
)

//GlobalManifest is a struct defining the Global manifest having package list
type GlobalManifest struct {
	ProductVersion string    `json:"productVersion,omitempty"`
	Version        string    `json:"version,omitempty"`
	SupportedOS    []OS      `json:"supportedOS,omitempty"`
	SupportedArch  []string  `json:"supportedArch,omitempty"`
	Packages       []Package `json:"packages,omitempty"`
}

//OS is a struct defining the Operating System Details
type OS struct {
	Name    string `json:"name,omitempty" cql:"name"`
	Type    string `json:"type,omitempty" cql:"type"`
	Version string `json:"version,omitempty" cql:"version"`
}

//Package is a struct defining the Package Details
type Package struct {
	Name      string `json:"name,omitempty" cql:"name"`
	Type      string `json:"type,omitempty" cql:"type"`
	Version   string `json:"version,omitempty" cql:"version"`
	SourceURL string `json:"sourceUrl,omitempty" cql:"source_url"`
	Operation string `json:"operation,omitempty"`
}
