package web

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ServerConfig stores the configuration of the web server.
type ServerConfig struct {
	URLPathPrefix        string
	ListenURL            string
	CertificateFile      string
	CertificateKeyFile   string
	ReadTimeoutMinute    int
	WriteTimeoutMinute   int
	IdleTimeoutMinute    int
	MaxHandlers          int
	MaxConcurrentStreams uint32
}

//Server interface sets up routes, handlers and listening.
type Server interface {
	SetupRoutes(res []*RouteConfig)
	ListenAndServe()
	GetRouter() *mux.Router
	HTTP2ListenAndServeTLS() error
}

//RequestContext that will be provided by router to request handlers
type RequestContext interface {
	GetResponse() http.ResponseWriter
	GetRequest() *http.Request
	GetVars() map[string]string
	GetRequestDcDateTimeUTC() time.Time
	GetData() (data []byte, err error)
	GetRemoteAddr() (string, error)
}

//ServerFactory interface to for a Factory impmementation
type ServerFactory interface {
	GetServer(cfg *ServerConfig) Server
}

//ServerFactoryImpl A factory implementation for the HTTP server creation
type ServerFactoryImpl struct{}

//GetServer implements Server interface
func (ServerFactoryImpl) GetServer(cfg *ServerConfig) Server {
	mcfg := muxConfig{
		serverCfg: cfg,
		router:    mux.NewRouter(),
	}
	return &mcfg
}
