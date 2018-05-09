package omisego

import (
	"github.com/icrowley/fake"
	"net/http"
	"testing"
)

// Move this to a dotenv file https://github.com/joho/godotenv
var (
	apiKeyId = "api_01cd29gxqbk0a7c859t5v8g4bx"
	apiKey   = "C-b0WYz2L6gvUB-HAwBlcANu0ktoMFTCxJkzKlo3FmU"
	email    = "admin@example.com"
	pwd      = "u22rNF38veC5acIDS1flgA"

	aua = &AdminUserAuth{
		apiKeyId:      apiKeyId,
		apiKey:        apiKey,
		userId:        "usr_01cd29gyb4yrtnf1dmxqm33kbs",
		userAuthToken: "Hz8znbje5ahmVGlnmhG0kh8yIfdr5qr6R-sq9tUgQ08",
	}
	aca = &AdminClientAuth{
		apiKeyId: apiKeyId,
		apiKey:   apiKey,
	}
	sa = &ServerAuth{
		accessKey: "68HazVGtFNAw4rbSa7k2oz3UvSWOG6MCydXyuPmYoqg",
		secretKey: "AVqnuxAlbOtpPIer89BjPCLTMQh_PY8g0wd_Dxd-pGU",
	}
)

func TestLoginLogout(t *testing.T) {
	adminClient := AdminAPI{
		auth:      aca,
		c:         &http.Client{},
		baseUrl:   "http://localhost:4000",
		serverUrl: "/admin/api/",
	}

	body := LoginParams{
		Email:    email,
		Password: pwd,
	}
	res, err := adminClient.Login(body)
	if err != nil {
		t.Fatalf("%s %+v", err, res)
	}

	loggedInAdmin := &AdminUserAuth{
		apiKeyId:      apiKeyId,
		apiKey:        apiKey,
		userId:        res.Data["user_id"].(string),
		userAuthToken: res.Data["authentication_token"].(string),
	}
	adminUser := AdminAPI{
		auth:      loggedInAdmin,
		c:         &http.Client{},
		baseUrl:   "http://localhost:4000",
		serverUrl: "/admin/api/",
	}
	res, err = adminUser.Logout()
	if err != nil {
		t.Fatalf("%s %+v", err, res)
	}
}

func TestCreateUser(t *testing.T) {
	serverUser := EWalletAPI{
		auth:      sa,
		c:         &http.Client{},
		baseUrl:   "http://localhost:4000",
		serverUrl: "/api/",
	}

	body := UserCreateParams{
		ProviderUserId: fake.CharactersN(10),
		Username:       fake.UserName(),
		Metadata: map[string]interface{}{
			"first_name": fake.FirstName(),
			"last_name":  fake.LastName(),
		},
	}

	res, err := serverUser.UserCreate(body)
	if err != nil {
		t.Fatalf("%s %+v", err, res)
	}
}
