package omisego

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	// log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Response struct {
	Version string                 `json:"version"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

type ErrorResponse struct {
	Object      string            `json:"object"`
	Code        string            `json:"code"`
	Description string            `json:"description"`
	Messages    map[string]string `json:"messages"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%+v", *e)
}

type Client struct {
	auth       Authorization
	BaseURL    *url.URL
	httpClient *http.Client
	Log        io.Writer // If user set log file name all requests will be logged there
}

// NewClient returns new Client struct
func NewClient(apiKeyId string, apiKey string, BaseURL *url.URL) (*Client, error) {
	if apiKeyId == "" || apiKey == "" || BaseURL == nil {
		return nil, errors.New("ApiKeyID, ApiKey and BaseURL are required to create a Client")
	}

	return &Client{
		httpClient: &http.Client{},
		auth: &AdminClientAuth{
			apiKey:   apiKey,
			apiKeyId: apiKeyId,
		},
		BaseURL: BaseURL,
		Log:     nil,
	}, nil
}

func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: c.BaseURL.Path + path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	}

	req.Header.Set("Authorization", c.auth.CreateAuthorizationHeader())
	req.Header.Set("accept", "application/vnd.omisego.v1+json")

	return req, nil
}

func (c *Client) do(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	c.log(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 500 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	r := &Response{}
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	if !r.Success {
		errorResp := &ErrorResponse{}
		err = mapstructure.Decode(r.Data, &errorResp)
		if err != nil {
			return nil, err
		}
		return nil, errorResp
	}
	return r, err
}

// SetHTTPClient sets *http.Client to current client
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// SetLog will set/change the output destination.
// If log file is set client will log all requests and responses to this Writer
func (c *Client) SetLog(log io.Writer) {
	c.Log = log
}

// log will dump request and response to the log file
func (c *Client) log(r *http.Request, resp *http.Response) {
	if c.Log != nil {
		var (
			reqDump  string
			respDump []byte
		)

		if r != nil {
			reqDump = fmt.Sprintf("%s %s. Data: %s", r.Method, r.URL.String(), r.Form.Encode())
		}
		if resp != nil {
			respDump, _ = httputil.DumpResponse(resp, true)
		}

		c.Log.Write([]byte(fmt.Sprintf("Request: %s\nResponse: %s\n", reqDump, string(respDump))))
	}
}
