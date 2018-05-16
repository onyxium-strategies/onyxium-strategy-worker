package main

import (
	"bitbucket.org/onyxium/onyxium-strategy-worker/models"
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
func parseBinaryTree(tree []interface{}) (*models.Tree, error) {
	// TODO can we exclude this root node from the tree?
	// root := &models.Tree{Left: nil, Right: nil, Conditions: []models.Condition{}, Action: models.Action{}}
	root, err := _parseBinaryTree(tree, nil, 0)
	if err != nil {
		return nil, err
	}
	// root.Left = left

	return root, nil
}

// Recursive version of parseBinaryTree
func _parseBinaryTree(siblings []interface{}, root *models.Tree, i int) (*models.Tree, error) {
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
		root = &models.Tree{Left: nil, Right: nil, Conditions: conditionsInstance, Action: actionInstance}

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

// Decode json array of models.Condition structs
func createConditionsFromSlice(conditionsSlice []interface{}) ([]models.Condition, error) {
	conditions := []models.Condition{}
	for _, condition := range conditionsSlice {
		conditionInstance, err := createConditionFromMap(condition.(map[string]interface{}))
		if err != nil {
			return []models.Condition{}, err
		}
		conditions = append(conditions, conditionInstance)
	}
	return conditions, nil
}

// Decode json to models.Condition struct
func createConditionFromMap(m map[string]interface{}) (models.Condition, error) {
	var result models.Condition
	err := mapstructure.Decode(m, &result)
	if err != nil {
		log.Info(err)
		return models.Condition{}, err
	}

	result.Validate()

	return result, nil
}

// Decode json to Action struct
func createActionFromMap(m map[string]interface{}) (models.Action, error) {
	var result models.Action
	err := mapstructure.Decode(m, &result)
	if err != nil {
		log.Error(err)
		return models.Action{}, err
	}

	result.Validate()

	return result, nil
}
