package patching

// CleanUpInputMessage is a definition of /resources/patching/patchingCleanUpMessage.json
type CleanUpInputMessage struct {
	DownloadFolderPath string `json:"downloadFolderPath" description:"Path to downloaded folder"`
}
