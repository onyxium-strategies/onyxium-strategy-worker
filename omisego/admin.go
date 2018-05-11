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

// type LoginResponse struct {
// 	UserId              string                 `mapstructure:"user_id"`
// 	User                map[string]interface{} `mapstructure:"user"`
// 	Object              string                 `mapstructure:"object"`
// 	MasterAdmin         bool                   `mapstructure:"master_admin"`
// 	AuthenticationToken string                 `mapstructure:"authentication_token"`
// 	AccountId           string                 `mapstructure:"account_id"`
// 	Account             map[string]interface{} `mapstructure:"account"`
// }

func (a *AdminAPI) Login(reqBody LoginParams) (*AuthenicationToken, error) {
	req, err := a.newRequest("POST", "/login", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data AuthenicationToken
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

// type AuthTokenSwitchAccountResponse struct {
// 	Object              string                 `mapstructure:"object"`
// 	AuthenticationToken string                 `mapstructure:"authentication_token"`
// 	UserId              string                 `mapstructure:"user_id"`
// 	User                map[string]interface{} `mapstructure:"user"`
// 	AccountId           string                 `mapstructure:"account_id"`
// 	Account             map[string]interface{} `mapstructure:"account"`
// }

func (a *AdminAPI) AuthTokenSwitchAccount(reqBody AuthTokenSwitchAccountParams) (*AuthenicationToken, error) {
	req, err := a.newRequest("POST", "/auth_token.switch_account", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data AuthenicationToken
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

// type MintedTokenAllResponse struct {
// 	Object     string                 `mapstructure:"object"`
// 	Data       []MintedToken          `mapstructure:"data"`
// 	Pagination map[string]interface{} `mapstructure:"pagination"`
// }

func (a *AdminAPI) MintedTokenAll(reqBody ListParams) (*MintedTokenList, error) {
	req, err := a.newRequest("POST", "/minted_token.all", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data MintedTokenList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type MintedTokenGetParams struct {
	Id string `json:"id"`
}

// type MintedTokenGetResponse struct {
// 	MintedToken
// }

func (a *AdminAPI) MintedTokenGet(reqBody MintedTokenGetParams) (*MintedToken, error) {
	req, err := a.newRequest("POST", "/minted_token.get", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data MintedToken
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
	SubunitToUnit        int                    `json:"subunit_to_unit,omitempty"`
	Amount               int                    `json:"amount,omitempty"`
	IsoCode              string                 `json:"iso_code,omitempty"`
	ShortSymbol          string                 `json:"short_symbol,omitempty"`
	Subunit              string                 `json:"subunit,omitempty"`
	SymbolFirst          bool                   `json:"symbol_first,omitempty"`
	HtmlEntity           string                 `json:"html_entity,omitempty"`
	IsoNumeric           string                 `json:"iso_numeric,omitempty"`
	SmallestDenomination int                    `json:"smallest_denomination,omitempty"`
	Metadata             map[string]interface{} `json:"id,omitempty"`
	EncryptedMetadata    map[string]interface{} `json:"id,omitempty"`
}

// type MintedTokenCreateResponse struct {
// 	MintedToken
// }

func (a *AdminAPI) MintedTokenCreate(reqBody MintedTokenCreateParams) (*MintedToken, error) {
	req, err := a.newRequest("POST", "/minted_token.create", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data MintedToken
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

// type MintedTokenMintResponse struct {
// 	MintedToken
// }

func (a *AdminAPI) MintedTokenMint(reqBody MintedTokenMintParams) (*MintedToken, error) {
	req, err := a.newRequest("POST", "/minted_token.mint", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data MintedToken
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

/////////////////
// Account
/////////////////
func (a *AdminAPI) AccountAll(reqBody ListParams) (*AccountList, error) {
	req, err := a.newRequest("POST", "/account.all", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data AccountList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type AccountGetParams struct {
	Id string `json:"id"`
}

func (a *AdminAPI) AccountGet(reqBody AccountGetParams) (*Account, error) {
	req, err := a.newRequest("POST", "/account.get", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data Account
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type AccountCreateParams struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description,omitempty"`
	ParentId          string                 `json:"parent_id,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	EncryptedMetadata map[string]interface{} `json:"encrypted_metadata,omitempty"`
}

func (a *AdminAPI) AccountCreate(reqBody AccountCreateParams) (*Account, error) {
	req, err := a.newRequest("POST", "/account.create", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data Account
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type AccountUpdateParams struct {
	Id                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	EncryptedMetadata map[string]interface{} `json:"encrypted_metadata,omitempty"`
}

func (a *AdminAPI) AccountUpdate(reqBody AccountUpdateParams) (*Account, error) {
	req, err := a.newRequest("POST", "/account.update", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data Account
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type AccountUploadAvatarParams struct {
	Id     string `json:"id"`
	Avatar string `json:"avatar"`
}

func (a *AdminAPI) AccountUploadAvatar(reqBody AccountUploadAvatarParams) (*Account, error) {
	req, err := a.newRequest("POST", "/account.upload_avatar", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data Account
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type AccountListUsersParams struct {
	AccountId string `json:"account_id"`
}

func (a *AdminAPI) AccountListUsers(reqBody AccountListUsersParams) (*UserList, error) {
	req, err := a.newRequest("POST", "/account.list_users", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := a.do(req)
	if err != nil {
		return nil, err
	}

	var data UserList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}

	return &data, err
}

type AccountAssignUserParams struct {
	UserId      string `json:"user_id,omitempty"`
	AccountId   string `json:"account_id"`
	RoleName    string `json:"role_name"`
	RedirectUrl string `json:"redirect_url,omitempty"`
	Email       string `json:"email,omitempty"`
}

func (a *AdminAPI) AccountAssignUser(reqBody AccountAssignUserParams) error {
	req, err := a.newRequest("POST", "/account.assign_user", reqBody)
	if err != nil {
		return err
	}

	_, err = a.do(req)
	return err
}

type AccountUnassignUserParams struct {
	UserId    string `json:"user_id"`
	AccountId string `json:"account_id"`
}

func (a *AdminAPI) AccountUnassignUser(reqBody AccountUnassignUserParams) error {
	req, err := a.newRequest("POST", "/account.unassign_user", reqBody)
	if err != nil {
		return err
	}

	_, err = a.do(req)
	return err
}

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
