package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmailConfirm(t *testing.T) {
	testCases := []struct {
		id               string
		token            string
		shouldPass       bool
		expectedResponse string
	}{
		{"5af6d0de24979001180b8424", "$2a$04$vjyoz/l.OA1B4hne7x8rd.5l8h/oSldGB9oZwechaalB.0Z6spYuq", true, `{"success":true}`},
		{"5af6d0de24979001180b842", "$2a$04$vjyoz/l.OA1B4hne7x8rd.5l8h/oSldGB9oZwechaalB.0Z6spYuq", false, `{"error":"Incorrect id hex received: 5af6d0de24979001180b842"}`},
		{"5af6d0de24979001180b8424", "$2a$04$vjyoz/l.OA1B4hne7x8rd.5l8h/oSldGB9oZwechaalB.0Z6spYu", false, `{"error":"crypto/bcrypt: hashedPassword is not the hash of the given password"}`},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("POST", "/api/confirm-email", nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": tc.id, "token": tc.token})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(EmailConfirm)

		handler.ServeHTTP(rr, req)

		// In this case, our EmailConfirm returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code == http.StatusOK && !tc.shouldPass {
			t.Errorf("handler should have failed on with id %s and token %s: got %v want %v",
				tc.id, tc.token, rr.Code, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := tc.expectedResponse
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}
}
