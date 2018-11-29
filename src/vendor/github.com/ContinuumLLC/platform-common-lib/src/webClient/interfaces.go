//Package webClient package abstracts the underlying http packages used;
//this would help abstract cross cutting concerns like Encryption, Compression etc
package webClient

import "net/http"

//ClientFactory provides implementation of ClientService
type ClientFactory interface {
	GetClientService(f HTTPClientFactory, config ClientConfig) ClientService
	GetClientServiceByType(clientType ClientType, config ClientConfig) ClientService
}

//ClientService to be implemented by HTTP Web client
type ClientService interface {
	Do(request *http.Request) (*http.Response, error)
}

//HTTPClientFactory provides implementation of HttpCommandService
type HTTPClientFactory interface {
	GetHTTPClient(config ClientConfig) HTTPClientService
}

//HTTPClientService provies methods for posting data using the net/http package
type HTTPClientService interface {
	//Post(string, string, io.Reader) (*http.Response, error)
	Do(*http.Request) (*http.Response, error)
}

//ClientConfig Http Client Configuration for HTTP connection
type ClientConfig struct {
	MaxIdleConns                int
	MaxIdleConnsPerHost         int
	IdleConnTimeoutMinute       int
	TimeoutMinute               int
	DialTimeoutSecond           int
	DialKeepAliveSecond         int
	TLSHandshakeTimeoutSecond   int
	ExpectContinueTimeoutSecond int
}

//MessageType would specify whether mesage needs to be sent to
//Broker or some other location, this way location of the server can be configured at
//single location
type MessageType uint32

const (
	//Broker Message Type
	Broker MessageType = 1
)

//HTTPMethod would specify the http method to be executed
type HTTPMethod uint32

const (
	//Post method of HTTP
	Post HTTPMethod = 1
)

//Message would contain details of the data that needs to sent as a part of
//Http method call
// type Message struct {
// 	Method      HTTPMethod
// 	ContentType string
// 	Data        []byte
// 	//MessageType MessageType
// 	URLSuffix string
// }

// Response is the HTTP response received by the web client
// type Response struct {
// 	Header           map[string][]string // Content Type, Content Length, DateTime
// 	Status           string              // e.g. "200 OK"
// 	StatusCode       int                 // e.g. 200
// 	Body             io.ReadCloser
// 	TransferEncoding []string
// }

// const (
// 	brokerURL = "http://localhost:8081"
// )

//Error Codes
const (
	ErrorInvalidHTTPMethod  = "ErrInvalidHTTPMethod"
	ErrorEmptyContentType   = "ErrEmptyContentType"
	ErrorNilURL             = "ErrorNilURL"
	ErrorBlankHttpMethod    = "BlankHttpMethod"
	ErrorNilData            = "ErrNilData"
	ErrorInvalidMessageType = "ErrInvalidMessageType"
	ErrorEmptyURLSuffix     = "ErrEmptyURLSuffix"
	ErrorInvalidURLSuffix   = "ErrInvalidURLSuffix"

	BasicClient ClientType = 10
	TlsClient   ClientType = 20
)

type ClientType int
