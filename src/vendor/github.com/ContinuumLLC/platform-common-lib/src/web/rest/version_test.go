package rest

import (
	"net/http"
	"testing"
)

func TestRegistryVersion(t *testing.T) {
	oldVersionData := versionData
	defer func() { versionData = oldVersionData }()

	RegistryVersion(&Version{Type: versionType})
	if versionData.(*Version).Type != versionType {
		t.Errorf("expected Type %s, but got %s", versionType, versionData.(*Version).Type)
	}
}

func TestHandlerVersion_OK(t *testing.T) {
	oldVersionData := versionData
	versionData = &Version{}
	defer func() { versionData = oldVersionData }()
	versionLastTimeSuccess = 0

	mock := &mockResponseWriter{dataHeader: http.Header{}}
	HandlerVersion(mock, nil)
	if mock.dataWriteHeader != http.StatusOK {
		t.Errorf("expected code %d, but got %d", http.StatusOK, mock.dataWriteHeader)
	}
}
