package main

import (
	"time"
)

//Config contains
type Config struct {
	Source      string
	Delta       string
	Destination string
	Action      string
	mode        IMergeMode
}

//Schedules contains
type Schedules struct {
	Schedule []struct {
		Name       string    `json:"name"`
		Type       string    `json:"type"`
		Time       time.Time `json:"timestampUTC"`
		Version    string    `json:"version"`
		Task       string    `json:"task"`
		Path       string    `json:"path"`
		ExecuteNow string    `json:"executeNow"`
		Schedue    string    `json:"schedule"`
	} `json:"schedules"`
}

//MergeService provides
type MergeService interface {
	Merge(cfg *Config) error
}

// //Handler provides
// type Handler interface {
// 	AddingNew(map[string]interface{}, string, interface{})
// 	Removing(map[string]interface{}, string)
// 	ActualMerge(map[string]interface{}, string, interface{})
// }

//IMergeMode provides
type IMergeMode interface {
	SourceHasKey(src, dest *map[string]interface{}, key string, value interface{})
	SourceNotHasKey(src, dest *map[string]interface{}, key string, value interface{})
	IfArray(src, dest *map[string]interface{}, key string)
}
