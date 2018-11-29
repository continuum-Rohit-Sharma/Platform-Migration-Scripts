package http

import (
	"bufio"
	"io"
	nh "net/http"
	nht "net/http/httptest"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
)

func responseSerializeRaw(res *protocol.Response, dst io.Writer) (err error) {
	httpRes := nht.NewRecorder()

	for k, v := range res.Headers {
		httpRes.HeaderMap[getHTTPHeader(k)] = v
	}

	httpRes.HeaderMap.Set(getHTTPHeader(protocol.HdrTransactionID), res.TransactionID)

	if res.Status != protocol.Ok {
		httpRes.WriteHeader(int(res.Status))
	}

	if res.Body != nil {
		_, err = io.Copy(httpRes, res.Body)
		if err != nil {
			return err
		}
	}

	err = httpRes.Result().Write(dst)

	return
}

func responseDeserializeRaw(src io.Reader) (res *protocol.Response, err error) {
	httpRes, err := nh.ReadResponse(bufio.NewReader(src), nil)
	if err != nil {
		return
	}
	res, err = newResponseFrom(httpRes)
	res.TransactionID = httpRes.Header.Get(getHTTPHeader(protocol.HdrTransactionID))
	return
}
