package main

import (
	"bytes"
	"net/http"
)

func eventUserSignUp() error {
	body := `{"text":"A new user has just signed up.","icon_emoji":":dancer::skin-tone-2:"}`
	_, err := http.Post("https://hooks.slack.com/services/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", "application/json", bytes.NewBuffer([]byte(body)))
	return err
}

func eventUserActivated() error {
	body := `{"text":"A new user has been activated.","icon_emoji":":dancer::skin-tone-5:"}`
	_, err := http.Post("https://hooks.slack.com/services/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", "application/json", bytes.NewBuffer([]byte(body)))
	return err
}
