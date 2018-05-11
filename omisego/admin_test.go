package omisego_test

import (
	omg "bitbucket.org/visa-startups/coinflow-strategy-worker/omisego"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	// "net/http"
	"net/url"
	"testing"
)

// Testing variables
// Use a dotenv file during development https://github.com/joho/godotenv
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
	// aca = &omg.AdminClientAuth{
	// 	ApiKey:   apiKey,
	// 	ApiKeyId: apiKeyId,
	// }
	// adminClient = omg.AdminAPI{
	// 	Client: &omg.Client{
	// 		BaseURL:    adminURL,
	// 		Auth:       aca,
	// 		HttpClient: &http.Client{},
	// 	},
	// }

	loginBody = omg.LoginParams{
		Email:    email,
		Password: pwd,
	}
)

func TestStuff(t *testing.T) {
	input := `{"object":"list","data":[{"object":"minted_token","id":"tok_ABC_01cbfge9qhmsdbjyb7a8e8pxt3","symbol":"ABC","name":"ABC Point","subunit_to_unit":100,"created_at":"2018-01-01T00:00:00Z","updated_at":"2018-01-01T10:00:00Z"}],"pagination":{"per_page":10,"current_page":1,"is_first_page":true,"is_last_page":true}}`

	data, ok := gjson.Parse(input).Value().(map[string]interface{})
	log.Info(data)
	if !ok {
		t.Fatal("Json input is not a slice")
	}
	var i omg.MintedTokenList
	err := mapstructure.Decode(data, &i)
	log.Infof("%#v", i)
	if err != nil {
		t.Fatal(err)
	}
	// log.Info(i.Data[1].(Transaction).Id)
}

func TestLogin(t *testing.T) {
	c, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{
		Client: c,
	}
	body := omg.LoginParams{
		Email:    email,
		Password: pwd,
	}
	_, err := adminClient.Login(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogout(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	adminClient.Login(loginBody)
	err := adminClient.Logout()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccessKeyCreate(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	adminClient.Login(loginBody)
	_, err := adminClient.AccessKeyCreate()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthTokenSwitchAccount(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	adminClient.Login(loginBody)
	body := omg.AuthTokenSwitchAccountParams{
		AccountId: "the_account_id",
	}
	_, err := adminClient.AuthTokenSwitchAccount(body)
	if err.Error() != "{Code:account:not_found Description:There is no user corresponding to the provided account id Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestPasswordReset(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	body := omg.PasswordResetParams{
		Email:       "test@example.com",
		RedirectUrl: "https://example.com/admin/update_password?email={email}&token={token}",
	}
	err := adminClient.PasswordReset(body)
	if err.Error() != "{Code:user:email_not_found Description:There is no user corresponding to the provided email Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestPasswordUpdate(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	body := omg.PasswordUpdateParams{
		Email:                "test@example.com",
		Token:                "26736ca1-43a0-442b-803e-76220cd3cb1d",
		Password:             "nZi9Enc5$l#",
		PasswordConfirmation: "nZi9Enc5$l#",
	}
	err := adminClient.PasswordUpdate(body)
	if err.Error() != "{Code:user:email_not_found Description:There is no user corresponding to the provided email Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestMintedTokenAll(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	adminClient.Login(loginBody)
	body := omg.ListParams{
		Page:    1,
		PerPage: 10,
	}
	_, err := adminClient.MintedTokenAll(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMintedTokenGet(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	adminClient.Login(loginBody)
	body := omg.ByIdParam{
		Id: "tok_ABC_01cbfge9qhmsdbjyb7a8e8pxt3",
	}
	_, err := adminClient.MintedTokenGet(body)
	if err.Error() != "{Code:minted_token:id_not_found Description:There is no minted token corresponding to the provided id Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestMintedTokenCreate(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	adminClient.Login(loginBody)
	body := omg.MintedTokenCreateParams{
		Symbol:        "OMG",
		Name:          "Omisego",
		Description:   "desc",
		SubunitToUnit: 100,
	}
	_, err := adminClient.MintedTokenCreate(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMintedTokenMint(t *testing.T) {
	client, _ := omg.NewClient(apiKeyId, apiKey, adminURL)
	adminClient := omg.AdminAPI{client}
	adminClient.Login(loginBody)
	body := omg.MintedTokenMintParams{
		Id:     "ce3982f5-4a27-498d-a91b-7bb2e2a8d3d1",
		Amount: 1000,
	}
	_, err := adminClient.MintedTokenMint(body)
	if err.Error() != "{Code:minted_token:id_not_found Description:There is no minted token corresponding to the provided id Messages:map[]}" {
		t.Fatal(err)
	}
}
