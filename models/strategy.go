package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Strategy struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	JsonTree string        `json:"jsonTree" bson:"jsonTree"`
	BsonTree *Tree         `json:"bsonTree" bson:"bsonTree"`
	Status   string        `json:"status" bson:"status"`
	State    int           `json:"state" bson:"state"`
}

type Tree struct {
	Id         int
	Left       *Tree
	Right      *Tree
	Conditions []Condition
	Action     Action
}

type Condition struct {
	ConditionType string  `validate:"required,oneof=percentage-decrease percentage-increase greater-than-or-equal-to less-than-or-equal-to"`
	BaseCurrency  string  `validate:"required,nefield=QuoteCurrency"`
	QuoteCurrency string  `validate:"required",nefield=BaseCurrency`
	TimeframeInMS int     `validate:"omitempty,gt=0"`
	BaseMetric    string  `validate:"required,oneof=price-ask price-bid price-last volume"`
	Value         float64 `validate:"required,gte=0"`
}

type Action struct {
	OrderType        string  `validate:"required,oneof=limit-buy limit-sell"`
	ValueType        string  `validate:"required,oneof=absolute relative-above relative-below percentage-above percentage-below"`
	ValueQuoteMetric string  `validate:"omitempty,oneof=price-ask price-bid price-last"`
	BaseCurrency     string  `validate:"required,nefield=QuoteCurrency"`
	QuoteCurrency    string  `validate:"required,nefield=BaseCurrency"`
	Quantity         float64 `validate:"required,gt=0"`
	Value            float64 `validate:"required"gt=0`
}

func (db *MGO) StrategyCreate(name string, jsonTree string, bsonTree *Tree) error {
	strategy := Strategy{
		Id:       bson.NewObjectId(),
		Name:     name,
		JsonTree: jsonTree,
		BsonTree: bsonTree,
		Status:   "paused",
		State:    bsonTree.Id,
	}
	c := db.DB("coinflow").C("strategy")
	err := c.Insert(strategy)
	return err
}

func (db *MGO) GetPausedStrategies() ([]Strategy, error) {
	var strategies []Strategy
	c := db.DB("coinflow").C("strategy")
	err := c.Find(bson.M{"status": "paused"}).All(strategies)
	if err != nil {
		return nil, err
	}
	return strategies, nil
}
