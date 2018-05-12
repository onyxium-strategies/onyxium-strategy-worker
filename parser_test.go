package main

import (
	// "github.com/stretchr/testify/assert"
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
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

	conditions, err := createConditionsFromSlice(data)
	if err != nil {
		t.Fatal(err)
	}
	expectedConditions := []models.Condition{{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.01990991}, {ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.02}}

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

	condition, err := createConditionFromMap(data)
	if err != nil {
		t.Fatal(err)
	}
	expectedCondition := models.Condition{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.01990991, TimeframeInMS: 3600000}

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

	action, err := createActionFromMap(data)
	if err != nil {
		t.Fatal(err)
	}
	expectedAction := models.Action{OrderType: "limit-buy", ValueType: "percentage-above", ValueQuoteMetric: "price-ask", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 10, Value: 0.02}

	if !reflect.DeepEqual(action, expectedAction) {
		t.Errorf("Expected condition to be %+v but it was %+v", expectedAction, action)
	}
}

// Right now the test is dependent on property OrderType of Action. If the property changes this test breaks.
// Have to find a way to test the tree traversal without using OrderType
func TestParseBinaryTree(t *testing.T) {
	var jsonInput = `[{"conditions":[{}],"action":{"OrderType": "A"},"then":[{"conditions":[{}],"action":{"OrderType": "B"},"then":[]}, {"conditions":[{}],"action":{"OrderType": "C"},"then":[]}]}, {"conditions":[{}],"action":{"OrderType": "D"},"then":[]}, {"conditions":[{}],"action":{"OrderType": "E"},"then":[]}]`
	// var jsonInput = `[{"then":[{"then":[]}, {"then":[]}]}, {"then":[]}, {"then":[]}]`
	data, err := parseJsonArray(jsonInput)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("%+v", data)
	tree, err := parseBinaryTree(data)
	if err != nil {
		t.Fatal(err)
	}
	// t.Logf("%+v", tree)
	size := 6
	ch := make(chan string, size)
	inOrderTraverse(t, &tree, ch)

	expectedCh := make(chan string, size)
	expectedCh <- "" //root
	expectedCh <- "A"
	expectedCh <- "B"
	expectedCh <- "C"
	expectedCh <- "D"
	expectedCh <- "E"
	close(ch)
	close(expectedCh)

	for i := range ch {
		if j := <-expectedCh; i != j {
			t.Fatalf("Incorrected tree structure expected node to be %s but it was %s", j, i)
		}
	}
}

func inOrderTraverse(t *testing.T, node *models.Tree, ch chan string) {
	if node != nil {
		ch <- node.Action.OrderType
		inOrderTraverse(t, node.Left, ch)
		inOrderTraverse(t, node.Right, ch)
	}
}
