package main

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

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
func parseBinaryTree(tree []interface{}) (Tree, error) {
	// TODO can we exclude this root node from the tree?
	root := Tree{Left: nil, Right: nil, Conditions: []Condition{}, Action: Action{}}
	left, err := _parseBinaryTree(tree, root.Left, 0)
	if err != nil {
		return Tree{}, err
	}
	root.Left = left

	return root, nil
}

// Recursive version of parseBinaryTree
func _parseBinaryTree(siblings []interface{}, root *Tree, i int) (*Tree, error) {
	if i < len(siblings) {
		conditions, ok := siblings[i].(map[string]interface{})["conditions"]
		if !ok {
			conditions = make([]interface{}, 0)
		}
		action, ok := siblings[i].(map[string]interface{})["action"]
		if !ok {
			action = make(map[string]interface{})
		}

		conditionsInstance, err := createConditionsFromSlice(conditions.([]interface{}))
		if err != nil {
			return nil, err
		}
		actionInstance, err := createActionFromMap(action.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		root = &Tree{Left: nil, Right: nil, Conditions: conditionsInstance, Action: actionInstance}

		if then, ok := siblings[i].(map[string]interface{})["then"]; ok {
			left, err := _parseBinaryTree(then.([]interface{}), root.Left, 0)
			if err != nil {
				return nil, err
			}
			root.Left = left
		} else {
			empty := make([]interface{}, 0)
			left, err := _parseBinaryTree(empty, root.Left, 0)
			if err != nil {
				return nil, err
			}
			root.Left = left
		}

		i += 1

		right, err := _parseBinaryTree(siblings, root.Right, i)
		if err != nil {
			return nil, err
		}
		root.Right = right
	}
	return root, nil
}

// Decode json array of Condition structs
func createConditionsFromSlice(conditionsSlice []interface{}) ([]Condition, error) {
	conditions := []Condition{}
	for _, condition := range conditionsSlice {
		conditionInstance, err := createConditionFromMap(condition.(map[string]interface{}))
		if err != nil {
			return []Condition{}, err
		}
		conditions = append(conditions, conditionInstance)
	}
	return conditions, nil
}

// Decode json to Condition struct
func createConditionFromMap(m map[string]interface{}) (Condition, error) {
	var result Condition
	err := mapstructure.Decode(m, &result)
	if err != nil {
		log.Info(err)
		return Condition{}, err
	}

	validate = validator.New()
	err = validate.Struct(result)
	if err != nil {
		log.Info(err)
		return Condition{}, err
	}

	return result, nil
}

// Decode json to Action struct
func createActionFromMap(m map[string]interface{}) (Action, error) {
	var result Action
	err := mapstructure.Decode(m, &result)
	if err != nil {
		log.Error(err)
		return Action{}, err
	}

	validate = validator.New()
	err = validate.Struct(result)
	if err != nil {
		log.Info(err)
		return Action{}, err
	}
	return result, nil
}
