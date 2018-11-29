package wmi

import (
	"github.com/StackExchange/wmi"
)

//Wrapper for the wmi to execute/create query
type Wrapper interface {
	//Query to execute the WMI query and wrap the output to the dst type(struct) which must be passed as reference
	Query(query string, dst interface{}, connectServerArgs ...interface{}) error
	//CreateQuery is to create the WMI query based on the src type
	CreateQuery(src interface{}, where string) string
}

//GetWrapper returns the implementation of WMI Wrapper
func GetWrapper() StackExchangeWMI {
	return StackExchangeWMI{}
}

//StackExchangeWMI is the implementation for the WMI Wrapper
type StackExchangeWMI struct{}

//Query to execute the WMI query and wrap the output to the dst type(struct) which must be passed as reference
func (StackExchangeWMI) Query(query string, dst interface{}, connectServerArgs ...interface{}) error {
	return wmi.Query(query, dst, connectServerArgs...)
}

//CreateQuery is to create the WMI query based on the src type
func (StackExchangeWMI) CreateQuery(src interface{}, where string) string {
	return wmi.CreateQuery(src, where)
}
