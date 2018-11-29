package http

import (
	"bufio"
	"io"
	nh "net/http"

	"github.com/ContinuumLLC/platform-common-lib/src/plugin/protocol"
	"github.com/gorilla/mux"
)

// ServerHTTPFactory is HTTP implementation of protocol.ServerFactory
type ServerHTTPFactory struct{}

// GetServer returns an instance of ServerHTTP
func (cf ServerHTTPFactory) GetServer(req io.Reader, res io.Writer) protocol.Server {
	return GetServer(req, res)
}

// GetServer returns an instance of ServerHTTP
func GetServer(req io.Reader, res io.Writer) protocol.Server {
	return &serverHTTP{
		reqStream: req,
		resStream: res,
	}
}

// serverHTTP is implementation of model.Server
type serverHTTP struct {
	reqStream io.Reader
	resStream io.Writer

	routeMap  map[*mux.Route]*protocol.Route
	muxRouter *mux.Router
}

// ReceiveRequest returns a request on the stream
func (ch *serverHTTP) ReceiveRequest() (req *protocol.Request, err error) {
	reqStream := bufio.NewReader(ch.reqStream)
	httpReq, err := nh.ReadRequest(reqStream)
	if err != nil {
		return
	}

	req = newRequestFrom(httpReq)

	if ch.muxRouter != nil {
		matched := &mux.RouteMatch{}
		if ch.muxRouter.Match(httpReq, matched) {
			req.Route = ch.routeMap[matched.Route]
			req.PathParams = matched.Vars
		}
	}

	return
}

// ReceiveResponse sends response on the stream
func (ch *serverHTTP) SendResponse(res *protocol.Response) (err error) {
	return responseSerializeRaw(res, ch.resStream)
}

// RegisterRoutes registers routes with Router
func (ch *serverHTTP) RegisterRoutes(newRoutes ...*protocol.Route) {
	if ch.muxRouter == nil {
		ch.muxRouter = mux.NewRouter()
		ch.routeMap = make(map[*mux.Route]*protocol.Route, len(newRoutes))
	}
	for _, r1 := range newRoutes {
		ch.routeMap[ch.muxRouter.Path(r1.Path)] = r1
	}
}

// SetReqRespStream sets request and response stream for http server
func (ch *serverHTTP) SetReqRespStream(req io.Reader, res io.Writer) {
	ch.reqStream = req
	ch.resStream = res
}
