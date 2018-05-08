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
	c *http.Client
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
	data := []byte(AdminUserAuth.apiKeyId + ":" + AdminUserAuth.apiKey + ":" + AdminUserAuth + ":" + AdminUserAuth.userAuthToken)
	str := base64.StdEncoding.EncodeToString(data)
	return str
}

func ExampleReqeust() {
	client := &http.Client{}

	var body = []byte(`{"page": 1,"per_page": 10,"search_term": "","sort_by": "email","sort_dir": "asc"}`)
	req, err := http.NewRequest("POST", "http://localhost:4000/admin/api/admin.all", bytes.NewBuffer(body))

	req.Header.Set("Authorization", "OMGAdmin MWNlOTNkYWItYmFmMi00MzEzLWIzNDgtZjI5NjJiZGY5MDFkOlBtWE81N0p3N09pWVJxYk1Ha3pXQk4yZC1yaDYwbExWcDZJTzBJLUYyZ286NTk5YmE3ZTAtMmU2MC00ODk4LThlYTEtNzhiYjJmZTNlYWZiOmN5a2RDVDFWSXFpci1ybHlmeDR4cEw0VjJsazVIT1FXRm5XQzdURFNIY1k=")
	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	req.Header.Set("accept", "application/vnd.omisego.v1+json")

	log.Info("Request: ", req)
	res, err := client.Do(req)
	log.Info("Res: ", res)
	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Response: %s", string(resBody))
}
