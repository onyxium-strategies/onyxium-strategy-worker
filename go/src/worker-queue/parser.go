package main

import (
	"github.com/tidwall/gjson"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/mitchellh/mapstructure"
)

// TODO return err obj
func parseJson(jsonInput string) ([]interface{}, string) {
	data, ok := gjson.Parse(jsonInput).Value().([]interface{})
	if !ok {
		return data, "Invalid json structure"
	}
	return data, ""
}

/*
[map[then:[map[conditions:[map[quoteCurrency:OMG timeframeInMS:3.6e+06 baseMetric:price value:20 conditionType:percentage-increase baseCurrency:ETH] map[quoteCurrency:OMG timeframeInMS:3.6e+06 baseMetric:price quoteMetric:boilinger-band-upper-bound value:20 conditionType:percentage-increase baseCurrency:ETH]] action:map[quoteCurrency:OMG quantity:100 value:0.012 orderType:limit-buy orderValueType:absolute baseCurrency:ETH] priority:0] map[conditions:[map[quoteCurrency:OMG timeframeInMS:7.2e+06 baseMetric:rsi value:5 conditionType:absolute-increase baseCurrency:ETH]] action:map[quoteCurrency:OMG quantity:100 value:0.012 orderType:limit-sell orderValueType:absolute baseCurrency:ETH] priority:1]] priority:0 conditions:[map[conditionType:percentage-increase baseCurrency:ETH quoteCurrency:OMG timeframeInMS:3.6e+06 baseMetric:price value:20]] action:map[quoteCurrency:OMG quantity:100 value:0.012 orderType:limit-buy orderValueType:absolute baseCurrency:ETH]] map[conditions:[map[quoteCurrency:OMG timeframeInMS:7.2e+06 baseMetric:rsi value:5 conditionType:absolute-increase baseCurrency:ETH]] action:map[quoteCurrency:OMG quantity:100 value:0.012 orderType:limit-buy orderValueType:absolute baseCurrency:ETH] priority:1] map[priority:2 conditions:[map[timeframeInMS:7.2e+06 baseMetric:rsi value:10 conditionType:absolute-decrease baseCurrency:ETH quoteCurrency:OMG]] action:map[orderValueType:absolute baseCurrency:ETH quoteCurrency:OMG quantity:100 value:0.012 orderType:limit-sell]]]
*/
func parseBinaryTree(tree []interface{}) Tree {
    // conditions := tree[0].(map[string]interface{})["conditions"]
	// action := tree[0].(map[string]interface{})["action"]
    // root := Tree{Left: nil, Right: nil, Conditions: createConditionsFromSlice(conditions.([]interface{})), Action: createActionFromMap(action.(map[string]interface {})) }
    root := Tree{ Left: nil, Right: nil, Conditions: []Condition{{ConditionType: "true"}}, Action: Action{} }
    fmt.Println(root)
	root.Left = _parseBinaryTree(tree, &root.Left, 0)
	return root
}

// https://www.geeksforgeeks.org/construct-complete-binary-tree-given-array/
func _parseBinaryTree(siblings []interface{}, root *Tree, i int) Tree{
    if i < len(siblings) {
        conditions := siblings[i].(map[string]interface{})["conditions"]
        action := siblings[i].(map[string]interface{})["action"]
        newNode := Tree{Left: nil, Right: nil, Conditions: createConditionsFromSlice(conditions.([]interface{})), Action: createActionFromMap(action.(map[string]interface {})) }
        root := newNode

        if then, ok := siblings[i].(map[string]interface{})["then"]; ok {
            root.Left = _parseBinaryTree(then, &root.Left, 0)
        } else {
            root.Left = _parseBinaryTree([]int, &root.Left, 0)
        }

        tmp := i
        tmp++
        root.Right = _parseBinaryTree(siblings[i], &root.Right, tmp)
    }
    return root
}

func createConditionsFromSlice(conditionsSlice []interface{}) []Condition {
	conditions := []Condition{}
	for _, condition := range conditionsSlice {
		conditions = append(conditions, createConditionFromMap(condition.(map[string]interface {})))
	}
	return conditions
}

func createConditionFromMap(m map[string]interface{}) Condition {
    var result Condition
    err := mapstructure.Decode(m, &result)
    if err != nil {
    	fmt.Printf("Error: %v\n", err)
    }
    return result
}

func createActionFromMap(m map[string]interface{}) Action {
    var result Action
    err := mapstructure.Decode(m, &result)
    if err != nil {
    	fmt.Printf("Error: %v\n", err)
    }
    return result
}

func main() {
	file, e := ioutil.ReadFile("./tree-example.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	myJson := string(file)

	tree, err := parseJson(myJson)
	if err != "" {
		fmt.Printf(err)
	}
	root := parseBinaryTree(tree)
    fmt.Println(root)
}
