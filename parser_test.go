package main

import (
	// "github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"reflect"
	"testing"
)

func TestParseJsonArray(t *testing.T) {
	var jsonInput = `[{"name":{"first":"Janet","last":"Prichard"},"age":47}]`
	_, err := parseJsonArray(jsonInput)
	if err != nil {
		t.Fatalf("Not able to parse a valid nested array json string structure with error: %s.", err)
	}
}

func TestCreateConditionsFromSlice(t *testing.T) {
	var jsonInput = `[{"conditionType":"greater-than-or-equal-to","baseCurrency":"BTC","quoteCurrency":"ETH","baseMetric":"price-last","value":0.01990991},{"conditionType":"less-than-or-equal-to","baseCurrency":"BTC","quoteCurrency":"ETH","baseMetric":"price-last","value":0.02}]`

	data, ok := gjson.Parse(jsonInput).Value().([]interface{})
	if !ok {
		t.Fatal("Json input is not a slice")
	}

	conditions := createConditionsFromSlice(data)
	expectedConditions := []Condition{{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.01990991}, {ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.02}}

	if !reflect.DeepEqual(conditions, expectedConditions) {
		t.Errorf("Expected conditions to be %+v but it was %+v", expectedConditions, conditions)
	}
}

func TestCreateConditionFromMap(t *testing.T) {
	var jsonInput = `{"conditionType":"percentage-increase","baseCurrency":"BTC","quoteCurrency":"ETH","baseMetric":"price-last","value":0.01990991, "TimeframeInMS": 3600000}`

	data, ok := gjson.Parse(jsonInput).Value().(map[string]interface{})
	if !ok {
		t.Fatal("Json input is not a map")
	}

	condition := createConditionFromMap(data)
	expectedCondition := Condition{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.01990991, TimeframeInMS: 3600000}

	if !reflect.DeepEqual(condition, expectedCondition) {
		t.Errorf("Expected condition to be %+v but it was %+v", expectedCondition, condition)
	}
}

func TestCreateActionFromMap(t *testing.T) {
	var jsonInput = `{"orderType":"limit-buy","valueType":"percentage-above","valueQuoteMetric":"price-ask","baseCurrency":"BTC","quoteCurrency":"ETH","quantity":10,"value":0.02}`

	data, ok := gjson.Parse(jsonInput).Value().(map[string]interface{})
	if !ok {
		t.Fatal("Json input is not a map")
	}

	action := createActionFromMap(data)
	expectedAction := Action{OrderType: "limit-buy", ValueType: "percentage-above", ValueQuoteMetric: "price-ask", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 10, Value: 0.02}

	if !reflect.DeepEqual(action, expectedAction) {
		t.Errorf("Expected condition to be %+v but it was %+v", expectedAction, action)
	}
}

// TODO make jsonInput a valid tree structure
// Walk over the tree and see if the order is correct
// Dont care about condition and action structs only that the order is correct.
func TestParseBinaryTree(t *testing.T) {
	var jsonInput = `[{"conditions":[{}],"action":{},"then":[]}]`

	data, ok := gjson.Parse(jsonInput).Value().([]interface{})
	if !ok {
		t.Fatal("Json input is not a slice")
	}

	tree := parseBinaryTree(data)

	t.Logf("%+v", tree.Left)

}
