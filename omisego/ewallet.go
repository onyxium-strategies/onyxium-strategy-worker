package omisego

import (
	"bytes"
	"crypto/rand"
	"fmt"
	log "github.com/sirupsen/logrus"
	// "io/ioutil"
	"encoding/json"
	"errors"
	"net/http"
)

type EWalletAPI struct {
	auth      Authorization
	c         *http.Client
	baseUrl   string
	serverUrl string
}

func (e *EWalletAPI) addDefaultHeaders(req *http.Request) {
	req.Header.Set("Authorization", e.auth.CreateAuthorizationHeader())
	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
	req.Header.Set("accept", "application/vnd.omisego.v1+json")
}

/////////////////
// User
/////////////////
type UserCreateParams struct {
	ProviderUserId    string                 `json:"provider_user_id"`
	Username          string                 `json:"username"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	EncryptedMetadata map[string]interface{} `json:"encrypted_metadata,omitempty"`
}

func (e *EWalletAPI) UserCreate(reqBody UserCreateParams) (Response, error) {
	log.Infof("%+v", reqBody)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reqBody)
	log.Infof("%+v", reqBody)
	req, err := http.NewRequest("POST", e.baseUrl+e.serverUrl+"user.create", b)
	e.addDefaultHeaders(req)
	req.Header.Set("Idempotency-Token", newIdempotencyToken())
	res, err := e.c.Do(req)
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

// func ExampleAdminUserReqeust() {
// 	client := EWalletAPI{
// 		auth: aca,
// 		c:    &http.Client{},
// 	}

// 	var body = []byte(`{"page": 1,"per_page": 10,"search_term": "","sort_by": "email","sort_dir": "asc"}`)
// 	req, err := http.NewRequest("POST", "http://localhost:4000/admin/api/admin.all", bytes.NewBuffer(body))

// 	req.Header.Set("Authorization", client.auth.CreateAuthorizationHeader())
// 	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
// 	req.Header.Set("accept", "application/vnd.omisego.v1+json")

// 	log.Infof("Request: %+v", req)
// 	res, err := client.c.Do(req)
// 	log.Infof("Res: %+v", res)
// 	resBody, err := ioutil.ReadAll(res.Body)
// 	res.Body.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Infof("Response: %s", string(resBody))
// }

// func ExampleServerReqeust() {
// 	client := EWalletAPI{
// 		auth: sa,
// 		c:    &http.Client{},
// 	}

// 	var body = []byte(`{"provider_user_id":"123456","username":"alain@omise.co","metadata":{"first_name":"Alain","last_name":"Koning"},"encrypted_metadata":{}}`)
// 	req, err := http.NewRequest("POST", "http://localhost:4000/api/user.create", bytes.NewBuffer(body))

// 	req.Header.Set("Authorization", client.auth.CreateAuthorizationHeader())
// 	req.Header.Set("Content-Type", "application/vnd.omisego.v1+json")
// 	req.Header.Set("accept", "application/vnd.omisego.v1+json")
// 	req.Header.Set("Idempotency-Token", newIdempotencyToken())

// 	log.Infof("Request: %+v", req)
// 	res, err := client.c.Do(req)
// 	log.Infof("Res: %+v", res)
// 	resBody, err := ioutil.ReadAll(res.Body)
// 	res.Body.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Infof("Response: %s", string(resBody))
// }

// newIdempotencyToken returns a randomly generated idempotency token.
func newIdempotencyToken() string {
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
