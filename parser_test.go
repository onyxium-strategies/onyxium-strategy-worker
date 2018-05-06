package main

import (
	"testing"
)

func TestParseJsonArray(t *testing.T) {
	var jsonInput = `[{"name":{"first":"Janet","last":"Prichard"},"age":47}]`
	_, err := parseJsonArray(jsonInput)
	if err != nil {
		t.Fatal("Not able to parse a valid nested array json string structure.")
	}
}
