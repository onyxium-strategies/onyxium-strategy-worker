package omisego

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	// log "github.com/sirupsen/logrus"
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
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}

	var data User
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

func (e *EWalletAPI) MeGet() (*User, error) {
	req, err := e.newRequest("POST", "/me.get", nil)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data User
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

/////////////////
// Balance
/////////////////
type AddressList struct {
	Data []Address `mapstructure:"data"`
}

func (e *EWalletAPI) UserListBalances(reqBody ProviderUserIdParam) (*AddressList, error) {
	req, err := e.newRequest("POST", "/user.list_balances", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data AddressList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

type BalanceAdjustmentParams struct {
	PrividerUserId        string                 `json:"provider_user_id"`
	TokenId               string                 `json:"token_id"`
	Amount                int                    `json:"amount"`
	AccountId             string                 `json:"account_id,omitempty"`
	BurnBalanceIdentifier string                 `json:"burn_balance_identifier,omitempty"`
	Metadata              map[string]interface{} `json:"id,omitempty"`
	EncryptedMetadata     map[string]interface{} `json:"id,omitempty"`
}

func (e *EWalletAPI) UserCreditBalance(reqBody BalanceAdjustmentParams) (*AddressList, error) {
	req, err := e.newRequest("POST", "/user.credit_balance", reqBody)
	req.Header.Set("Idempotency-Token", NewIdempotencyToken())
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data AddressList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

func (e *EWalletAPI) UserDebitBalance(reqBody BalanceAdjustmentParams) (*AddressList, error) {
	req, err := e.newRequest("POST", "/user.debit_balance", reqBody)
	req.Header.Set("Idempotency-Token", NewIdempotencyToken())
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data AddressList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

type TransferParams struct {
	FromAddress       string                 `json:"from_address"`
	ToAddress         string                 `json:"to_address"`
	TokenId           string                 `json:"token_id"`
	Amount            int                    `json:"amount"`
	Metadata          map[string]interface{} `json:"id,omitempty"`
	EncryptedMetadata map[string]interface{} `json:"id,omitempty"`
}

func (e *EWalletAPI) Transfer(reqBody TransferParams) (*AddressList, error) {
	req, err := e.newRequest("POST", "/transfer", reqBody)
	req.Header.Set("Idempotency-Token", NewIdempotencyToken())
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data AddressList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

/////////////////
// Settings
/////////////////
func (e *EWalletAPI) GetSettings() (*Settings, error) {
	req, err := e.newRequest("POST", "/get_settings", nil)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data Settings
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

/////////////////
// Transaction
/////////////////
func (e *EWalletAPI) TransactionAll(reqBody ListParams) (*TransactionList, error) {
	req, err := e.newRequest("POST", "/transaction.all", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data TransactionList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

type UserListTransactionsParams struct {
	ProviderUserId string                 `json:"provider_user_id"`
	Address        string                 `json:"address"`
	Page           int                    `json:"page,omitempty"`
	PerPage        int                    `json:"per_page,omitempty"`
	SearchTerm     string                 `json:"search_term,omitempty"`
	SearchTerms    map[string]interface{} `json:"search_terms,omitempty"`
	SortBy         string                 `json:"sort_by,omitempty"`
	SortDir        string                 `json:"sort_dir,omitempty"`
}

func (e *EWalletAPI) UserListTransactions(reqBody UserListTransactionsParams) (*TransactionList, error) {
	req, err := e.newRequest("POST", "/user.list_transactions", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data TransactionList
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

/////////////////
// Transaction Request
/////////////////
