package webClient

import (
	"net/http"
	"net/url"
	"testing"
)

func TestValidateRequestNilNoContent(t *testing.T) {
	request := &http.Request{}
	err := validateRequest(request)
	if err == nil {
		t.Errorf("Expected error %s, not returned", err)
		return
	}
	if err.Error() != ErrorEmptyContentType {
		t.Errorf("Unexpected error returned, Expected Error: %s, Returned Error: %v", ErrorEmptyContentType, err)
		return
	}
}

func TestValidateRequestNilUrl(t *testing.T) {
	request := &http.Request{}
	request.Header = http.Header{}
	request.Header.Set("Content-Type", "text")
	err := validateRequest(request)
	if err == nil {
		t.Errorf("Expected error %s, not returned", err)
		return
	}
	if err.Error() != ErrorNilURL {
		t.Errorf("Unexpected error returned, Expected Error: %s, Returned Error: %v", ErrorNilURL, err)
		return
	}
}

func TestValidateRequestNilMethod(t *testing.T) {
	request := &http.Request{
		URL: &url.URL{},
	}
	request.Header = http.Header{}
	request.Header.Set("Content-Type", "text")
	err := validateRequest(request)
	if err == nil {
		t.Errorf("Expected error %s, not returned", err)
		return
	}
	if err.Error() != ErrorBlankHttpMethod {
		t.Errorf("Unexpected error returned, Expected Error: %s, Returned Error: %v", ErrorBlankHttpMethod, err)
		return
	}
}

func TestValidateRequestSuccess(t *testing.T) {
	request := &http.Request{
		URL:    &url.URL{},
		Method: "GET",
	}
	request.Header = http.Header{}
	request.Header.Set("Content-Type", "text")
	err := validateRequest(request)
	if err != nil {
		t.Errorf("Expected nil error, %s returned", err)
		return
	}
}

func TestGetClientService(t *testing.T) {
	clientFact := ClientFactoryImpl{}
	httpClientFact := HTTPClientFactoryImpl{}
	clientService := clientFact.GetClientService(httpClientFact, ClientConfig{})
	if clientService == nil {
		t.Error("Expected ClientFactory returned nil")
	}
}

func TestGetClientServiceDo(t *testing.T) {
	clientFact := ClientFactoryImpl{}
	httpClientFact := HTTPClientFactoryImpl{}
	clientService := clientFact.GetClientService(httpClientFact, ClientConfig{})
	request, _ := http.NewRequest("GET", "http://local", nil)
	_, err := clientService.Do(request)
	if err == nil {
		t.Error("Expected Error returned nil")
	}
}

func TestGetTLSClientServiceDo(t *testing.T) {
	clientFact := ClientFactoryImpl{}
	clientService := clientFact.GetClientServiceByType(TlsClient, ClientConfig{})
	request, _ := http.NewRequest("GET", "http://local", nil)
	_, err := clientService.Do(request)
	if err == nil {
		t.Error("Expected Error returned nil")
	}
}
