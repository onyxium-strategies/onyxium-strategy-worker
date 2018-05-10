package omisego

import (
	"github.com/icrowley/fake"
	"net/http"
	"net/url"
	"testing"
)

var (
	sa = &ServerAuth{
		accessKey: "68HazVGtFNAw4rbSa7k2oz3UvSWOG6MCydXyuPmYoqg",
		secretKey: "AVqnuxAlbOtpPIer89BjPCLTMQh_PY8g0wd_Dxd-pGU",
	}
	ewalletURL = &url.URL{
		Scheme: "http",
		Host:   "localhost:4000",
		Path:   "/api",
	}
	serverUser = EWalletAPI{
		Client: Client{
			auth:       sa,
			httpClient: &http.Client{},
			BaseURL:    ewalletURL,
		},
	}
)

func TestCreateUser(t *testing.T) {
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
		t.Fatal(err)
	}
	if !res.Success {
		t.Fatalf("%+v", res)
	}
}
