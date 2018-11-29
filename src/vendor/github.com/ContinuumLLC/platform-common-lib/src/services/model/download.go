package model

//DownloadFileInfo is a common struct for Downloading a file using common API struct instance for any Service
type DownloadFileInfo struct {
	Name        string
	Path        string
	ContentType string
}
