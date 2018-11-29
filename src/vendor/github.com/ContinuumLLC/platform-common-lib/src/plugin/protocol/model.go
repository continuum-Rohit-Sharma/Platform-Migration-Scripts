package protocol

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
)

// Client is the protocol initiater, Agent Core
type Client interface {
	SendRequest(req *Request) error
	ReceiveResponse() (*Response, error)
}

// Server is the protocol responder, Agent Plugin
type Server interface {
	ReceiveRequest() (*Request, error)
	SendResponse(res *Response) error
	RegisterRoutes(routes ...*Route)
	SetReqRespStream(req io.Reader, res io.Writer)
}

// ClientFactory is factory of Client objects
type ClientFactory interface {
	GetClient(req io.Writer, res io.Reader) Client
	GetResponseSerializer() ResponseSerializer
}

// ServerFactory is factory of server objects
type ServerFactory interface {
	GetServer(req io.Reader, res io.Writer) Server
}

// Request is request sent from core to plugin
type Request struct {
	Path          string
	Headers       Headers
	Body          io.Reader
	Params        Parameters
	Route         *Route
	PathParams    map[string]string
	TransactionID string
}

// Response is response returned from plugin to core for sent request
type Response struct {
	Headers       Headers
	Body          io.Reader
	Status        ResponseStatus
	TransactionID string
}

// SetError sets up error information on the response object.
// Sets status code, assigns error header and prepends rootCause stackTrace to body
func (resp *Response) SetError(code ResponseStatus, rootCause error) {
	resp.Status = code
	if rootCause != nil {
		// set header
		resp.Headers.SetKeyValue(HdrErrorCode, rootCause.Error())
		resp.Headers.SetKeyValue(HdrContentType, "text/plain")
		// append stacktrace
		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, "%+v", rootCause)
		if resp.Body != nil {
			fmt.Fprintln(buf, "\nContent (Base64 Encoded):")
			enc := base64.NewEncoder(base64.StdEncoding, buf)
			io.Copy(enc, resp.Body)
			enc.Close()
		}
		resp.Body = buf
	} else {
		resp.Headers.SetKeyValue(HdrErrorCode, string(HdrErrorCode))
	}
}

// NewRequest constructs and returns a Request
func NewRequest() (req *Request) {
	req = &Request{}
	req.Headers = make(Headers)
	req.Params = make(Parameters)
	return
}

// NewResponse constructs a response object
func NewResponse() (res *Response) {
	res = &Response{
		Status: Ok,
	}
	res.Headers = make(Headers)
	return
}

// HandleRoute is plugin path route handler
type HandleRoute func(req *Request) (res *Response, err error)

// Route is a Plugin Path configuration
type Route struct {
	Path   string
	Handle HandleRoute
}
