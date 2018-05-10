package omisego

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	// log "github.com/sirupsen/logrus"
)

type AdminAPI struct {
	*Client
}

/////////////////
// Session
/////////////////
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserId              string                 `mapstructure:"user_id"`
	User                map[string]interface{} `mapstructure:"user"`
	Object              string                 `mapstructure:"object"`
	MasterAdmin         bool                   `mapstructure:"master_admin"`
	AuthenticationToken string                 `mapstructure:"authentication_token"`
	AccountId           string                 `mapstructure:"account_id"`
	Account             map[string]interface{} `mapstructure:"account"`
}

func (a *AdminAPI) Login(reqBody LoginParams) (*LoginResponse, error) {
	req, err := a.newRequest("POST", "/login", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data LoginResponse
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	// Log the user in with new authentication
	a.auth = &AdminUserAuth{
		apiKey:        a.auth.(*AdminClientAuth).apiKey,
		apiKeyId:      a.auth.(*AdminClientAuth).apiKeyId,
		userId:        data.UserId,
		userAuthToken: data.AuthenticationToken,
	}

	return &data, err
}

func (a *AdminAPI) Logout() error {
	req, err := a.newRequest("POST", "/logout", nil)
	if err != nil {
		return err
	}

	_, err = a.do(req)
	return err
}

type AuthTokenSwitchAccountParams struct {
	AccountId string `json:"account_id"`
}

type AuthTokenSwitchAccountResponse struct {
	Object              string                 `mapstructure:"object"`
	AuthenticationToken string                 `mapstructure:"authentication_token"`
	UserId              string                 `mapstructure:"user_id"`
	User                map[string]interface{} `mapstructure:"user"`
	AccountId           string                 `mapstructure:"account_id"`
	Account             map[string]interface{} `mapstructure:"account"`
}

func (a *AdminAPI) AuthTokenSwitchAccount(reqBody AuthTokenSwitchAccountParams) (*AuthTokenSwitchAccountResponse, error) {
	req, err := a.newRequest("POST", "/auth_token.switch_account", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data AuthTokenSwitchAccountResponse
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type PasswordResetParams struct {
	Email       string `json:"email"`
	RedirectUrl string `json:"redirect_url"`
}

func (a *AdminAPI) PasswordReset(reqBody PasswordResetParams) error {
	req, err := a.newRequest("POST", "/password.reset", reqBody)
	if err != nil {
		return err
	}

	_, err = a.do(req)
	return err
}

type PasswordUpdateParams struct {
	Email                string `json:"email"`
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (a *AdminAPI) PasswordUpdate(reqBody PasswordUpdateParams) error {
	req, err := a.newRequest("POST", "/password.update", reqBody)
	if err != nil {
		return err
	}

	_, err = a.do(req)
	return err
}

/////////////////
// Minted Token
/////////////////
type MintedToken struct {
	Object        string `mapstructure:"object"`
	Id            string `mapstructure:"id"`
	Symbol        string `mapstructure:"symbol"`
	SubunitToUnit int    `mapstructure:"subunit_to_unit"`
	CreatedAt     string `mapstructure:"created_at"`
	UpdatedAt     string `mapstructure:"updated_at"`
}

type MintedTokenAllParams struct {
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	SearchTerm string `json:"search_term"`
	SortBy     string `json:"sort_by"`
	SortDir    string `json:"sort_dir"`
}

type MintedTokenAllResponse struct {
	Object     string                 `mapstructure:"object"`
	Data       []MintedToken          `mapstructure:"data"`
	Pagination map[string]interface{} `mapstructure:"pagination"`
}

func (a *AdminAPI) MintedTokenAll(reqBody MintedTokenAllParams) (*MintedTokenAllResponse, error) {
	req, err := a.newRequest("POST", "/minted_token.all", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data MintedTokenAllResponse
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type MintedTokenGetParams struct {
	Id string `json:"id"`
}

type MintedTokenGetResponse struct {
	MintedToken
}

func (a *AdminAPI) MintedTokenGet(reqBody MintedTokenGetParams) (*MintedTokenGetResponse, error) {
	req, err := a.newRequest("POST", "/minted_token.get", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data MintedTokenGetResponse
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type MintedTokenCreateParams struct {
	Name                 string                 `json:"name"`
	Symbol               string                 `json:"symbol"`
	Description          string                 `json:"description"`
	SubunitToUnit        int                    `json:"subunit_to_unit"`
	Amount               int                    `json:"amount"`
	IsoCode              string                 `json:"iso_code"`
	ShortSymbol          string                 `json:"short_symbol"`
	Subunit              string                 `json:"subunit"`
	SymbolFirst          bool                   `json:"symbol_first"`
	HtmlEntity           string                 `json:"html_entity"`
	IsoNumeric           string                 `json:"iso_numeric"`
	SmallestDenomination int                    `json:"smallest_denomination"`
	Metadata             map[string]interface{} `json:"id"`
	EncryptedMetadata    map[string]interface{} `json:"id"`
}

type MintedTokenCreateResponse struct {
	MintedToken
}

func (a *AdminAPI) MintedTokenCreate(reqBody MintedTokenCreateParams) (*MintedTokenCreateResponse, error) {
	req, err := a.newRequest("POST", "/minted_token.create", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data MintedTokenCreateResponse
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type MintedTokenMintParams struct {
	Id     string `json:"id"`
	Amount int    `json:"amount"`
}

type MintedTokenMintResponse struct {
	MintedToken
}

func (a *AdminAPI) MintedTokenMint(reqBody MintedTokenMintParams) (*MintedTokenMintResponse, error) {
	req, err := a.newRequest("POST", "/minted_token.mint", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data MintedTokenMintResponse
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

/////////////////
// Transaction
/////////////////

/////////////////
// API Access
/////////////////
type AccessKeyCreateResponse struct {
	Object    string `mapstructure:"object"`
	Id        string `mapstructure:"id"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	AccountId string `mapstructure:"account_id"`
	CreatedAt string `mapstructure:"created_at"`
	UpdatedAt string `mapstructure:"updated_at"`
	DeletedAt string `mapstructure:"deleted_at"`
}

func (a *AdminAPI) AccessKeyCreate() (*AccessKeyCreateResponse, error) {
	req, err := a.newRequest("POST", "/access_key.create", nil)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	var data AccessKeyCreateResponse
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}
