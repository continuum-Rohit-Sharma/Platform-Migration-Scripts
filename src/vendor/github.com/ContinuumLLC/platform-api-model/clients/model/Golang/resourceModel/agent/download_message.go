package agent

import "time"

// DownloadMessageType message type
const DownloadMessageType = "DOWNLOAD"

// DownloadMessage download message
type DownloadMessage struct {
	ExecutionID        string                `json:"executionID"`
	DownloadFolderPath string                `json:"downloadFolderPath"`
	Items              []DownloadMessageItem `json:"items"`
}

// DownloadMessageItem download message item
type DownloadMessageItem struct {
	URL            string            `json:"url"`
	Checksum       string            `json:"checksum"`
	KeepOriginName bool              `json:"keepOriginName"`
	ExpiredAt      time.Time         `json:"expiredAt"`
	Metadata       string            `json:"metadata"`
	Header         map[string]string `json:"header"`
}
