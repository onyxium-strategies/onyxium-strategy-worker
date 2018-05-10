package omisego

import ()

type EWalletAPI struct {
	Client
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

func (e *EWalletAPI) UserCreate(reqBody UserCreateParams) (Response, error) {
	req, err := e.newRequest("POST", "/user.create", reqBody)
	req.Header.Set("Idempotency-Token", newIdempotencyToken())
	if err != nil {
		return Response{}, err
	}

	var res Response
	_, err = e.do(req, &res)
	return res, err
}
