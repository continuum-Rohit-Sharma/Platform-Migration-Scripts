package rest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pkg/errors"
)

type mockResponseWriter struct {
	dataWriteHeader int
	dataWrite       []byte
	dataHeader      http.Header
}

func (m *mockResponseWriter) Header() http.Header {
	return m.dataHeader
}

func (m *mockResponseWriter) Write(b []byte) (int, error) {
	m.dataWrite = b
	return len(m.dataWrite), nil
}

func (m *mockResponseWriter) WriteHeader(code int) {
	m.dataWriteHeader = code
}

type TestCase struct {
	name string
	http.ResponseWriter
	code          int
	response      []byte
	message       string
	err           error
	expectMessage string
}

func mockCreateErrorGenerator() func(msg string) interface{} {
	var flag bool
	return func(msg string) interface{} {
		if !flag {
			flag = true
			return make(chan int)
		}
		return ""
	}
}

var cases = []TestCase{
	{name: "TestStatusInternalServerError", ResponseWriter: http.ResponseWriter(&mockResponseWriter{dataHeader: http.Header{}}), code: http.StatusInternalServerError, response: []byte{}, message: "", err: errors.New("error"), expectMessage: `{"error":{"message":"Internal Server Error"}}`},
	{name: "TestStatusNotFound", ResponseWriter: http.ResponseWriter(&mockResponseWriter{dataHeader: http.Header{}}), code: http.StatusNotFound, response: []byte("test"), message: "message1", err: nil, expectMessage: `{"error":{"message":"message1"}}`},
	{name: "TestStatusOK", ResponseWriter: http.ResponseWriter(&mockResponseWriter{dataHeader: http.Header{}}), code: http.StatusOK, response: []byte("test bla-bla"), message: "message2", err: nil, expectMessage: `{"error":{"message":"message2"}}`},
}

func TestRender(t *testing.T) {
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			render(http.ResponseWriter(test.ResponseWriter), test.code, test.response)
			mock := test.ResponseWriter.(*mockResponseWriter)
			if mock.dataWriteHeader != test.code {
				t.Errorf("expected code %d, but got %d", test.code, mock.dataWriteHeader)
			}
			if len(mock.dataWrite) != len(test.response) {
				t.Errorf("expected response len %d, but got %d", len(test.response), len(mock.dataWrite))
			}
			if mock.dataHeader.Get(contentType) != applicationJSON {
				t.Errorf("expected %q = %q, but got %q", contentType, applicationJSON, mock.dataHeader.Get(contentType))
			}
		})
	}
}

func TestSendError(t *testing.T) {
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			SendError(http.ResponseWriter(test.ResponseWriter), test.code, test.message, test.err)
			mock := test.ResponseWriter.(*mockResponseWriter)
			if mock.dataWriteHeader != test.code {
				t.Errorf("expected code %d, but got %d", test.code, mock.dataWriteHeader)
			}
			if string(mock.dataWrite) != test.expectMessage {
				t.Errorf("expected msg %s, but got %s", test.expectMessage, string(mock.dataWrite))
			}
		})
	}
}

func TestSendError_InternalServerError(t *testing.T) {
	oldCreateError := createError
	createError = mockCreateErrorGenerator()
	defer func() {
		createError = oldCreateError
	}()
	mock := &mockResponseWriter{dataHeader: http.Header{}}
	SendError(http.ResponseWriter(mock), http.StatusOK, "message", nil)

	if mock.dataWriteHeader != http.StatusOK {
		t.Errorf("expected code %d, but got %d", http.StatusOK, mock.dataWriteHeader)
	}
}

func TestSendInternalServerError(t *testing.T) {
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			SendInternalServerError(http.ResponseWriter(test.ResponseWriter), test.message, test.err)
			mock := test.ResponseWriter.(*mockResponseWriter)
			if mock.dataWriteHeader != http.StatusInternalServerError {
				t.Errorf("expected code %d, but got %d", http.StatusInternalServerError, mock.dataWriteHeader)
			}
			if string(mock.dataWrite) != test.expectMessage {
				t.Errorf("expected msg %s, but got %s", test.expectMessage, string(mock.dataWrite))
			}
		})
	}
}

func TestRenderJSON(t *testing.T) {
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			RenderJSON(http.ResponseWriter(test.ResponseWriter), test.message)
			mock := test.ResponseWriter.(*mockResponseWriter)
			if mock.dataWriteHeader != http.StatusOK {
				t.Errorf("expected code %d, but got %d", http.StatusOK, mock.dataWriteHeader)
			}
			if string(mock.dataWrite) != fmt.Sprintf("%q", test.message) {
				t.Errorf("expected msg %s, but got %s", fmt.Sprintf("%q", test.message), string(mock.dataWrite))
			}
		})
	}
}

func Test_RenderJSON_Fail_Marshal(t *testing.T) {
	mock := &mockResponseWriter{dataHeader: http.Header{}}
	renderJSON(http.ResponseWriter(mock), http.StatusOK, make(chan int))
	if mock.dataWriteHeader != http.StatusInternalServerError {
		t.Errorf("expected code %d, but got %d", http.StatusInternalServerError, mock.dataWriteHeader)
	}
}
