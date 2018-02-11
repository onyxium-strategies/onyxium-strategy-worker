package main

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
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
[
    map[
        conditions:[
            map[
                value:20
                conditionType:percentage-increase
                baseCurrency:ETH
                quoteCurrency:OMG
                timeframeInMS:3.6e+06
                baseMetric:price
            ]
        ]
        action:
            map[
                value:0.012
                orderType:limit-buy
                orderValueType:absolute
                baseCurrency:ETH
                quoteCurrency:OMG
                quantity:100
            ]
        then:[
            map[
                conditions:[
                    map[
                        timeframeInMS:3.6e+06
                        baseMetric:price
                        value:20
                        conditionType:percentage-increase
                        baseCurrency:ETH
                        quoteCurrency:OMG
                    ]
                    map[
                        baseCurrency:ETH
                        quoteCurrency:OMG
                        timeframeInMS:3.6e+06
                        baseMetric:price
                        quoteMetric:boilinger-band-upper-bound
                        value:20
                        conditionType:percentage-increase
                    ]
                ]
                action:
                    map[
                        value:0.012
                        orderType:limit-buy
                        orderValueType:absolute
                        baseCurrency:ETH
                        quoteCurrency:OMG
                        quantity:100
                    ]
            ]
            map[conditions:[map[quoteCurrency:OMG timeframeInMS:7.2e+06 baseMetric:rsi value:5 conditionType:absolute-increase baseCurrency:ETH]] action:map[quoteCurrency:OMG quantity:100 value:0.012 orderType:limit-sell orderValueType:absolute baseCurrency:ETH]]]] map[conditions:[map[timeframeInMS:7.2e+06 baseMetric:rsi value:5 conditionType:absolute-increase baseCurrency:ETH quoteCurrency:OMG]] action:map[orderValueType:absolute baseCurrency:ETH quoteCurrency:OMG quantity:100 value:0.012 orderType:limit-buy]] map[conditions:[map[timeframeInMS:7.2e+06 baseMetric:rsi value:10 conditionType:absolute-decrease baseCurrency:ETH quoteCurrency:OMG]] action:map[orderType:limit-sell orderValueType:absolute baseCurrency:ETH quoteCurrency:OMG quantity:100 value:0.012]]]
*/
func parseBinaryTree(tree []interface{}) Tree {
	root := Tree{Left: nil, Right: nil, Conditions: []Condition{{ConditionType: "true"}}, Action: Action{}}
	root.Left = _parseBinaryTree(tree, root.Left, 0)
	return root
}

// for debugging
// func insert(root *Tree, i int) *Tree {
// 	if false {
// 		root = &Tree{Left: nil, Right: nil, Conditions: []Condition{{ConditionType: "false"}}, Action: Action{}}
// 	}
// 	return root
// }

func _parseBinaryTree(siblings []interface{}, root *Tree, i int) *Tree {
	if i < len(siblings) {
		conditions := siblings[i].(map[string]interface{})["conditions"]
		action := siblings[i].(map[string]interface{})["action"]
		root = &Tree{Left: nil, Right: nil, Conditions: createConditionsFromSlice(conditions.([]interface{})), Action: createActionFromMap(action.(map[string]interface{}))}

		if then, ok := siblings[i].(map[string]interface{})["then"]; ok {
			root.Left = _parseBinaryTree(then.([]interface{}), root.Left, 0)
		} else {
			empty := make([]interface{}, 0)
			root.Left = _parseBinaryTree(empty, root.Left, 0)
		}

		i += 1
		root.Right = _parseBinaryTree(siblings, root.Right, i)
	}
	return root
}

func createConditionsFromSlice(conditionsSlice []interface{}) []Condition {
	conditions := []Condition{}
	for _, condition := range conditionsSlice {
		conditions = append(conditions, createConditionFromMap(condition.(map[string]interface{})))
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

// for debugging the tree
func walk(tree *Tree) {
	if tree != nil {
		fmt.Println(tree.Action)
		walk(tree.Left)
		walk(tree.Right)
	}
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
	walk(&root)
}
