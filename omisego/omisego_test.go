package omisego

import (
	"net/http"
	"testing"
)

// Move this to a dotenv file https://github.com/joho/godotenv
var (
	apiKeyId = "api_01cd29gxqbk0a7c859t5v8g4bx"
	apiKey   = "C-b0WYz2L6gvUB-HAwBlcANu0ktoMFTCxJkzKlo3FmU"
	email    = "admin@example.com"
	pwd      = "u22rNF38veC5acIDS1flgA"

	// aua = &AdminUserAuth{
	// 	apiKeyId: apiKeyId,
	// 	apiKey:   apiKey,
	// 	userId:        "e4c6087c-034e-40cf-b23f-820a865689a7",
	// 	userAuthToken: "hdTAcBwCJkp1Py8qZacf294cwAhQiQmSaXj0SbCGpfw",
	// }
	aca = &AdminClientAuth{
		apiKeyId: apiKeyId,
		apiKey:   apiKey,
	}
	// adminUser = AdminAPI{
	// 	auth:      aua,
	// 	c:         &http.Client{},
	// 	baseUrl:   "http://localhost:4000",
	// 	serverUrl: "/admin/api/",
	// }
	adminClient = AdminAPI{
		auth:      aca,
		c:         &http.Client{},
		baseUrl:   "http://localhost:4000",
		serverUrl: "/admin/api/",
	}
)

func TestStuff(t *testing.T) {
	body := LoginArgs{
		Email:    email,
		Password: pwd,
	}
	res := adminClient.Login(body)
	t.Log(res.Data["user_id"])
	t.Log(res.Data["authentication_token"])

	aua := &AdminUserAuth{
		apiKeyId:      apiKeyId,
		apiKey:        apiKey,
		userId:        res.Data["user_id"].(string),
		userAuthToken: res.Data["authentication_token"].(string),
	}
	adminUser := AdminAPI{
		auth:      aua,
		c:         &http.Client{},
		baseUrl:   "http://localhost:4000",
		serverUrl: "/admin/api/",
	}
	t.Logf("%+v", adminUser)
	res = adminUser.Logout()
	t.Logf("%+v", res)

	// ExampleAdminUserReqeust()
	// ExampleServerReqeust()
}
