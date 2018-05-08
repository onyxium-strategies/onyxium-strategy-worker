package main

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// Parse nested array json string to json object
func parseJsonArray(jsonInput string) ([]interface{}, error) {
	if !gjson.Valid(jsonInput) {
		return nil, errors.New("Invalid json structure")
	}

	data, ok := gjson.Parse(jsonInput).Value().([]interface{})
	if !ok {
		return nil, errors.New("Input is not a json array")
	}
	return data, nil
}

// Parse json tree from frontend to a binary tree for backend
func parseBinaryTree(tree []interface{}) Tree {
	// TODO can we exclude this root node from the tree?
	root := Tree{Left: nil, Right: nil, Conditions: []Condition{}, Action: Action{}}
	root.Left = _parseBinaryTree(tree, root.Left, 0)
	return root
}

// Recursive version of parseBinaryTree
func _parseBinaryTree(siblings []interface{}, root *Tree, i int) *Tree {
	if i < len(siblings) {
		conditions, ok := siblings[i].(map[string]interface{})["conditions"]
		if !ok {
			conditions = make([]interface{}, 0)
		}
		action, ok := siblings[i].(map[string]interface{})["action"]
		if !ok {
			action = make(map[string]interface{})
		}
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

// Decode json array of Condition structs
func createConditionsFromSlice(conditionsSlice []interface{}) []Condition {
	conditions := []Condition{}
	for _, condition := range conditionsSlice {
		conditions = append(conditions, createConditionFromMap(condition.(map[string]interface{})))
	}
	return conditions
}

// Decode json to Condition struct
func createConditionFromMap(m map[string]interface{}) Condition {
	var result Condition
	err := mapstructure.Decode(m, &result)
	if err != nil {
		log.Error(err)
	}
	return result
}

// Decode json to Action struct
func createActionFromMap(m map[string]interface{}) Action {
	var result Action
	err := mapstructure.Decode(m, &result)
	if err != nil {
		log.Error(err)
	}
	return result
}
