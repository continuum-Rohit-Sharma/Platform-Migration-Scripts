package procParser

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestGetInt64(t *testing.T) {
	input := "1"
	val, err := GetInt64(input)
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}
	if val != 1 {
		t.Errorf("Unexpected value returned: %d", val)
	}
}

func TestGetInt64InvalidInput(t *testing.T) {
	input := "invalid"
	_, err := GetInt64(input)
	if err == nil || !strings.Contains(err.Error(), "parsing") {
		t.Errorf("Unexpected error returned: %v", err)
	}
}

func TestGetBytes(t *testing.T) {

	var err error
	err = testGetBytes(1, "kB", 1024)
	if err != nil {
		t.Error(err.Error())
	}
	err = testGetBytes(1, "mB", 1024*1024)
	if err != nil {
		t.Error(err.Error())
	}
	err = testGetBytes(1, "gB", 1024*1024*1024)
	if err != nil {
		t.Error(err.Error())
	}

	err = testGetBytes(1, "INVALID", 0)
	if err != nil {
		t.Error(err.Error())
	}

	/*	err = testGetBytes("InvalidData", "", 0)
		if err != nil {
			t.Error(err.Error())
		}*/
}

func TestGetFormattedPercentUptoTwoDecimals(t *testing.T) {
	var testData = []struct {
		fValue    float64
		fexpValue float64
	}{
		{10, 10},
		{12.2, 12.2},
		{15.0567, 15.05},
		{20.10, 20.10},
		{18.999999, 18.99},
	}

	for _, val := range testData {
		fTest := GetFormattedPercentUptoTwoDecimals(val.fValue)
		if fTest != val.fexpValue {
			t.Errorf("Expected value is %f,got the following value %f", val.fexpValue, fTest)
		}
	}

}

func TestGetUint64(t *testing.T) {
	var testData = []struct {
		strValue   string
		uiexpValue uint64
	}{
		{"10", 10},
		{"20", 20},
		{"0", 0},
		{"30", 30},
		{"40", 40},
	}

	for _, val := range testData {
		uiTest, err := GetUint64(val.strValue)
		if nil != err {
			t.Error(err.Error())
		}
		if uiTest != val.uiexpValue {
			t.Errorf("Expected value is %d,got the following value %d", val.uiexpValue, uiTest)
		}
	}

}
func testGetBytes(data int64, measure string, bytesValue int64) error {
	val := GetBytes(data, measure)
	if val != bytesValue {
		errMsg := fmt.Sprintf("Invalid bytes conversion for "+measure+" returned value %d", val)
		return errors.New(errMsg)
	}
	return nil
}
