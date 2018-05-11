package omisego

import (
	"fmt"
)

type (
	BaseResponse struct {
		Version string                 `json:"version"`
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}

	ErrorResponse struct {
		Code        string            `json:"code"`
		Description string            `json:"description"`
		Messages    map[string]string `json:"messages"`
	}

	AuthenicationToken struct {
		UserId              string `mapstructure:"user_id"`
		User                `mapstructure:"user"`
		MasterAdmin         bool   `mapstructure:"master_admin"`
		AuthenticationToken string `mapstructure:"authentication_token"`
		AccountId           string `mapstructure:"account_id"`
		Account             `mapstructure:"account"`
	}

	User struct {
		Id                string                 `mapstructure:"id"`
		ProviderUserId    string                 `mapstructure:"provider_user_id"`
		Username          string                 `mapstructure:"username"`
		Email             string                 `mapstructure:"email"`
		Metadata          map[string]interface{} `mapstructure:"metadata"`
		EncryptedMetadata map[string]interface{} `mapstructure:"encrypted_metadata"`
		Avatar            map[string]interface{} `mapstructure:"avatar"`
		CreatedAt         string                 `mapstructure:"created_at"`
		UpdatedAt         string                 `mapstructure:"updated_at"`
	}

	UserList struct {
		Data       []User `mapstructure:"data"`
		Pagination `mapstructure:"pagination"`
	}

	Account struct {
		Id                string                 `mapstructure:"id"`
		ParentId          string                 `mapstructure:"parent_id"`
		Name              string                 `mapstructure:"name"`
		Description       string                 `mapstructure:"description"`
		Master            bool                   `mapstructure:"master"`
		Metadata          map[string]interface{} `mapstructure:"metadata"`
		EncryptedMetadata map[string]interface{} `mapstructure:"encrypted_metadata"`
		Avatar            map[string]interface{} `mapstructure:"avatar"`
		CreatedAt         string                 `mapstructure:"created_at"`
		UpdatedAt         string                 `mapstructure:"updated_at"`
	}

	AccountList struct {
		Data       []Account `mapstructure:"data"`
		Pagination `mapstructure:"pagination"`
	}

	MintedToken struct {
		Id                string                 `mapstructure:"id"`
		Symbol            string                 `mapstructure:"symbol"`
		Name              string                 `mapstructure:"name"`
		SubunitToUnit     int                    `mapstructure:"subunit_to_unit"`
		CreatedAt         string                 `mapstructure:"created_at"`
		UpdatedAt         string                 `mapstructure:"updated_at"`
		Metadata          map[string]interface{} `mapstructure:"metadata"`
		EncryptedMetadata map[string]interface{} `mapstructure:"encrypted_metadata"`
	}

	MintedTokenList struct {
		Data       []MintedToken `mapstructure:"data"`
		Pagination `mapstructure:"pagination"`
	}

	Transaction struct {
		Id                string            `mapstructure:"id"`
		From              TransactionSource `mapstructure:"from"`
		To                TransactionSource `mapstructure:"to"`
		Exchange          `mapstructure:"exchange"`
		Metadata          map[string]interface{} `mapstructure:"metadata"`
		EncryptedMetadata map[string]interface{} `mapstructure:"encrypted_metadata"`
		Status            string                 `mapstructure:"status"`
		CreatedAt         string                 `mapstructure:"created_at"`
		UpdatedAt         string                 `mapstructure:"updated_at"`
	}

	TransactionList struct {
		Data       []Transaction `mapstructure:"data"`
		Pagination `mapstructure:"pagination"`
	}

	TransactionSource struct {
		Address     string `mapstructure:"address"`
		Amount      int    `mapstructure:"amount"`
		MintedToken `mapstructure:"minted_token"`
	}

	Exchange struct {
		Rate int `mapstructure:"rate"`
	}

	Pagination struct {
		PerPage     int  `mapstructure:"per_page"`
		CurrentPage int  `mapstructure:"current_page"`
		IsFirstPage bool `mapstructure:"is_first_page"`
		IsLastPage  bool `mapstructure:"is_last_page"`
	}

	AccessKey struct {
		Id        string `mapstructure:"id"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		AccountId string `mapstructure:"account_id"`
		CreatedAt string `mapstructure:"created_at"`
		UpdatedAt string `mapstructure:"updated_at"`
		DeletedAt string `mapstructure:"deleted_at"`
	}

	AccessKeyList struct {
		Data       []AccessKey `mapstructure:"data"`
		Pagination `mapstructure:"pagination"`
	}

	APIKey struct {
		Id        string `mapstructure:"id"`
		Key       string `mapstructure:"key"`
		OwnerApp  string `mapstructure:"owner_app"`
		AccountId string `mapstructure:"account_id"`
		CreatedAt string `mapstructure:"created_at"`
		UpdatedAt string `mapstructure:"updated_at"`
		DeletedAt string `mapstructure:"deleted_at"`
	}

	APIKeyList struct {
		Data       []APIKey `mapstructure:"data"`
		Pagination `mapstructure:"pagination"`
	}

	////////////
	// Request body parameters
	////////////
	ListParams struct {
		Page       int    `json:"page,omitempty"`
		PerPage    int    `json:"per_page,omitempty"`
		SearchTerm string `json:"search_term,omitempty"`
		SortBy     string `json:"sort_by,omitempty"`
		SortDir    string `json:"sort_dir,omitempty"`
	}

	GetByIdParam struct {
		Id string `json:"id"`
	}
)

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%+v", *e)
}
