package http

import (
	"bytes"
	"io"
	nh "net/http"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

func newRequestFrom(httpReq *nh.Request) (req *protocol.Request) {
	req = protocol.NewRequest()
	req.Path = httpReq.URL.String()
	req.Body = httpReq.Body
	for k, v := range httpReq.Header {
		req.Headers[getProtocolHeader(k)] = v
	}
	return
}

func newResponseFrom(httpRes *nh.Response) (res *protocol.Response, err error) {
	res = protocol.NewResponse()
	if httpRes.Body != nil {
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, httpRes.Body)
		if err != nil {
			return
		}
		res.Body = buf
	}
	res.Status = protocol.ResponseStatus(httpRes.StatusCode)
	for k, v := range httpRes.Header {
		res.Headers[getProtocolHeader(k)] = v
	}
	return
}
