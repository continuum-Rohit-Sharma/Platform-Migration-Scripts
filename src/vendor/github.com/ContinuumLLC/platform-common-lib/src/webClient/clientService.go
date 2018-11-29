package webClient

import (
	"errors"
	"net/http"
)

//Client implementation using go net/http package
type clientServiceImpl struct {
	httpClientFact HTTPClientFactory
	config         ClientConfig
}

func (cs clientServiceImpl) Do(request *http.Request) (*http.Response, error) {
	client := cs.httpClientFact.GetHTTPClient(cs.config)
	return client.Do(request)
}

func validateRequest(request *http.Request) error {
	if request.Header.Get("Content-Type") == "" {
		return errors.New(ErrorEmptyContentType)
	}
	if request.URL == nil {
		return errors.New(ErrorNilURL)
	}
	if request.Method == "" {
		return errors.New(ErrorBlankHttpMethod)
	}
	return nil
}
