package omisego

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"testing"
)

// Move this to a dotenv file https://github.com/joho/godotenv
var (
	apiKeyId = "api_01cd29gxqbk0a7c859t5v8g4bx"
	apiKey   = "C-b0WYz2L6gvUB-HAwBlcANu0ktoMFTCxJkzKlo3FmU"
	email    = "admin@example.com"
	pwd      = "u22rNF38veC5acIDS1flgA"
	userId   = "usr_01cd29gyb4yrtnf1dmxqm33kbs"

	adminURL = &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/admin/api",
	}

	aua = &AdminUserAuth{
		apiKeyId:      apiKeyId,
		apiKey:        apiKey,
		userId:        userId,
		userAuthToken: "",
	}
	aca = &AdminClientAuth{
		apiKeyId: apiKeyId,
		apiKey:   apiKey,
	}

	adminClient = AdminAPI{
		Client: Client{
			auth:       aca,
			httpClient: &http.Client{},
			BaseURL:    adminURL,
		},
	}
	adminUser = AdminAPI{
		Client: Client{
			auth:       aua,
			httpClient: &http.Client{},
			BaseURL:    adminURL,
		},
	}
)

// We have to login to get the auth token before we can start testing
func TestMain(m *testing.M) {
	body := LoginParams{
		Email:    email,
		Password: pwd,
	}
	res, err := adminClient.Login(body)
	if err != nil {
		log.Fatal(err)
	}
	if !res.Success {
		log.Fatalf("%+v", res)
	}
	aua.userAuthToken = res.Data["authentication_token"].(string)
	os.Exit(m.Run())
}

// func TestLogin(t *testing.T) {
// 	body := LoginParams{
// 		Email:    email,
// 		Password: pwd,
// 	}
// 	res, err := adminClient.Login(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !res.Success {
// 		t.Fatalf("%+v", res)
// 	}
// 	t.Log(res)
// }

func TestLogout(t *testing.T) {
	res, err := adminUser.Logout()
	if err != nil {
		t.Fatal(err)
	}
	if !res.Success {
		t.Fatalf("%+v", res)
	}
}
