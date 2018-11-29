package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

func TestToString(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected string
	}{
		{nil, ""},
		{"a", "a"},
		{1, ""},
	}
	for _, d := range testCases {
		returned := ToString(d.v)
		if returned != d.expected {
			t.Errorf("Unexpected object/value in converting interface to string")
		}
	}
}

func TestToInt64(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected int64
	}{
		{nil, 0},
		{"a", 0},
		{int64(1), 1}, //int can not be type casted to int64 that is why it needs to be converted to int64
	}
	for _, d := range testCases {
		returned := ToInt64(d.v)
		if returned != d.expected {
			t.Errorf("Unexpected object/value in converting interface to int64")
		}
	}
}

func TestToInt(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected int
	}{
		{nil, 0},
		{"a", 0},
		{1, 1},
	}
	for _, d := range testCases {
		returned := ToInt(d.v)
		if returned != d.expected {
			t.Errorf("Unexpected object/value in converting interface to int")
		}
	}
}

func TestToStringArray(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expected []string
	}{
		{nil, []string{}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{0, []string{}},
	}
	for _, d := range testCases {
		returned := ToStringArray(d.v)
		if !reflect.DeepEqual(returned, d.expected) {
			t.Errorf("Unexpected object/value in converting interface to string array")
			break
		}
	}
}

func TestGetTransactionIDFromResponse(t *testing.T) {
	id := "1"
	r := &http.Response{}
	r.Header = make(http.Header, 1)
	r.Header.Set(string(protocol.HdrTransactionID), id)

	tid := GetTransactionIDFromResponse(r)
	if tid != id {
		t.Errorf("Unexpected transactionId returned , Expected:%s, Returned:%s", id, tid)
	}

}

func TestGetTransactionIDFromResponseNilCheck(t *testing.T) {

	tid := GetTransactionIDFromResponse(nil)
	if tid != "" {
		t.Errorf("Unexpected transactionId returned , Returned:%s", tid)
	}

}

func TestGetTransactionIDFromRequest(t *testing.T) {
	id := "1"
	r := &http.Request{}
	r.Header = make(http.Header, 1)
	r.Header.Set(string(protocol.HdrTransactionID), id)

	tid := GetTransactionIDFromRequest(r)
	if tid != id {
		t.Errorf("Unexpected transactionId returned , Expected:%s, Returned:%s", id, tid)
	}

}

func TestGetTransactionIDFromRequestNilCheck(t *testing.T) {

	tid := GetTransactionIDFromRequest(nil)
	if tid != "" {
		t.Errorf("Unexpected transactionId returned , Returned:%s", tid)
	}

}

func TestGetChecksumFromRequest(t *testing.T) {

	tid := GetChecksumFromRequest(nil)
	if tid != "" {
		t.Errorf("Unexpected transactionId returned , Returned:%s", tid)
	}

}

func TestExecuteWithTimeoutSucess(t *testing.T) {
	i, err := ExecuteWithTimeout(systemCall, systemCallTimeoutHandler, 5*time.Second)
	if err != nil && i != nil {
		t.Errorf("Expected NIL but GOT Error : %v", err)
	}
}

func TestExecuteWithTimeoutFail(t *testing.T) {
	_, err := ExecuteWithTimeout(systemCall, systemCallTimeoutHandler, 2*time.Second)
	if err == nil {
		t.Errorf("Expected Error but GOT nil")
	}
}
func TestGetChecksum(t *testing.T) {
	validChecksum := "0cbc6611f5540bd0809a388dc95a615b"
	message := []byte("Test")
	checksum := GetChecksum(message)
	if checksum != validChecksum {
		t.Errorf("Failed!:%v", checksum)
	}
}

func TestGetChecksumBlank(t *testing.T) {
	validChecksum := "d41d8cd98f00b204e9800998ecf8427e"
	message := []byte{}
	checksum := GetChecksum(message)
	if checksum != validChecksum {
		t.Errorf("Failed!: Unexpected checksum %s", checksum)
	}
}

func TestValidateMessageFail(t *testing.T) {
	checksum := "TestString"
	flag, _ := ValidateMessage([]byte("TestString"), checksum)
	if flag {
		t.Errorf("Failed!: Function should return false when checksum is invalid")
	}
}

func TestValidateMessageSuccess(t *testing.T) {
	message := []byte("TestString")
	checksum := GetChecksum(message)
	flag, _ := ValidateMessage(message, checksum)
	if !flag {
		t.Errorf("Failed!: Function should return true when checksum %s", checksum)
	}
}

func TestValidateMessageBlankCheckSum(t *testing.T) {
	message := []byte("TestString")
	flag, _ := ValidateMessage(message, "")
	if !flag {
		t.Errorf("Failed!: Function should return true when checksum blank")
	}
}

func systemCall() (interface{}, error) {
	fmt.Println("In systemCall method")
	time.Sleep(3 * time.Second)
	return "Success", nil
}

func systemCallTimeoutHandler() {
	fmt.Println("In systemCallTimeoutHandler method")
}

type RecoverTest struct {
	Panic bool
	T     *testing.T
}

func (r RecoverTest) PanicHandler(err error) {
	if r.Panic && err == nil {
		r.T.Errorf("Expected Error but Got Nil")
	}

	if !r.Panic && err != nil {
		r.T.Errorf("Expected nil but Got Error %v", err)
	}
}

func (r RecoverTest) PanicGenerator() {
	if r.Panic {
		panic(fmt.Errorf("Generating Panic"))
	}
}

func TestWithRecover(t *testing.T) {

	type args struct {
		fn      func()
		handler PanicHandler
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{fn: RecoverTest{Panic: true, T: t}.PanicGenerator, handler: RecoverTest{Panic: true, T: t}.PanicHandler},
		},
		{
			name: "1",
			args: args{fn: RecoverTest{Panic: false, T: t}.PanicGenerator, handler: RecoverTest{Panic: false, T: t}.PanicHandler},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WithRecover(tt.args.fn, tt.args.handler)
		})
	}
}

func TestToBool(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TC1_True",
			args: args{
				v: true,
			},
			want: true,
		},
		{
			name: "TC2_NIL",
			args: args{
				v: nil,
			},
			want: false,
		},
		{
			name: "TC3_Int",
			args: args{
				v: 1,
			},
			want: false,
		},
		{
			name: "TC4_False",
			args: args{
				v: false,
			},
			want: false,
		},
		{
			name: "TC5_Float64",
			args: args{
				v: float64(1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBool(tt.args.v); got != tt.want {
				t.Errorf("ToBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "TC1_String",
			args: args{
				v: "abc",
			},
			want: 0,
		},
		{
			name: "TC2_Nil",
			args: args{
				v: nil,
			},
			want: 0,
		},
		{
			name: "TC3_float64",
			args: args{
				v: float64(1),
			},
			want: 1,
		},
		{
			name: "TC4_false",
			args: args{
				v: false,
			},
			want: 0,
		},
		{
			name: "TC5_true",
			args: args{
				v: true,
			},
			want: 0,
		},
		{
			name: "TC6_true",
			args: args{
				v: int64(23),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFloat64(tt.args.v); got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStringMap(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "TC1",
			args: args{v: nil},
			want: nil,
		},
		{
			name: "TC2",
			args: args{v: map[string]string{
				"name":  "a",
				"value": "a",
			},
			},
			want: map[string]string{
				"name":  "a",
				"value": "a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStringMap(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
