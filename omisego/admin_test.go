package omisego

import (
	"net/url"
	// "os"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
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

func TestStuff(t *testing.T) {
	input := `{"user_id":"usr_01cd29gyb4yrtnf1dmxqm33kbs","user":{"username":null,"updated_at":"2018-05-09T11:10:33.316287Z","socket_topic":"user:usr_01cd29gyb4yrtnf1dmxqm33kbs","provider_user_id":null,"object":"user","metadata":{},"id":"usr_01cd29gyb4yrtnf1dmxqm33kbs","encrypted_metadata":{},"email":"admin@example.com","created_at":"2018-05-09T11:10:33.316267Z","avatar":{"thumb":null,"small":null,"original":null,"large":null}},"object":"authentication_token","master_admin":true,"authentication_token":"J_Nh_NxHuZd_Hy6P5Y-bSdAqmV7c_Uk9GvdBsI2CYaQ","account_id":"acc_01cd29gxgyymjbvf03bygcb076","account":{"updated_at":"2018-05-09T11:10:32.484352Z","socket_topic":"account:acc_01cd29gxgyymjbvf03bygcb076","parent_id":null,"object":"account","name":"master_account","metadata":{},"master":true,"id":"acc_01cd29gxgyymjbvf03bygcb076","encrypted_metadata":{},"description":"Master Account","created_at":"2018-05-09T11:10:32.484335Z","avatar":{"thumb":null,"small":null,"original":null,"large":null}}}`

	data, ok := gjson.Parse(input).Value().(map[string]interface{})
	log.Info(data)
	if !ok {
		t.Fatal("Json input is not a slice")
	}
	var i AuthenicationToken
	err := mapstructure.Decode(data, &i)
	log.Infof("%#v", i)
	if err != nil {
		t.Fatal(err)
	}
	// log.Info(i.Data[1].(Transaction).Id)
}

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

func TestAuthTokenSwitchAccount(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	body := AuthTokenSwitchAccountParams{
		AccountId: "the_account_id",
	}
	_, err := adminClient.AuthTokenSwitchAccount(body)
	if err.Error() != "{Object:error Code:account:not_found Description:There is no user corresponding to the provided account id Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestPasswordReset(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	body := PasswordResetParams{
		Email:       "test@example.com",
		RedirectUrl: "https://example.com/admin/update_password?email={email}&token={token}",
	}
	err := adminClient.PasswordReset(body)
	if err.Error() != "{Object:error Code:user:email_not_found Description:There is no user corresponding to the provided email Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestPasswordUpdate(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	body := PasswordUpdateParams{
		Email:                "test@example.com",
		Token:                "26736ca1-43a0-442b-803e-76220cd3cb1d",
		Password:             "nZi9Enc5$l#",
		PasswordConfirmation: "nZi9Enc5$l#",
	}
	err := adminClient.PasswordUpdate(body)
	if err.Error() != "{Object:error Code:user:email_not_found Description:There is no user corresponding to the provided email Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestMintedTokenAll(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	body := MintedTokenAllParams{
		Page:    1,
		PerPage: 10,
	}
	_, err := adminClient.MintedTokenAll(body)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMintedTokenGet(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	body := MintedTokenGetParams{
		Id: "tok_ABC_01cbfge9qhmsdbjyb7a8e8pxt3",
	}
	_, err := adminClient.MintedTokenGet(body)
	if err.Error() != "{Object:error Code:minted_token:id_not_found Description:There is no minted token corresponding to the provided id Messages:map[]}" {
		t.Fatal(err)
	}
}

func TestMintedTokenCreate(t *testing.T) {
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	body := MintedTokenCreateParams{
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
	client, _ := NewClient(apiKeyId, apiKey, adminURL)
	adminClient := AdminAPI{client}
	adminClient.Login(loginBody)
	body := MintedTokenMintParams{
		Id:     "ce3982f5-4a27-498d-a91b-7bb2e2a8d3d1",
		Amount: 1000,
	}
	_, err := adminClient.MintedTokenMint(body)
	if err.Error() != "{Object:error Code:minted_token:id_not_found Description:There is no minted token corresponding to the provided id Messages:map[]}" {
		t.Fatal(err)
	}
}
