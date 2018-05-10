package omisego

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type EWalletAPI struct {
	*Client
}

/////////////////
// User
/////////////////
type UserCreateParams struct {
	ProviderUserId    string                 `json:"provider_user_id"`
	Username          string                 `json:"username"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	EncryptedMetadata map[string]interface{} `json:"encrypted_metadata,omitempty"`
}

type UserCreateResponse struct {
	Object            string                 `mapstructure:"object"`
	Id                string                 `mapstructure:"id"`
	ProviderUserId    string                 `mapstructure:"provider_user_id"`
	Username          string                 `mapstructure:"username"`
	Metadata          map[string]interface{} `mapstructure:"metadata"`
	EncryptedMetadata map[string]interface{} `mapstructure:"encrypted_metadata"`
}

func (e *EWalletAPI) UserCreate(reqBody UserCreateParams) (*UserCreateResponse, error) {
	req, err := e.newRequest("POST", "/user.create", reqBody)
	req.Header.Set("Idempotency-Token", NewIdempotencyToken())
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)

	var data UserCreateResponse
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}
