package http

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"strings"

	"github.com/ContinuumLLC/platform-common-lib/src/exception"
	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

func TestResponseError(t *testing.T) {
	ser := ClientHTTPFactory{}.GetResponseSerializer()
	buf := &bytes.Buffer{}

	resp := protocol.NewResponse()
	data := &bytes.Buffer{}
	data.WriteString("This string is Base64 Encoded")
	resp.Body = data
	resp.SetError(protocol.PathNotFound, exception.New("SomeError", nil))
	ser.Serialize(resp, buf)

	// pser is resp deserialized back
	pser, _ := ser.Deserialize(buf)

	fmt.Println(pser)
	if pser.Status != protocol.PathNotFound {
		t.Error("Status Incorrect")
		return
	}
	checkHeaderForValue(t, protocol.HdrErrorCode, pser, true, "SomeError")
	checkHeaderForValue(t, protocol.HdrContentType, pser, true, "text/plain")

	buf = &bytes.Buffer{}
	io.Copy(buf, pser.Body)
	pserStr := buf.String()
	if !strings.HasSuffix(pserStr, "VGhpcyBzdHJpbmcgaXMgQmFzZTY0IEVuY29kZWQ=") {
		t.Error("Base64 Encoded String Missing")
	}
}

func checkHeaderForValue(t *testing.T, header protocol.HeaderKey, resp *protocol.Response, checkValue bool, value string) {
	valArr, ok := resp.Headers[header]
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
