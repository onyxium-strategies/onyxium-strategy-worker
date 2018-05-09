package main

import (
	"bytes"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Authorization interface {
	CreateAuthorizationHeader() string
}

type OmisegoAPI struct {
	auth Authorization
	c    *http.Client
}

type AdminClientAuth struct {
	apiKeyId string
	apiKey   string
}

type AdminUserAuth struct {
	apiKeyId      string
	apiKey        string
	userId        string
	userAuthToken string
}

func (a *AdminClientAuth) CreateAuthorizationHeader() string {
	data := []byte(a.apiKeyId + ":" + a.apiKey)
	str := base64.StdEncoding.EncodeToString(data)
	return "OMGAdmin " + str
}

func (a *AdminUserAuth) CreateAuthorizationHeader() string {
	data := []byte(a.apiKeyId + ":" + a.apiKey + ":" + a.userId + ":" + a.userAuthToken)
	str := base64.StdEncoding.EncodeToString(data)
	return "OMGAdmin " + str
}

func ExampleReqeust() {
	client := OmisegoAPI{
		auth: &AdminUserAuth{
			apiKeyId:      "e8a74b60-4959-40a6-92d7-dbb6986f80c2",
			apiKey:        "j8cPwzVTYG1l8i-4KNG4NvTHBiQU8TTBEGCIerrDkY8",
			userId:        "e4c6087c-034e-40cf-b23f-820a865689a7",
			userAuthToken: "hdTAcBwCJkp1Py8qZacf294cwAhQiQmSaXj0SbCGpfw",
		},
		c: &http.Client{},
	}

	var body = []byte(`{"page": 1,"per_page": 10,"search_term": "","sort_by": "email","sort_dir": "asc"}`)
	req, err := http.NewRequest("POST", "http://localhost:4000/admin/api/admin.all", bytes.NewBuffer(body))

	req.Header.Set("Authorization", client.auth.CreateAuthorizationHeader())
	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	req.Header.Set("accept", "application/vnd.omisego.v1+json")

	log.Info("Request: ", req)
	res, err := client.c.Do(req)
	log.Info("Res: ", res)
	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Response: %s", string(resBody))
}

func main() {
	ExampleReqeust()
}
