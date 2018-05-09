package omisego

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"errors"
	// log "github.com/sirupsen/logrus"
	"net/http"
)

type AdminAPI struct {
	auth      Authorization
	c         *http.Client
	baseUrl   string
	serverUrl string
}

type Response struct {
	Version string                 `json:"version"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

func (a *AdminAPI) addDefaultHeaders(req *http.Request) {
	req.Header.Set("Authorization", a.auth.CreateAuthorizationHeader())
	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	req.Header.Set("accept", "application/vnd.omisego.v1+json")
}

/////////////////
// Session
/////////////////
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AdminAPI) Login(reqBody LoginParams) (Response, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	req, err := http.NewRequest("POST", a.baseUrl+a.serverUrl+"login", b)
	a.addDefaultHeaders(req)
	res, err := a.c.Do(req)
	if err != nil {
		return Response{}, err
	}

	resBody := Response{}
	json.NewDecoder(res.Body).Decode(&resBody)
	if !resBody.Success {
		return resBody, errors.New("Unsuccessful request.")
	}
	return resBody, nil
}

func (a *AdminAPI) Logout() (Response, error) {
	req, err := http.NewRequest("POST", a.baseUrl+a.serverUrl+"logout", nil)
	a.addDefaultHeaders(req)
	res, err := a.c.Do(req)
	if err != nil {
		return Response{}, err
	}

	resBody := Response{}
	json.NewDecoder(res.Body).Decode(&resBody)
	if !resBody.Success {
		return resBody, errors.New("Unsuccessful request.")
	}
	return resBody, nil
}

/////////////////
// API Access
/////////////////
func (a *AdminAPI) AccessKeyCreate() (Response, error) {
	req, err := http.NewRequest("POST", a.baseUrl+a.serverUrl+"access_key.create", nil)
	a.addDefaultHeaders(req)
	res, err := a.c.Do(req)
	if err != nil {
		return Response{}, errors.New("Unsuccesful request")
	}

	resBody := Response{}
	json.NewDecoder(res.Body).Decode(&resBody)
	if !resBody.Success {
		return resBody, errors.New("Unsuccesful request")
	}
	return resBody, nil
}
