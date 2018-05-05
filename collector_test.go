package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCollectorHandler_ServeHTTP(t *testing.T) {
	// golden file test
	content, err := ioutil.ReadFile("tree-example.json")
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/work", bytes.NewBuffer(content))
	if err != nil {
		t.Fatalf("Request error: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := CollectorHandler{}

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

}
