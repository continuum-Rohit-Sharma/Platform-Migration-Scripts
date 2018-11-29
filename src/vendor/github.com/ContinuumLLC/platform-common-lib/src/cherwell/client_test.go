package cherwell

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newMockHandler(method, path, resp string, statusCode int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(resp))
	})
}

func newTestServer() (*httptest.Server, *http.ServeMux) {
	defaultTokenResponse := []byte(`{
		"access_token": "access_token",
		"token_type": "bearer",
		"expires_in": 14399,
		"refresh_token": "refresh_token",
		"as:client_id": "client_id",
		"username": "username",
		".issued": "Tue, 31 Jul 2018 14:46:46 GMT",
		".expires": "Tue, 31 Jul 2018 18:46:46 GMT"
	  }`)

	mux := http.NewServeMux()
	mux.HandleFunc(tokenEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write(defaultTokenResponse)
	})
	server := httptest.NewServer(mux)
	return server, mux
}

func newFailedTestServer() (*httptest.Server, *http.ServeMux) {
	defaultTokenResponse := []byte(`{
		"error": "access_token",		
		"error_description": "Tue, 31 Jul 2018 18:46:46 GMT"
	  }`)

	mux := http.NewServeMux()
	mux.HandleFunc(tokenEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Write(defaultTokenResponse)
	})
	server := httptest.NewServer(mux)
	return server, mux
}

func TestNewClientSuccess(t *testing.T) {
	server, _ := newTestServer()
	conf := Config{
		Host: server.URL,
	}
	_, err := NewClient(conf, &http.Client{Transport: &http.Transport{}})
	assert.NoError(t, err)
}

func TestNewClientFailNilHTTPClient(t *testing.T) {
	conf := Config{
		Host: "empty/URL",
	}
	_, err := NewClient(conf, nil)
	assert.EqualError(t, err, "cherwell.NewClient: httpClient can not be nil")
}

func TestNewClientFailObtainAccesToken(t *testing.T) {
	conf := Config{
		Host: "http://Invalid hosthname",
	}
	_, err := NewClient(conf, &http.Client{Transport: &http.Transport{}})
	assert.EqualError(t, err, "cherwell.NewClient authentication error: getAccessToken authentication request failed: parse http://Invalid hosthname/token: invalid character \" \" in host name")
}

func TestObtainAccesTokenSuccess(t *testing.T) {
	server, _ := newTestServer()
	conf := Config{
		Host: server.URL,
	}
	client := &Client{config: conf, httpClient: &http.Client{Transport: &http.Transport{}}}
	assert.NoError(t, client.getAccessToken())
}

func TestObtainAccesTokenFailOnRequest(t *testing.T) {
	conf := Config{
		Host: "http://Invalid hosthname",
	}
	client := &Client{config: conf, httpClient: &http.Client{Transport: &http.Transport{}}}
	err := client.getAccessToken()
	assert.EqualError(t, err, "getAccessToken authentication request failed: parse http://Invalid hosthname/token: invalid character \" \" in host name")
}

func TestObtainAccesTokenFailOnBadResponseCode(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(errorResponse{Err: "Invalid request", ErrorDescription: "password is empty"})
		assert.NoError(t, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))

	conf := Config{
		Host: server.URL,
	}
	client := &Client{config: conf, httpClient: &http.Client{Transport: &http.Transport{}}}
	err := client.getAccessToken()
	assert.EqualError(t, err, "Error: Invalid request\nDescription: password is empty\n")
}

func TestClientRefreshAccessToken(t *testing.T) {
	server, _ := newTestServer()
	failedServer, _ := newFailedTestServer()
	conf := Config{
		Host: server.URL,
	}
	failedConf := Config{
		Host: failedServer.URL,
	}
	client := &Client{config: conf, httpClient: &http.Client{Transport: &http.Transport{}}}

	tests := []struct {
		name          string
		tokenResponse *tokenResponse
		httpClient    *http.Client
		config        Config
		wantErr       bool
	}{
		{name: "Success token refresh",
			tokenResponse: nil,
			httpClient:    client.httpClient,
			config:        conf,
			wantErr:       false,
		},
		{name: "Failed token refresh",
			tokenResponse: nil,
			httpClient:    client.httpClient,
			config:        failedConf,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				tokenResponse: tt.tokenResponse,
				httpClient:    tt.httpClient,
				config:        tt.config,
			}
			if err := c.refreshAccessToken(); (err != nil) != tt.wantErr {
				t.Errorf("Client.refreshAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRaceConditionOnRefreshToken(t *testing.T) {
	server, mux := newTestServer()
	conf := Config{
		Host: server.URL,
	}
	client := &Client{config: conf, httpClient: &http.Client{Transport: &http.Transport{}}}
	client.getAccessToken()

	resp := `{
          "busObPublicId": "pub_id_1",
          "busObRecId": "rec_id_1",
          "cacheKey": "string",
          "fieldValidationErrors": [],
          "notificationTriggers": [],
          "errorCode": "",
          "errorMessage": "",
          "hasError": false
    }`

	test := struct {
		name   string
		method string
		path   string
	}{
		name:   "Success token refresh",
		method: http.MethodGet,
		path:   fmt.Sprintf(getBOByRecIDEndpoint, "12", "123"),
	}

	mux.Handle(test.path, newMockHandler(test.method, test.path, resp, 401))

	t.Run(test.name, func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, method string, path string) {
				for i := 0; i < 500; i++ {
					respEntity := new(businessObjectResponse)
					client.performRequest(test.method, test.path, nil, respEntity)
				}
				wg.Done()
			}(&wg, test.method, test.path)
		}
		wg.Wait()
	})
}
