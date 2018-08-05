package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollector(t *testing.T) {
	// golden file test
	content, err := ioutil.ReadFile("strategy-example.json")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/work", bytes.NewBuffer(content))
	if err != nil {
		t.Fatalf("Request error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(NewStrategyCollector)
	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusCreated {
		body, _ := ioutil.ReadAll(rec.Body)
		t.Errorf("handler returned wrong status code: got %v want %v with response %s.",
			status, http.StatusCreated, string(body))
	}
}
