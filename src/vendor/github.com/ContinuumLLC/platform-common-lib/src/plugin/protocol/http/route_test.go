package http

import (
	"testing"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

func BenchmarkRoute(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c, s, ios := setupProto(false, false)
		registerRoutes(s)
		sendRoutedRequest(c)
		receiveRoutedRequest(b, s)
		sendResponse(b, s)
		ios.resW.Close()
		receiveResponse(b, c)
	}
}

func TestProtoRoute(t *testing.T) {
	// fmt.Println("Hello World")
	c, s, ios := setupProto(false, false)
	registerRoutes(s)
	sendRoutedRequest(c)
	receiveRoutedRequest(t, s)
	sendResponse(t, s)
	ios.resW.Close()
	receiveResponse(t, c)
}

func registerRoutes(s protocol.Server) {
	s.RegisterRoutes(&protocol.Route{Path: "/performanceMemory", Handle: nil},
		&protocol.Route{Path: "/performanceNetwork", Handle: nil})
}

func sendRoutedRequest(c protocol.Client) {
	req := protocol.NewRequest()
	req.Path = "/performanceNetwork"
	req.Headers.SetKeyValue(protocol.HdrUserAgent, "UnitTest")
	c.SendRequest(req)
}

func receiveRoutedRequest(t testing.TB, s protocol.Server) {
	req, _ := s.ReceiveRequest()
	pr := req.Route
	if pr == nil {
		t.Error("Did not find Protocol Route")
		return
	}
	if pr.Path != "/performanceNetwork" {
		t.Error("Matched a different Route: " + pr.Path)
		return
	}
}

func TestProtoRouteComplex(t *testing.T) {
	// fmt.Println("Hello World")
	c, s, _ := setupProto(false, false)
	registerRoutesComplex(s)
	sendRoutedRequestComplex(c)
	receiveRoutedRequestComplex(t, s)
	// sendResponse(t, s)
	// receiveResponse(t, c)
}

func registerRoutesComplex(s protocol.Server) {
	s.RegisterRoutes(&protocol.Route{Path: "/performance/processor/{procIndex}/core/{coreIndex}", Handle: nil})
}

func sendRoutedRequestComplex(c protocol.Client) {
	req := protocol.NewRequest()
	req.Path = "/performance/processor/0/core/1"
	req.Headers.SetKeyValue(protocol.HdrUserAgent, "UnitTest")
	c.SendRequest(req)
}

func receiveRoutedRequestComplex(t testing.TB, s protocol.Server) {
	req, _ := s.ReceiveRequest()
	pr := req.Route
	if pr == nil {
		t.Error("Did not find Protocol Route")
		return
	}
	if req.PathParams == nil {
		t.Error("Did not find Path Params")
		return
	}
	if req.PathParams["procIndex"] != "0" || req.PathParams["coreIndex"] != "1" {
		t.Error("Did not get variables: ")
		t.Error(req.PathParams)
		return
	}
}
