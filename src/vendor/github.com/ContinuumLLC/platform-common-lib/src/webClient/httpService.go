package webClient

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

//httpClientServiceImpl implements HTTPCommandService
type httpClientServiceImpl struct {
	config     ClientConfig
	httpClient *http.Client
}

//Do sends Request to the Server
func (hc httpClientServiceImpl) Do(request *http.Request) (*http.Response, error) {
	if hc.httpClient == nil {
		hc.httpClient = createClient(hc.config, nil)
	}

	resp, err := hc.httpClient.Do(request)
	if err != nil {
		err = checkOffline(err)
	}
	return resp, err
}

func createClient(config ClientConfig, tlsConfig *tls.Config) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(config.DialTimeoutSecond) * time.Second,
			KeepAlive: time.Duration(config.DialKeepAliveSecond) * time.Second,
		}).DialContext,
		MaxIdleConns:          config.MaxIdleConns,
		IdleConnTimeout:       time.Duration(config.IdleConnTimeoutMinute) * time.Minute,
		TLSHandshakeTimeout:   time.Duration(config.TLSHandshakeTimeoutSecond) * time.Second,
		ExpectContinueTimeout: time.Duration(config.ExpectContinueTimeoutSecond) * time.Second,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
	}
	client := &http.Client{
		Timeout:   time.Duration(config.TimeoutMinute) * time.Minute,
		Transport: transport,
	}
	return client
}
