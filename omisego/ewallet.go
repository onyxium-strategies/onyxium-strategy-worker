package omisego

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type EWalletAPI struct {
	*Client
}

/////////////////
// Session
/////////////////
type ProviderUserIdParam struct {
	PrividerUserId string `json:"provider_user_id"`
}

func (e *EWalletAPI) Login(reqBody ProviderUserIdParam) (*AuthenicationToken, error) {
	req, err := e.newRequest("POST", "/login", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
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

func (e *EWalletAPI) Logout() error {
	req, err := e.newRequest("POST", "/logout", nil)
	if err != nil {
		return err
	}

	_, err = e.do(req)
	return err
}

/////////////////
// User
/////////////////
type UserParams struct {
	ProviderUserId    string                 `json:"provider_user_id"`
	Username          string                 `json:"username"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	EncryptedMetadata map[string]interface{} `json:"encrypted_metadata,omitempty"`
}

func (e *EWalletAPI) UserCreate(reqBody UserParams) (*User, error) {
	req, err := e.newRequest("POST", "/user.create", reqBody)
	req.Header.Set("Idempotency-Token", NewIdempotencyToken())
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)

	var data User
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

func (e *EWalletAPI) UserUpdate(reqBody UserParams) (*User, error) {
	req, err := e.newRequest("POST", "/user.update", reqBody)
	req.Header.Set("Idempotency-Token", NewIdempotencyToken())
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)

	var data User
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

func (e *EWalletAPI) UserGet(reqBody ProviderUserIdParam) (*User, error) {
	req, err := e.newRequest("POST", "/user.get", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)

	var data User
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

func (e *EWalletAPI) MeGet() (*User, error) {
	req, err := e.newRequest("POST", "/user.get", nil)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)

	var data User
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}
