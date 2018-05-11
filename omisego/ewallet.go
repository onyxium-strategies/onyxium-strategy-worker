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
	ProviderUserId string `json:"provider_user_id"`
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
	ProviderUserId        string                 `json:"provider_user_id"`
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
type ServerCreateTransactionRequestParams struct {
	Type                string                 `json:"type"`
	TokenId             string                 `json:"token_id"`
	Amount              int                    `json:"amount,omitempty"`
	CorrelationId       string                 `json:"correlation_id,omitempty"`
	AccountId           string                 `json:"account_id,omitempty"`
	ProviderUserId      string                 `json:"provider_user_id,omitempty"`
	Address             string                 `json:"address,omitempty"`
	RequireConfirmation bool                   `json:"require_confirmation,omitempty"`
	MaxConsumptions     int                    `json:"max_consumptions,omitempty"`
	ConsumptionLifetime int                    `json:"consumption_lifetime,omitempty"`
	ExpirationDate      string                 `json:"expiration_date,omitempty"`
	AllowAmountOverride bool                   `json:"allow_amount_override,omitempty"`
	Metadata            map[string]interface{} `json:"id,omitempty"`
	EncryptedMetadata   map[string]interface{} `json:"id,omitempty"`
}

type TransactionRequest struct {
	Version             string                 `mapstructure:"version"`
	Success             bool                   `mapstructure:"success"`
	Data                map[string]interface{} `mapstructure:"data"`
	Id                  string                 `mapstructure:"id"`
	SocketTopic         string                 `mapstructure:"socket_topic"`
	Type                string                 `mapstructure:"type"`
	Amount              string                 `mapstructure:"amount"`
	Status              string                 `mapstructure:"status"`
	CorrelationId       string                 `mapstructure:"correlation_id"`
	MintedTokenId       string                 `mapstructure:"minted_token_id"`
	MintedToken         map[string]interface{} `mapstructure:"minted_token"`
	AccountId           string                 `mapstructure:"account_id"`
	UserId              string                 `mapstructure:"user_id"`
	Address             string                 `mapstructure:"address"`
	RequireConfirmation bool                   `mapstructure:"require_confirmation"`
	MaxConsumptions     int                    `mapstructure:"max_consumptions"`
	ConsumptionLifetime int                    `mapstructure:"consumption_lifetime"`
	ExpirationReason    string                 `mapstructure:"expiration_reason"`
	ExpirationDate      string                 `mapstructure:"expiration_date"`
	AllowAmountOverride bool                   `mapstructure:"allow_amount_override"`
	Metadata            map[string]interface{} `mapstructure:"id"`
	EncryptedMetadata   map[string]interface{} `mapstructure:"id"`
	CreatedAt           string                 `mapstructure:"created_at"`
	UpdatedAt           string                 `mapstructure:"updated_at"`
	ExpiredAt           string                 `mapstructure:"expired_at"`
}

func (e *EWalletAPI) TransactionRequestCreate(reqBody ServerCreateTransactionRequestParams) (*TransactionRequest, error) {
	req, err := e.newRequest("POST", "/transaction_request.create", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data TransactionRequest
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

func (e *EWalletAPI) TransactionRequestGet(reqBody ByIdParam) (*TransactionRequest, error) {
	req, err := e.newRequest("POST", "/transaction_request.get", reqBody)
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data TransactionRequest
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

type ServerTransactionRequestConsumeParams struct {
	TransactionRequestId string                 `json:"transaction_request_id"`
	TokenId              string                 `json:"token_id,omitempty"`
	Amount               int                    `json:"amount,omitempty"`
	CorrelationId        string                 `json:"correlation_id,omitempty"`
	AccountId            string                 `json:"account_id,omitempty"`
	ProviderUserId       string                 `json:"provider_user_id,omitempty"`
	Address              string                 `json:"address,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	EncryptedMetadata    map[string]interface{} `json:"encrypted_metadata,omitempty"`
}

type TransactionComsumption struct {
	Version              string                 `mapstructure:"version"`
	Success              bool                   `mapstructure:"success"`
	Data                 map[string]interface{} `mapstructure:"data"`
	Id                   string                 `mapstructure:"id"`
	SocketTopic          string                 `mapstructure:"socket_topic"`
	Amount               string                 `mapstructure:"amount"`
	Status               string                 `mapstructure:"status"`
	CorrelationId        string                 `mapstructure:"correlation_id"`
	MintedTokenId        string                 `mapstructure:"minted_token_id"`
	MintedToken          map[string]interface{} `mapstructure:"minted_token"`
	IdempotencyToken     string                 `mapstructure:"idempotency_token"`
	TransactionId        string                 `mapstructure:"transaction_id"`
	Transaction          map[string]interface{} `mapstructure:"transaction"`
	UserId               string                 `mapstructure:"user_id"`
	User                 map[string]interface{} `mapstructure:"user"`
	AccountId            string                 `mapstructure:"account_id"`
	Account              map[string]interface{} `mapstructure:"account"`
	TransactionRequestId string                 `mapstructure:"transaction_request_id"`
	TransactionRequest   map[string]interface{} `mapstructure:"transaction_request"`
	Address              string                 `mapstructure:"address"`
	Metadata             map[string]interface{} `mapstructure:"metadata"`
	EncryptedMetadata    map[string]interface{} `mapstructure:"encrypted_metadata"`
	CreatedAt            string                 `mapstructure:"created_at"`
	UpdatedAt            string                 `mapstructure:"updated_at"`
	ExpiredAt            string                 `mapstructure:"expired_at"`
	ApprovedAt           string                 `mapstructure:"approved_at"`
	RejectedAt           string                 `mapstructure:"rejected_at"`
	ConfirmedAt          string                 `mapstructure:"confirmed_at"`
	FailedAt             string                 `mapstructure:"failed_at"`
}

func (e *EWalletAPI) TransactionRequestConsume(reqBody ServerTransactionRequestConsumeParams) (*TransactionComsumption, error) {
	req, err := e.newRequest("POST", "/transaction_request.consume", reqBody)
	req.Header.Set("Idempotency-Token", NewIdempotencyToken())
	if err != nil {
		return nil, err
	}

	res, err := e.do(req)
	if err != nil {
		return nil, err
	}

	var data TransactionComsumption
	err = mapstructure.Decode(res.Data, &data)
	if err != nil {
		return nil, fmt.Errorf("Something went wrong with decoding %+v to %T", res.Data, data)
	}
	return &data, nil
}

/////////////////
// Transaction Consumption
/////////////////
