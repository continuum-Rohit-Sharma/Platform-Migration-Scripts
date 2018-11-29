package protocol

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest()
	if req.Headers == nil {
		t.Error("Headers is nil")
		return
	}

	if req.Params == nil {
		t.Error("Params is nil")
		return
	}
}

func TestNewResponse(t *testing.T) {
	res := NewResponse()
	if res.Headers == nil {
		t.Error("Headers is nil")
		return
	}
}

func TestResponseSimpleError(t *testing.T) {
	resp := setupResponse(true)
	oldBody := resp.Body
	resp.SetError(PathNotFound, errors.New("SampleError"))

	if oldBody == resp.Body {
		t.Error("Response.Body was not altered")
		return
	}

	buf := &bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	actualStr := buf.String()
	if !strings.HasPrefix(actualStr, "SampleError") {
		t.Error("Error Code Mismatch [" + actualStr + "]")
		return
	}
}

func TestResponseNilError(t *testing.T) {
	resp := setupResponse(true)
	oldBody := resp.Body
	resp.SetError(PathNotFound, nil)

	if oldBody != resp.Body {
		t.Error("Response.Body was altered")
		return
	}
	checkErrorHeaderForValue(t, resp, true, string(HdrErrorCode))
}

func checkErrorHeaderForValue(t *testing.T, resp *Response, checkValue bool, value string) {
	valArr, ok := resp.Headers[HdrErrorCode]
	if !ok {
		t.Error("Error header missing")
		return
	}

	if checkValue {
		if len(valArr) == 0 || valArr[0] != value {
			t.Error("Error Value Different")
			if len(valArr) != 0 {
				t.Error(valArr)
			}
			return
		}
	}
}

func setupResponse(body bool) *Response {
	resp := NewResponse()

	if body {
		resp.Body = &bytes.Buffer{}
	}

	return resp
}
