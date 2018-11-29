package cherwell

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

// Client contains information about Cherwell API response after authorization request
type Client struct {
	tokenResponse *tokenResponse
	httpClient    *http.Client
	config        Config
	mutex         sync.RWMutex
}

// NewClient creates an instance of Client and obtains access token
func NewClient(conf Config, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		return nil, errors.New("cherwell.NewClient: httpClient can not be nil")
	}

	client := &Client{config: conf, httpClient: httpClient}

	if err := client.getAccessToken(); err != nil {
		return nil, fmt.Errorf("cherwell.NewClient authentication error: %s", err)
	}
	return client, nil
}

// getAccessToken is a func for authenticating using the Internal Mode. In this scenario, the User logs in to the
// REST API using CSM credentials. CSM returns a JSON response that includes information about the access token or error
// with status code.
func (c *Client) getAccessToken() error {
	var errRes errorResponse

	vals := url.Values{
		"grant_type": {passwordGrantType},
		"client_id":  {c.config.ClientID},
		"username":   {c.config.UserName},
		"password":   {c.config.Password},
		"auth_mode":  {c.config.AuthMode},
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	resp, err := c.httpClient.PostForm(c.config.Host+tokenEndpoint, vals)
	if err != nil {
		return fmt.Errorf("getAccessToken authentication request failed: %s", err)
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode == http.StatusOK {
		return unmarshalRespBody(resp, &c.tokenResponse)
	}

	if err = unmarshalRespBody(resp, &errRes); err != nil {
		return fmt.Errorf("getAccessToken failed with error: %s , during deserialization response body: %s", err, resp.Body)
	}
	return errRes
}

// unmarshalRespBody perform unmarshals Cherwell API response body into given structure
func unmarshalRespBody(resp *http.Response, out interface{}) (err error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &out)
	if err != nil {
		return fmt.Errorf("non-JSON response received, %s: %s", err.Error(), string(data))
	}

	return nil
}

func (c *Client) performRequest(method, path string, reqEntity interface{}, respEntity interface{}) error {
	resp, err := c.getResponse(method, path, reqEntity)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint: errcheck
	if resp.StatusCode == http.StatusUnauthorized {
		if err := c.refreshAccessToken(); err != nil {
			return err
		}
		resp, err = c.getResponse(method, path, reqEntity)
		if err != nil {
			return err
		}
	}
	if err := unmarshalRespBody(resp, respEntity); err != nil {
		return &GeneralFailure{Message: err.Error()}
	}
	return nil
}

func (c *Client) refreshAccessToken() error {
	retry := 0
	for retry < retryCount {
		err := c.getAccessToken()
		if err != nil {
			retry++
			continue
		} else {
			return nil
		}
	}
	return &CherwellError{Code: "500", Message: "cherwell server does not respond"}
}

func (c *Client) createRequest(method, path string, data io.Reader) (*http.Request, error) {
	requestPath := c.config.Host + path
	req, err := http.NewRequest(method, requestPath, data)
	if err != nil {
		return nil, fmt.Errorf("getResponse failed to create request: %s", err)
	}

	c.mutex.RLock()
	req.Header.Set("Authorization", "Bearer "+c.tokenResponse.AccessToken)
	c.mutex.RUnlock()
	return req, nil
}

func (c *Client) getResponse(method, path string, reqEntity interface{}) (resp *http.Response, err error) {
	var reqBody io.Reader

	if reqEntity != nil {
		data, err := json.Marshal(&reqEntity)

		if err != nil {
			return nil, fmt.Errorf("performRequest failed to marshal request: %s", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := c.createRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, fmt.Errorf("getResponse failed to create request: %s", err)
	}

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("getResponse failed to send request: %s", err)
	}

	return resp, err
}
