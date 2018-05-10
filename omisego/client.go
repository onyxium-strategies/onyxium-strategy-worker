package omisego

import (
	"bytes"
	"encoding/json"
	"errors"
	// log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type Response struct {
	Version string                 `json:"version"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

type Client struct {
	auth       Authorization
	BaseURL    *url.URL
	httpClient *http.Client
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

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 500 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
