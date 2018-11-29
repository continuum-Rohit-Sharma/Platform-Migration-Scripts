package http

import (
	"bytes"
	"os"
	"testing"

	"io"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

func TestProto(t *testing.T) {
	// fmt.Println("Hello World")
	c, s, ios := setupProto(false, false)
	sendRequest(c)
	receiveRequest(t, s)
	sendResponse(t, s)
	ios.resW.Close()
	receiveResponse(t, c)
}

func TestProtoErrRecvResp(t *testing.T) {
	// fmt.Println("Hello World")
	c, s, _ := setupProto(false, true)
	sendRequest(c)
	receiveRequest(t, s)
	sendResponse(t, s)
	receiveResponseErr(t, c)
}

func BenchmarkProto(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c, s, _ := setupProto(false, false)
		sendRequest(c)
		receiveRequest(b, s)
		sendResponse(b, s)
		receiveResponse(b, c)
	}
}

type iostreams struct {
	reqW io.WriteCloser
	reqR io.ReadCloser
	resW io.WriteCloser
	resR io.ReadCloser
}

// func TestProtoErrReqRecv(t *testing.T) {
// 	c, s := setupProto(true, false)
// 	sendRequest(c)
// 	receiveRequestErr(t, s)
// }

// func TestProtoSendRespErr(t *testing.T) {
// 	c, s := setupProto(false, false)
// 	sendRequest(c)
// 	receiveRequest(t, s)
// 	sendResponseErr(t, s)
// }

func setupProto(badReq bool, badRes bool) (c protocol.Client, s protocol.Server, ios *iostreams) {
	reqr, reqw, _ := os.Pipe()
	resr, resw, _ := os.Pipe()
	ios = &iostreams{reqw, reqr, resw, resr}

	c = ClientHTTPFactory{}.GetClient(reqw, resr)
	s = ServerHTTPFactory{}.GetServer(reqr, resw)

	if badReq {
		reqw.WriteString("BadRequest")
	}
	if badRes {
		resw.WriteString("BadResponse")
	}

	return
}

func sendRequest(c protocol.Client) {
	req := protocol.NewRequest()
	req.Path = "/test/path"
	req.Headers.SetKeyValue(protocol.HdrUserAgent, "UnitTest")
	c.SendRequest(req)
}

func receiveRequest(t testing.TB, s protocol.Server) {
	sreq, _ := s.ReceiveRequest()
	if sreq.Path != "/test/path" {
		t.Error("Unexpected Path")
		return
	}
	if sreq.Headers.GetKeyValue(protocol.HdrUserAgent) != "UnitTest" {
		t.Error("Unexpected Header")
		return
	}
}

func receiveRequestErr(t testing.TB, s protocol.Server) {
	req, err := s.ReceiveRequest()
	if err == nil {
		t.Error("Expected Error not Received")
		t.Error(req)
	}
}

func sendResponse(t testing.TB, s protocol.Server) {
	res := protocol.NewResponse()
	res.Headers.SetKeyValue(protocol.HdrContentType, "text/json")
	res.Body = bytes.NewBufferString("Hello World!")
	err := s.SendResponse(res)
	if err != nil {
		t.Error("Unexpected Error: " + err.Error())
	}
}

func sendResponseErr(t testing.TB, s protocol.Server) {
	res := protocol.NewResponse()
	res.Headers.SetKeyValue(protocol.HdrContentType, "text/json")
	err := s.SendResponse(res)
	if err != nil {
		t.Error("Unexpected Error: " + err.Error())
	}
}

func receiveResponse(t testing.TB, c protocol.Client) {
	cres, _ := c.ReceiveResponse()

	if cres.Headers.GetKeyValue(protocol.HdrContentType) != "text/json" {
		t.Error("Unexpected Header in Response: ")
		// t.Error(res)
		t.Error(cres)
		return
	}
}

func receiveResponseErr(t testing.TB, c protocol.Client) {
	_, err := c.ReceiveResponse()
	if err == nil {
		t.Error("Expected Error not Received")
	}
}
