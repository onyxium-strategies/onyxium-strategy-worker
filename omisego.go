package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
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

type ServerAuth struct {
	accessKey string
	secretKey string
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

func (s *ServerAuth) CreateAuthorizationHeader() string {
	data := []byte(s.accessKey + ":" + s.secretKey)
	str := base64.StdEncoding.EncodeToString(data)
	return "OMGServer " + str
}

var (
	aua = &AdminUserAuth{
		apiKeyId:      "e8a74b60-4959-40a6-92d7-dbb6986f80c2",
		apiKey:        "j8cPwzVTYG1l8i-4KNG4NvTHBiQU8TTBEGCIerrDkY8",
		userId:        "e4c6087c-034e-40cf-b23f-820a865689a7",
		userAuthToken: "hdTAcBwCJkp1Py8qZacf294cwAhQiQmSaXj0SbCGpfw",
	}
	aca = &AdminClientAuth{
		apiKeyId: "e8a74b60-4959-40a6-92d7-dbb6986f80c2",
		apiKey:   "j8cPwzVTYG1l8i-4KNG4NvTHBiQU8TTBEGCIerrDkY8",
	}
	sa = &ServerAuth{
		accessKey: "94Fk-3qQwWmRZ8id8c5Q3SEHE7lpmdeQXuOi5TV4Q4c",
		secretKey: "NDccDpezX7W2HJs7V4MLDprbhc9DW6G3kj_Xqg2UJmA",
	}
)

func ExampleAdminUserReqeust() {
	client := OmisegoAPI{
		auth: aua,
		c:    &http.Client{},
	}

	var body = []byte(`{"page": 1,"per_page": 10,"search_term": "","sort_by": "email","sort_dir": "asc"}`)
	req, err := http.NewRequest("POST", "http://localhost:4000/admin/api/admin.all", bytes.NewBuffer(body))

	req.Header.Set("Authorization", client.auth.CreateAuthorizationHeader())
	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	req.Header.Set("accept", "application/vnd.omisego.v1+json")

	log.Infof("Request: %+v", req)
	res, err := client.c.Do(req)
	log.Infof("Res: %+v", res)
	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Response: %s", string(resBody))
}

func ExampleServerReqeust() {
	client := OmisegoAPI{
		auth: sa,
		c:    &http.Client{},
	}

	var body = []byte(`{"provider_user_id":"123456","username":"alain@omise.co","metadata":{"first_name":"Alain","last_name":"Koning"},"encrypted_metadata":{}}`)
	req, err := http.NewRequest("POST", "http://localhost:4000/api/user.create", bytes.NewBuffer(body))

	req.Header.Set("Authorization", client.auth.CreateAuthorizationHeader())
	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	req.Header.Set("accept", "application/vnd.omisego.v1+json")
	req.Header.Set("Idempotency-Token", getIdempotencyToken())

	log.Infof("Request: %+v", req)
	res, err := client.c.Do(req)
	log.Infof("Res: %+v", res)
	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Response: %s", string(resBody))
}

// getIdempotencyToken returns a randomly generated idempotency token.
func getIdempotencyToken() string {
	b := make([]byte, 16)
	rand.Reader.Read(b)
	return uUIDVersion4(b)
}

// UUIDVersion4 returns a Version 4 random UUID from the byte slice provided
func uUIDVersion4(u []byte) string {
	// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
	// 13th character is "4"
	u[6] = (u[6] | 0x40) & 0x4F
	// 17th character is "8", "9", "a", or "b"
	u[8] = (u[8] | 0x80) & 0xBF

	return fmt.Sprintf(`%X-%X-%X-%X-%X`, u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func main() {
	// ExampleAdminUserReqeust()
	ExampleAdminUserReqeust2()
	// ExampleServerReqeust()
}
