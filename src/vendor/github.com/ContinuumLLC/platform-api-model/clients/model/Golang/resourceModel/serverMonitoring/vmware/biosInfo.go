package vmware

import (
	"time"
)

//BiosInfo describes Host BIOS info
type BiosInfo struct {
	BiosVersion          string     `json:"biosVersion,omitempty"`
	ReleaseDate          *time.Time `json:"releaseDate,omitempty"`
	Vendor               string     `json:"vendor,omitempty"`
	MajorRelease         int        `json:"majorRelease,omitempty"`
	MinorRelease         int        `json:"minorRelease,omitempty"`
	FirmwareMajorRelease int        `json:"firmwareMajorRelease,omitempty"`
	FirmwareMinorRelease int        `json:"firmwareMinorRelease,omitempty"`
}
