package main

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmailConfirm(t *testing.T) {
	testCases := []struct {
		name             string
		id               string
		token            string
		shouldPass       bool
		expectedResponse string
	}{
		{"Correct params", "5af6d0de24979001180b8424", "$2a$04$vjyoz/l.OA1B4hne7x8rd.5l8h/oSldGB9oZwechaalB.0Z6spYuq", true, `{"success":true}`},
		{"Incorrect id", "5af6d0de24979001180b842", "$2a$04$vjyoz/l.OA1B4hne7x8rd.5l8h/oSldGB9oZwechaalB.0Z6spYuq", false, `{"error":"Incorrect id hex received: 5af6d0de24979001180b842"}`},
		{"Incorrect token", "5af6d0de24979001180b8424", "$2a$04$vjyoz/l.OA1B4hne7x8rd.5l8h/oSldGB9oZwechaalB.0Z6spYu", false, `{"error":"crypto/bcrypt: hashedPassword is not the hash of the given password"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/confirm-email", nil)
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"id": tc.id, "token": tc.token})

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(EmailConfirm)

			handler.ServeHTTP(rec, req)

			// In this case, our EmailConfirm returns a non-200 response
			// for a route variable it doesn't know about.
			if rec.Code == http.StatusOK && !tc.shouldPass {
				t.Errorf("handler should have failed on with id %s and token %s: got %v want %v",
					tc.id, tc.token, rec.Code, http.StatusOK)
			}

			// Check the response body is what we expect.
			expected := tc.expectedResponse
			if rec.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rec.Body.String(), expected)
			}
		})
	}
}

func TestUserAll(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(UserAll)

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check the response body is what we expect.
	body := make([]models.User, 2)
	json.NewDecoder(rec.Body).Decode(&body)
	expectedResponse := []models.User{
		{Email: "test@gmail.com", Password: "pwd"},
		{Email: "test2@gmail.com", Password: "pwd2"},
	}
	assert.Equal(t, expectedResponse, body)
}

func TestUserCreate(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/user", bytes.NewBuffer([]byte(`{"email": "example@gmail.com","password":"goedpassword"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(UserCreate)

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rec.Code)

	// Check the response body is what we expect.
	body := models.User{}
	json.NewDecoder(rec.Body).Decode(&body)

	// assert that the password is hashed
	assert.NotEqual(t, "goedpassword", body.Password)
}
