package web

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"

	"github.com/gorilla/mux"
)

const (
	cProxyHeader = "X-FORWARDED-FOR"
)

// muxConfig structure is Mux Adapter for Server Interface
type muxConfig struct {
	serverCfg *ServerConfig
	router    *mux.Router
}

// SetupRoutes implementation of Server interface for muxConfig
func (mcfg muxConfig) SetupRoutes(routes []*RouteConfig) {
	for _, route := range routes {
		aHandler := muxRouteHandler{
			route: route,
		}
		mcfg.router.HandleFunc(mcfg.serverCfg.URLPathPrefix+route.URLPathSuffix,
			aHandler.handleFunc)
	}
}

//GetRouter returns the mux router
func (mcfg muxConfig) GetRouter() *mux.Router {
	return mcfg.router
}

// ListenAndServe implementation of Server interface for muxConfig
func (mcfg muxConfig) ListenAndServe() {
	srv := mcfg.createServerInstance()
	srv.ListenAndServe()
}

//HTTP2ListenAndServeTLS listen as HTTP2 (which is only in TLS)
func (mcfg muxConfig) HTTP2ListenAndServeTLS() error {
	srv := mcfg.createServerInstance()
	http2.ConfigureServer(srv, &http2.Server{
		IdleTimeout:          time.Duration(mcfg.serverCfg.IdleTimeoutMinute) * time.Minute,
		MaxHandlers:          mcfg.serverCfg.MaxHandlers,
		MaxConcurrentStreams: mcfg.serverCfg.MaxConcurrentStreams,
	})
	return srv.ListenAndServeTLS(mcfg.serverCfg.CertificateFile, mcfg.serverCfg.CertificateKeyFile)
}

func (mcfg muxConfig) createServerInstance() *http.Server {
	return &http.Server{
		Addr:         mcfg.serverCfg.ListenURL,
		Handler:      mcfg.router,
		ReadTimeout:  time.Duration(mcfg.serverCfg.ReadTimeoutMinute) * time.Minute,
		WriteTimeout: time.Duration(mcfg.serverCfg.WriteTimeoutMinute) * time.Minute,
	}
}

type muxRouteHandler struct {
	route *RouteConfig
}

type muxRequestContext struct {
	response      http.ResponseWriter
	request       *http.Request
	vars          map[string]string
	varsResolved  bool
	dcDateTimeUTC time.Time
}

func (ctx muxRequestContext) GetRequest() *http.Request {
	return ctx.request
}

func (ctx muxRequestContext) GetResponse() http.ResponseWriter {
	return ctx.response
}

func (ctx muxRequestContext) GetVars() map[string]string {
	if !ctx.varsResolved {
		ctx.varsResolved = true
		ctx.vars = mux.Vars(ctx.request)
	}
	return ctx.vars
}

func (ctx muxRequestContext) GetData() (data []byte, err error) {
	return ioutil.ReadAll(ctx.GetRequest().Body)
}

func (ctx muxRequestContext) GetRequestDcDateTimeUTC() time.Time {
	return ctx.dcDateTimeUTC
}

func (aHandler muxRouteHandler) handleFunc(w http.ResponseWriter, r *http.Request) {
	ctx := &muxRequestContext{
		response:      w,
		request:       r,
		dcDateTimeUTC: time.Now().UTC(),
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Connection", "close")

	switch r.Method {
	case http.MethodGet:
		aHandler.route.Res.Get(ctx)
	case http.MethodPost:
		aHandler.route.Res.Post(ctx)
	case http.MethodPut:
		aHandler.route.Res.Put(ctx)
	case http.MethodDelete:
		aHandler.route.Res.Delete(ctx)
	default:
		aHandler.route.Res.Others(ctx)
	}

	return
}

func (ctx muxRequestContext) GetRemoteAddr() (string, error) {
	ipAddress := "0.0.0.0"
	remoteProxy := ctx.request.Header.Get(cProxyHeader)
	remoteHostPort := ctx.request.RemoteAddr
	//If remoteProxy is set it means endpoint was hidden behind proxy, hence get the real IP from x-forwarded-for header.
	//else try to get it directly from RemoteAddress attribute of HTTP request
	if len(remoteProxy) > 0 {
		//X-Forwarded-For: client, proxy1, proxy2
		// where the value is a comma+space separated list of IP addresses
		ips := strings.Split(remoteProxy, ", ")
		if len(ips) > 0 {
			//TODO verify if first IP is the real IP.
			ip := net.ParseIP(ips[0])
			if ip != nil {
				ipAddress = ip.String()
			}
		}
	} else if len(remoteHostPort) > 0 {
		host, _, err := net.SplitHostPort(remoteHostPort)
		if err != nil {
			return ipAddress, err
		}
		ip := net.ParseIP(host)
		if ip != nil {
			ipAddress = ip.String()
		}
	}

	//TODO Need to check how IPv6 values for remoteIP, remoteProxy will look like
	//if with bracket, [2001:db8:85a3:8d3:1319:8a2e:370:7348]:443, above solution works

	//if without brackets, ipv6= 2001:db8:85a3:8d3:1319:8a2e:370:7348:443
	/*var port string
	ind := strings.LastIndex(header.RemoteAddr, ":")

	if ind > 0 {
		port = header.RemoteAddr[ind+1:]
		ip := net.ParseIP(header.RemoteAddr[:ind])
		if ip != nil {
			assetColl.Message.RemoteAddress = ip.String()
		}
	}*/

	return ipAddress, nil
}
