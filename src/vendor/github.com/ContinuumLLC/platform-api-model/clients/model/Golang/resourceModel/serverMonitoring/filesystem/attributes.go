package filesystem

import "time"

// Attributes represent file / directory attributes
type Attributes struct {
	IsDirectory    bool       `json:"isDirectory,omitempty"`
	SizeInBytes    int64      `json:"sizeInBytes,omitempty"`
	Created        *time.Time `json:"created,omitempty"`
	LastModified   *time.Time `json:"lastModified,omitempty"`
	LastAccessed   *time.Time `json:"lastAccessed,omitempty"`
	FileVersion    string     `json:"fileVersion,omitempty"`
	ProductVersion string     `json:"productVersion,omitempty"`
}
