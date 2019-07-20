package main

import (
	"bytes"
	"net/http"
)

func eventUserSignUp() error {
	body := `{"text":"A new user has just signed up.","icon_emoji":":golang:","username":"Go"}`
	_, err := http.Post("", "application/json", bytes.NewBuffer([]byte(body)))
	return err
}

func eventUserActivated() error {
	body := `{"text":"A new user has been activated.","icon_emoji":":golang:","username":"Go"}`
	_, err := http.Post("", "application/json", bytes.NewBuffer([]byte(body)))
	return err
}
