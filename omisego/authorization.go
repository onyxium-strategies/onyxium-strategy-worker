package omisego

import (
	"encoding/base64"
)

type Authorization interface {
	CreateAuthorizationHeader() string
}

type AdminClientAuth struct {
	apiKeyId string
	apiKey   string
}

type AdminUserAuth struct {
	apiKeyId      string
	apiKey        string
	userId        string
	userAuthToken string
}

type ServerAuth struct {
	accessKey string
	secretKey string
}

func (a *AdminClientAuth) CreateAuthorizationHeader() string {
	data := []byte(a.apiKeyId + ":" + a.apiKey)
	str := base64.StdEncoding.EncodeToString(data)
	return "OMGAdmin " + str
}

func (a *AdminUserAuth) CreateAuthorizationHeader() string {
	data := []byte(a.apiKeyId + ":" + a.apiKey + ":" + a.userId + ":" + a.userAuthToken)
	str := base64.StdEncoding.EncodeToString(data)
	return "OMGAdmin " + str
}

func (s *ServerAuth) CreateAuthorizationHeader() string {
	data := []byte(s.accessKey + ":" + s.secretKey)
	str := base64.StdEncoding.EncodeToString(data)
	return "OMGServer " + str
}
