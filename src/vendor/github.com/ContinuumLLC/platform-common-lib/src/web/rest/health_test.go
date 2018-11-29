package rest

import (
	"fmt"
	"net/http"
	"testing"
)

const listenURL = ":12124"

func TestRegistryHealth(t *testing.T) {
	const status = "testing"

	oldHealthData := healthData
	defer func() { healthData = oldHealthData }()

	RegistryHealth(&Health{Status: status})
	if healthData.(*Health).Status != status {
		t.Errorf("expected status %s, but got %s", status, healthData.(*Health).Status)
	}
}

func TestHealth(t *testing.T) {
	h := &Health{ListenURL: listenURL}
	if err := h.health(); err != nil {
		t.Error(err)
	}

	if h.Type != healthType {
		t.Errorf("expected Type %s, but got %s", healthType, h.Type)
	}

	if h.LastStartTimeUTC != TimeStampUTC {
		t.Errorf("expected LastStartTimeUTC %s, but got %s", TimeStampUTC, h.LastStartTimeUTC)
	}

	if len(h.OutboundConnectionStatus) != 0 {
		t.Error("OutboundConnectionStatus must be empty")
	}

	if h.Status != healthRunningStatus {
		t.Errorf("expected Status %s, but got %s", healthRunningStatus, h.Status)
	}

	h.ListenURL = ""
	if err := h.health(); err == nil {
		t.Error("Error can not be <nil>")
	}
}

func TestGetHealthStatus(t *testing.T) {
	var conns = []OutboundConnectionStatus{
		{ConnectionStatus: ConnectionStatusUnavailable},
	}
	if status := GetHealthStatus(conns); status != healthDegradedStatus {
		t.Errorf("expected Status %s, but got %s", healthDegradedStatus, status)
	}
}

type mockStatuser struct{}

func (mockStatuser) Status(status OutboundConnectionStatus) *OutboundConnectionStatus {
	return &status
}

func TestGetOutboundConnectionStatus(t *testing.T) {
	const serviceName = "testServiceName"
	var methods = []Statuser{mockStatuser{}}
	connections := GetOutboundConnectionStatus(methods, serviceName)
	name := fmt.Sprintf("%s-%s", serviceName, dbNameSuffix)
	if connections[0].Name != name {
		t.Errorf("expected Name %q, but got %q", name, connections[0].Name)
	}
}

func TestHandlerHealth_OK(t *testing.T) {
	oldHealthData := healthData
	healthData = &Health{ListenURL: listenURL}
	defer func() { healthData = oldHealthData }()
	healthLastTimeSuccess = 0

	mock := &mockResponseWriter{dataHeader: http.Header{}}
	HandlerHealth(mock, nil)
	if mock.dataWriteHeader != http.StatusOK {
		t.Errorf("expected code %d, but got %d", http.StatusOK, mock.dataWriteHeader)
	}
}

func TestHandlerHealth_InternalServerError(t *testing.T) {
	oldHealthData := healthData
	healthData = &Health{}
	defer func() { healthData = oldHealthData }()
	healthLastTimeSuccess = 0

	mock := &mockResponseWriter{dataHeader: http.Header{}}
	HandlerHealth(mock, nil)
	if mock.dataWriteHeader != http.StatusInternalServerError {
		t.Errorf("expected code %d, but got %d", http.StatusInternalServerError, mock.dataWriteHeader)
	}
}
