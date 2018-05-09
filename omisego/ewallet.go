package omisego

import (
	"bytes"
	"crypto/rand"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type EWalletAPI struct {
	auth      Authorization
	c         *http.Client
	baseUrl   string
	serverUrl string
}

var (
	sa = &ServerAuth{
		accessKey: "94Fk-3qQwWmRZ8id8c5Q3SEHE7lpmdeQXuOi5TV4Q4c",
		secretKey: "NDccDpezX7W2HJs7V4MLDprbhc9DW6G3kj_Xqg2UJmA",
	}
)

func ExampleAdminUserReqeust() {
	client := EWalletAPI{
		auth: aca,
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
	client := EWalletAPI{
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

// uUIDVersion4 returns a Version 4 random UUID from the byte slice provided
func uUIDVersion4(u []byte) string {
	// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
	// 13th character is "4"
	u[6] = (u[6] | 0x40) & 0x4F
	// 17th character is "8", "9", "a", or "b"
	u[8] = (u[8] | 0x80) & 0xBF

	return fmt.Sprintf(`%X-%X-%X-%X-%X`, u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
