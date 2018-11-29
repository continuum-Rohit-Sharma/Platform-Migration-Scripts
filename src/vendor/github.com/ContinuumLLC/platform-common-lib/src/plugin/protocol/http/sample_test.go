package http

// dummy code used for draft specification document

// import (
// 	"encoding/json"
// 	"io"

// 	"github.com/continuumLLC/platform-common-lib/src/plugin/protocol"
// )

// func sampleServer(reqR io.Reader, resW io.Writer, sf protocol.ServerFactory) {

// 	// create a server
// 	s := sf.NewServer(reqR, resW)

// 	// read request
// 	req, err := s.ReceiveRequest()

// 	// check headers and body, parse request
// 	reqObj := parseRequestToObject(req)

// 	// call business layer to process request
// 	respObj := businessLayerCall(reqObj)

// 	// create and populate protocol response
// 	resp := protocol.NewResponse()
// 	resp.Body := json.Marshal(respObj)
// 	resp.Status = protocol.Ok
// 	resp.Headers.SetKeyValue(protocol.HdrContentType, "text/json")

// 	// send response
// 	err = s.SendResponse(resp)
// }

// func sampleClient(reqW io.Writer, resR io.Reader, cf protocol.ClientFactory) {

// 	// create a client
// 	c := cf.NewClient(reqW, resR)

// 	// create request
// 	req := protocol.NewRequest()
// 	req.Path = "/plugin/path"
// 	req.Headers.SetKeyValue(protocol.HdrUserAgent, "Agent/1.1/Ubuntu/14.04")

// 	// if there's a payload, set the content type, and set the Body
// 	req.Headers.SetKeyValue(protocol.HdrContentType, "text/json")
// 	req.Body = PayloadReader()

// 	// send the request
// 	err := c.SendRequest(req)

// 	// receive the response
// 	resp, err := c.ReceiveResponse()

// 	// check the response status
// 	if resp.Status == protocol.Ok {
// 		// do something
// 	}

// 	contentType := resp.Headers.GetKeyValue(protocol.HdrContentType)
// 	if contentType != nil && contentType == "text/json" {
// 		// read body and handle content
// 		responsePayload := json.Unmarshal(resp.Body)
// 	}
// }

// func sampleRoutedServer(reqR io.Reader, resW io.Writer, sf protocol.ServerFactory) {

// 	// create a server
// 	s := sf.NewServer(reqR, resW)

// 	// register routes
// 	s.RegisterRoutes(&protocol.Route{Path: "/performance/memory", Handle: memoryHandler},
// 		&protocol.Route{Path: "/performance/network", Handle: networkHandler})

// 	// read request
// 	req, err := s.ReceiveRequest()

// 	// retrieve response
// 	var resp protocol.Response
// 	if req.Route != nil {
// 		resp, _ = req.Route.Handle(req)
// 	} else {
// 		// handle error for unregistered/invalid route
// 	}

// 	// send response
// 	err = s.SendResponse(resp)
// }

// func memoryHandler(*protocol.Request) (*protocol.Response, error) {
// 	// create and populate protocol response
// 	resp := protocol.NewResponse()
// 	resp.Body := json.Marshal("memory")
// 	resp.Status = protocol.Ok
// 	resp.Headers.SetKeyValue(protocol.HdrContentType, "text/json")

// 	return resp, nil
// }

// func ProcessProcessor(*protocol.Request) (*protocol.Response, error) {
// 	// create and populate protocol response
// 	resp := protocol.NewResponse()
// 	resp.Body := json.Marshal("network")
// 	resp.Status = protocol.Ok
// 	resp.Headers.SetKeyValue(protocol.HdrContentType, "text/json")

// 	return resp, nil
// }
