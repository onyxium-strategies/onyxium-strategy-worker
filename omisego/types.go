package omisego

import ()

type (
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

	TransactionSource struct {
		Address     string `mapstructure:"address"`
		Amount      int    `mapstructure:"amount"`
		MintedToken `mapstructure:"minted_token"`
	}

	Exchange struct {
		Rate int `mapstructure:"rate"`
	}

	TransactionList struct {
		Data       []Transaction `mapstructure:"data"`
		Pagination `mapstructure:"pagination"`
	}

	Pagination struct {
		PerPage     int  `mapstructure:"per_page"`
		CurrentPage int  `mapstructure:"current_page"`
		IsFirstPage bool `mapstructure:"is_first_page"`
		IsLastPage  bool `mapstructure:"is_last_page"`
	}
)
