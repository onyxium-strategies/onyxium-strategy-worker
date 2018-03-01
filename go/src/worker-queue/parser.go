package main

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
)

// TODO return err obj
func parseJson(jsonInput string) ([]interface{}, error) {
	data, ok := gjson.Parse(jsonInput).Value().([]interface{})
	if !ok {
		err := errors.New("Invalid json structure")
		return data, err
	}
	return data, nil
}

func parseBinaryTree(tree []interface{}) Tree {
	root := Tree{Left: nil, Right: nil, Conditions: []Condition{{ConditionType: "true"}}, Action: Action{}}
	root.Left = _parseBinaryTree(tree, root.Left, 0)
	return root
}

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
		fmt.Println(tree)
		walk(tree.Left)
		walk(tree.Right)
	}
}

func readTreeExample() Tree {
	file, e := ioutil.ReadFile("./tree-example.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	myJson := string(file)

	tree, err := parseJson(myJson)
	if err != nil {
		fmt.Println(err)
	}
	root := parseBinaryTree(tree)
	// walk(&root)
	return root
}
