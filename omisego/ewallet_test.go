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

func TestUserCreate(t *testing.T) {
	body := omg.UserParams{
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

func TestUserGet(t *testing.T) {
	body := omg.ProviderUserIdParam{
		PrividerUserId: "7x1hsxeryf",
	}

	_, err := serverUser.UserGet(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserListBalances(t *testing.T) {
	body := omg.ProviderUserIdParam{
		PrividerUserId: "7x1hsxeryf",
	}

	_, err := serverUser.UserListBalances(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserCreditBalance(t *testing.T) {
	body := omg.BalanceAdjustmentParams{
		PrividerUserId: "7x1hsxeryf",
		TokenId:        "BTC",
		Amount:         100,
	}

	_, err := serverUser.UserCreditBalance(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDebitBalance(t *testing.T) {
	body := omg.BalanceAdjustmentParams{
		PrividerUserId: "7x1hsxeryf",
		TokenId:        "BTC",
		Amount:         100,
	}

	_, err := serverUser.UserDebitBalance(body)
	if err != nil {
		t.Fatal(err)
	}
}
