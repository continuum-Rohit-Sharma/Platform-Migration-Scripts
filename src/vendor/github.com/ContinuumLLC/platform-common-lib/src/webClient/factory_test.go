package webClient

import (
	"errors"
	"testing"
)

func TestCheckOfflineTimeout(t *testing.T) {
	errorString := errors.New("connection timeout")
	err := checkOffline(errorString)

	if err.Error() != ErrorClientOffline {
		t.Errorf("Expected error: %v, Received Error: %v", ErrorClientOffline, err)
	}
}

func TestCheckOfflineError(t *testing.T) {
	errorString := errors.New("TestError")
	err := checkOffline(errorString)

	if err.Error() != "TestError" {
		t.Errorf("Expected error: %v, Received Error: %v", "TestError", err)
	}
}

func TestCheckOfflineConnectionRefused(t *testing.T) {
	errorString := errors.New("dial tcp [::1]:8081: getsockopt: connection refused")
	err := checkOffline(errorString)

	if err.Error() != ErrorClientOffline {
		t.Errorf("Expected error: %v, Received Error: %v", ErrorClientOffline, err)
	}
}

func TestGetClientServiceByTlsType(t *testing.T) {
	clientFact := ClientFactoryImpl{}
	clientService := clientFact.GetClientServiceByType(TlsClient, ClientConfig{})
	_, ok := clientService.(*tlsClientService)
	if !ok {
		t.Errorf("Expected tlsClientService returned %v", clientService)
	}
}

func TestGetClientServiceByBasicType(t *testing.T) {
	clientFact := ClientFactoryImpl{}
	clientService := clientFact.GetClientServiceByType(BasicClient, ClientConfig{})
	_, ok := clientService.(httpClientServiceImpl)
	if !ok {
		t.Errorf("Expected httpClientServiceImpl returned %v", clientService)
	}
}
