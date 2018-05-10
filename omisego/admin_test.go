package omisego

import (
	// log "github.com/sirupsen/logrus"
	// "net/http"
	"net/url"
	// "os"
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

	loginBody = LoginParams{
		Email:    email,
		Password: pwd,
	}
)

func TestLogin(t *testing.T) {
	c, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{
		Client: c,
	}
	body := LoginParams{
		Email:    email,
		Password: pwd,
	}
	_, err := adminClient.Login(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogout(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	err := adminClient.Logout()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccessKeyCreate(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	_, err := adminClient.AccessKeyCreate()
	if err != nil {
		t.Fatal(err)
	}
}
