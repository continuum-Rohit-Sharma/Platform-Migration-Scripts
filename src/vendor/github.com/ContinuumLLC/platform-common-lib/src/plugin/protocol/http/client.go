package http

import (
	"io"
	nh "net/http"
	nht "net/http/httptest"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

// ClientHTTPFactory is HTTP implementation of protocol.ClientFactory
type ClientHTTPFactory struct{}

// GetClient returns an instance of ClientHTTP
func (cf ClientHTTPFactory) GetClient(req io.Writer, res io.Reader) protocol.Client {
	return &clientHTTP{
		reqStream: req,
		resStream: res,
	}
}

// GetResponseSerializer ...
func (cf ClientHTTPFactory) GetResponseSerializer() protocol.ResponseSerializer {
	return &responseSerializerImpl{}
}

// clientHTTP is implementation of model.Client
type clientHTTP struct {
	reqStream io.Writer
	resStream io.Reader
}

// SendRequest sends a request on the stream
func (ch *clientHTTP) SendRequest(req *protocol.Request) (err error) {
	httpReq := ch.CreateHTTPRequest(req)
	err = httpReq.Write(ch.reqStream)
	return err
}

// CreateHTTPRequest ...
func (ch *clientHTTP) CreateHTTPRequest(req *protocol.Request) (httpReq *nh.Request) {
	httpReq = nht.NewRequest(nh.MethodPost, req.Path, req.Body)
	httpReq.Host = ""
	for k, v := range req.Headers {
		httpReq.Header[getHTTPHeader(k)] = v
	}
	return
}

// ReceiveResponse expects a response from the stream
func (ch *clientHTTP) ReceiveResponse() (res *protocol.Response, err error) {
	return responseDeserializeRaw(ch.resStream)
}
