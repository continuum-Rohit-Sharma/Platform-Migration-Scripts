package http

import (
	"bytes"
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

func TestSer(t *testing.T) {
	res := protocol.NewResponse()
	res.Headers.SetKeyValue(protocol.HeaderKey("Customheader"), "CustomValue")
	buf := new(bytes.Buffer)

	ser := ClientHTTPFactory{}.GetResponseSerializer()
	err := ser.Serialize(res, buf)
	if err != nil {
		t.Error("Unexpected Error in Response Serializer")
		t.Error(err)
	}

	res, err = ser.Deserialize(buf)
	if err != nil {
		t.Error("Unexpected Error in Response DeSerializer")
		t.Error(err)
	}

	if res.Headers.GetKeyValue(protocol.HeaderKey("Customheader")) != "CustomValue" {
		t.Error("Headers Lost:")
		t.Error(res)
	}
}
