package models

import (
	"fmt"
	// "github.com/tidwall/gjson"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
)

const StrategyCollection = "strategy"

type Strategy struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name" bson:"name"`
	JsonTree []interface{} `json:"jsonTree" bson:"jsonTree"`
	BsonTree *Tree         `json:"bsonTree" bson:"bsonTree"`
	Status   string        `json:"status" bson:"status"`
	State    int           `json:"state" bson:"state"`
	UserId   bson.ObjectId `json:"userId" bson:"userId"`
}

func NewStrategy(name, userId string, jsonTree []interface{}, bsonTree *Tree) (*Strategy, error) {
	// if !gjson.Valid(jsonTree) {
	// 	return nil, fmt.Errorf("Invalid json structure received: %s", jsonTree)
	// }
	strategy := &Strategy{
		Id:       bson.NewObjectId(),
		Name:     name,
		JsonTree: jsonTree,
		BsonTree: bsonTree,
		Status:   "paused",
		State:    bsonTree.Id,
		UserId:   bson.ObjectIdHex(userId),
	}
	return strategy, nil
}

type Tree struct {
	Id         int
	Left       *Tree
	Right      *Tree
	Conditions []Condition
	Action     Action
}

func (t *Tree) SetIdsForBinarySearch() {
	n := countNodes(t)
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
	inOrderTraverseReadChan(t, ch)
}

func countNodes(t *Tree) int {
	if t == nil {
		return 0
	} else {
		return 1 + countNodes(t.Left) + countNodes(t.Right)
	}
}

func inOrderTraverseReadChan(t *Tree, ch chan int) {
	if t != nil {
		inOrderTraverseReadChan(t.Left, ch)
		t.Id = <-ch
		inOrderTraverseReadChan(t.Right, ch)

	}
}

// TODO should return the string representation instead of print it
// String prints a visual representation of the tree
func (t *Tree) String() {
	fmt.Println("------------------------------------------------")
	stringify(t, 0)
	fmt.Println("------------------------------------------------")
}

// internal recursive function to print a tree
func stringify(t *Tree, level int) {
	if t != nil {
		format := ""
		for i := 0; i < level; i++ {
			format += "       "
		}
		format += "---[ "
		level++

		stringify(t.Right, level)
		fmt.Printf(format+"%d\n", t.Id)
		stringify(t.Left, level)

	}
}

// Search returns true if the Item t exists in the tree
func (t *Tree) Search(id int) (*Tree, error) {
	return search(t, id)
}

// internal recursive function to search an item in the tree
func search(t *Tree, id int) (*Tree, error) {
	if t == nil {
		return nil, fmt.Errorf("Did not find node with ID %d", id)
	}
	if id < t.Id {
		return search(t.Left, id)
	}
	if id > t.Id {
		return search(t.Right, id)
	}
	return t, nil
}

type Condition struct {
	ConditionType string  `validate:"required,oneof=percentage-decrease percentage-increase greater-than-or-equal-to less-than-or-equal-to"`
	BaseCurrency  string  `validate:"required,nefield=QuoteCurrency"`
	QuoteCurrency string  `validate:"required",nefield=BaseCurrency`
	TimeframeInMS int     `validate:"omitempty,gt=0"`
	BaseMetric    string  `validate:"required,oneof=price-ask price-bid price-last volume"`
	Value         float64 `validate:"required,gte=0"`
}

func (c *Condition) Validate() error {
	validate := validator.New()
	validate.RegisterStructValidation(customConditionValidation, Condition{})
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	return nil
}

// Helper for condition Validate()
func customConditionValidation(sl validator.StructLevel) {

	condition := sl.Current().Interface().(Condition)
	if (condition.ConditionType == "percentage-decrease") || (condition.ConditionType == "percentage-increase") {
		if condition.TimeframeInMS == 0 {
			sl.ReportError(condition.TimeframeInMS, "TimeframeInMS", "TimeframeInMS", "timeframerequired", "")
		}
	}
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

func (a *Action) Validate() error {
	validate := validator.New()
	validate.RegisterStructValidation(customActionValidation, Action{})
	err := validate.Struct(a)
	if err != nil {
		return err
	}
	return nil
}

// Helper for condition Validate()
func customActionValidation(sl validator.StructLevel) {
	action := sl.Current().Interface().(Action)
	if action.ValueType != "absolute" {
		if len(action.ValueQuoteMetric) == 0 {
			sl.ReportError(action.ValueQuoteMetric, "ValueQuoteMetric", "ValueQuoteMetric", "valuequotemetricrequired", "")
		}
	}
}

func (db *MGO) StrategyCreate(strategy *Strategy) error {
	c := db.DB(DatabaseName).C(StrategyCollection)
	err := c.Insert(strategy)
	return err
}

func (db *MGO) StrategyUpdate(strategy *Strategy) error {
	c := db.DB(DatabaseName).C(StrategyCollection)
	err := c.UpdateId(strategy.Id, strategy)
	return err
}

func (db *MGO) GetPausedStrategies() ([]Strategy, error) {
	var strategies []Strategy
	c := db.DB(DatabaseName).C(StrategyCollection)
	err := c.Find(bson.M{"status": "paused"}).All(&strategies)
	if err != nil {
		return nil, err
	}
	return strategies, nil
}
