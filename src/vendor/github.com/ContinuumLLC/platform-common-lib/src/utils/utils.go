//Package utils is to convert interface type to specific type
package utils

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"crypto/md5"
	"encoding/hex"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

//ToString converts an interface type to string
func ToString(v interface{}) string {
	t, ok := v.(string)
	if ok {
		return t
	}
	return ""
}

//ToTime converts an interface type to time
func ToTime(v interface{}) time.Time {
	t, ok := v.(time.Time)
	if ok {
		return t
	}
	return time.Time{}
}

//ToInt64 converts an interface type to int64
//interface{} holding an int will not be type casted to int64 and will return 0 as the result
func ToInt64(v interface{}) int64 {
	t, ok := v.(int64)
	if ok {
		return t
	}
	return 0
}

//ToInt converts an interface type to int
func ToInt(v interface{}) int {
	t, ok := v.(int)
	if ok {
		return t
	}
	return 0
}

//ToFloat64 converts an interface type to float64
func ToFloat64(v interface{}) float64 {
	t, ok := v.(float64)
	if ok {
		return t
	}
	return 0
}

//ToBool converts an interface type to bool
func ToBool(v interface{}) bool {
	t, ok := v.(bool)
	if ok {
		return t
	}
	return false
}

//ToStringArray converts an interface type to string array
func ToStringArray(v interface{}) []string {
	t, ok := v.([]string)
	if ok {
		return t
	}
	return []string{}
}

//ToStringMap converts an interface type to map[string]string
func ToStringMap(v interface{}) map[string]string {
	t, ok := v.(map[string]string)
	if ok {
		return t
	}
	return nil
}

//GetTransactionIDFromResponse retrieves transactionid from the Response header
func GetTransactionIDFromResponse(res *http.Response) string {
	if res == nil {
		return ""
	}
	return res.Header.Get(string(protocol.HdrTransactionID))
}

//GetTransactionIDFromRequest retrieves transactionID from the Request header
func GetTransactionIDFromRequest(req *http.Request) string {
	if req == nil {
		return ""
	}
	return req.Header.Get(string(protocol.HdrTransactionID))
}

//GetChecksumFromRequest retrives MD5 from the request header
func GetChecksumFromRequest(req *http.Request) string {
	if req == nil {
		return ""
	}
	return req.Header.Get(string(protocol.HdrContentMD5))
}

//TimeoutError This error will be send while function call timeout
const TimeoutError = "FunctionCallTimedOut"

//Call is a function to be executed
type Call func() (interface{}, error)

//TimeoutCall is a function to be called on timeout
type TimeoutCall func()

//ExecuteWithTimeout is a function to execute a function call with timeout
func ExecuteWithTimeout(call Call, timeoutCall TimeoutCall, duration time.Duration) (interface{}, error) {
	var err error
	var i interface{}

	ch := make(chan bool, 1)
	defer close(ch)
	go func() {
		i, err = call()
		ch <- true
	}()

	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		if timeoutCall != nil {
			timeoutCall()
		}
		return nil, errors.New(TimeoutError)
	case <-ch:
		return i, err
	}
}

//GetChecksum is a function to calculate MD5 hash value
func GetChecksum(message []byte) string {
	hasher := md5.New()
	hasher.Write(message)
	return hex.EncodeToString(hasher.Sum(nil))
}

//ValidateMessage checks if message is corrupted or not
func ValidateMessage(message []byte, receievedChecksum string) (bool, string) {
	hashValue := GetChecksum(message)
	if receievedChecksum != "" && hashValue != receievedChecksum {
		return false, hashValue
	}
	return true, hashValue
}

//PanicHandler is a function to handle in the panic
type PanicHandler func(err error)

//WithRecover is a function to call any function with recovery and gives callback to @PanicHandler by passing Error Message and Stack Trace
func WithRecover(fn func(), handler PanicHandler) {
	defer func() {
		if err := recover(); err != nil {
			if handler != nil {
				handler(fmt.Errorf("Message :: %v \nStack Trace :: \n%s", err, debug.Stack()))
			}
		}
	}()
	fn()
}
