package omisego_test

import (
	omg "bitbucket.org/visa-startups/coinflow-strategy-worker/omisego"
	"github.com/icrowley/fake"
	"net/http"
	"net/url"
	"testing"
)

var (
	ewalletURL = &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/api",
	}
	sa = &omg.ServerAuth{
		AccessKey: "68HazVGtFNAw4rbSa7k2oz3UvSWOG6MCydXyuPmYoqg",
		SecretKey: "AVqnuxAlbOtpPIer89BjPCLTMQh_PY8g0wd_Dxd-pGU",
	}
	serverUser = omg.EWalletAPI{
		Client: &omg.Client{
			Auth:       sa,
			HttpClient: &http.Client{},
			BaseURL:    ewalletURL,
		},
	}
)

func TestCreateUser(t *testing.T) {
	body := omg.UserCreateParams{
		ProviderUserId: fake.CharactersN(10),
		Username:       fake.UserName(),
		Metadata: map[string]interface{}{
			"first_name": fake.FirstName(),
			"last_name":  fake.LastName(),
		},
	}

	_, err := serverUser.UserCreate(body)
	if err != nil {
		t.Fatal(err)
	}
}
