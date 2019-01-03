package main

import (
	"fmt"
	"reflect"
	"time"
)

// Name of the struct tag used in examples
const tagName = "column_key"
const map_column_tag_name = "map_column"

type Asset struct {
	ID     int       `column_key:"id"`
	SId    string    `column_key:"sid"`
	RId    string    `column_key:"rid"`
	DCTime time.Time `column_key:"dc_time"`
}

//NewAsset is the New Table where it's name should corresponds to the new Table name
type NewAsset struct {
	EndpointID int       `column_key:"id",map_column:"id"`
	SiteID     string    `column_key:"site_id",map_column:"sid`
	RegID      string    `column_key:"reg_id",map_column:"rid`
	DCTimeUTC  time.Time `column_key:"dc_timestamp_utc",map_column:"dc_time`
}

func main() {
	asset := Asset{
		ID:     1,
		SId:    "s1",
		RId:    "r1",
		DCTime: time.Now(),
	}

	newAsset := NewAsset{}

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := reflect.TypeOf(asset)
	t2 := reflect.TypeOf(newAsset)
	// Get the type and kind of our user variable
	fmt.Println("Type:", t.Name())
	fmt.Println("Kind:", t.Kind())

	fmt.Println("new_Type:", t.Name())
	fmt.Println("new_Kind:", t.Kind())

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)
		field2 := t2.Field(i)
		// Get the field tag value
		tag := field.Tag.Get(tagName)
		tag_mapColumn := field2.Tag.Get(map_column_tag_name)
		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field2.Name, field2.Type.Name(), tag_mapColumn)
	}

}
