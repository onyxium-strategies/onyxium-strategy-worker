package models

import (
	"fmt"
	// "github.com/tidwall/gjson"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const StrategyCollection = "strategy"

var validate *validator.Validate

func NewStrategy(name, userId string, bsonTree *Tree) (*Strategy, error) {
	// if !gjson.Valid(jsonTree) {
	// 	return nil, fmt.Errorf("Invalid json structure received: %s", jsonTree)
	// }
	strategy := &Strategy{
		Id:        bson.NewObjectId(),
		Name:      name,
		BsonTree:  bsonTree,
		Status:    "paused",
		State:     bsonTree.Id,
		UserId:    bson.ObjectIdHex(userId),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return strategy, nil
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

func (t *Tree) ToKaryArray() ([]map[string]interface{}, error) {
	root, err := _toKaryArry(t)
	return root, err
}

func _toKaryArry(t *Tree) ([]map[string]interface{}, error) {
	if t == nil {
		return make([]map[string]interface{}, 0), nil
	}
	var node []map[string]interface{}
	if t.Left != nil {
		children, _ := _toKaryArry(t.Left)
		node = []map[string]interface{}{{
			"action":     t.Action,
			"conditions": t.Conditions,
			"id":         t.Id,
			"order":      t.Order,
			"then":       children,
		}}
	} else {
		node = []map[string]interface{}{{
			"action":     t.Action,
			"conditions": t.Conditions,
			"id":         t.Id,
			"order":      t.Order,
		}}
	}
	for t.Right != nil {
		node = append(node, _appendChild(t.Right))
		t = t.Right
	}
	return node, nil
}

func _appendChild(t *Tree) map[string]interface{} {
	if t == nil {
		return make(map[string]interface{})
	}
	if t.Left != nil {
		children, _ := _toKaryArry(t.Left)
		return map[string]interface{}{
			"action":     t.Action,
			"conditions": t.Conditions,
			"id":         t.Id,
			"order":      t.Order,
			"then":       children,
		}
	} else {
		return map[string]interface{}{
			"action":     t.Action,
			"conditions": t.Conditions,
			"id":         t.Id,
			"order":      t.Order,
		}
	}
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
	strategy.UpdatedAt = time.Now()
	err := c.UpdateId(strategy.Id, strategy)
	return err
}

// TODO: add user argument so you can retrieve the strategies specific to a user
//
func (db *MGO) StrategiesGetPaused() ([]Strategy, error) {
	var strategies []Strategy
	c := db.DB(DatabaseName).C(StrategyCollection)
	err := c.Find(bson.M{"status": "paused"}).All(&strategies)
	if err != nil {
		return nil, err
	}
	return strategies, nil
}

func (db *MGO) StrategyAll() ([]Strategy, error) {
	var strategies []Strategy
	c := db.DB(DatabaseName).C(StrategyCollection)
	err := c.Find(bson.M{}).All(&strategies)
	if err != nil {
		return nil, err
	}
	return strategies, nil
}

func (db *MGO) StrategyGet(id string) (*Strategy, error) {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return nil, fmt.Errorf("Incorrect IdHex received: %s", id)
	}
	c := db.DB(DatabaseName).C(StrategyCollection)
	strategy := &Strategy{}
	objectId := bson.ObjectIdHex(id)
	err := c.FindId(objectId).One(strategy)
	if err != nil {
		return nil, fmt.Errorf("Error getting strategy with message: %s", err)
	}
	return strategy, nil
}

func (db *MGO) StrategyDelete(id string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return fmt.Errorf("Incorrect id hex received: %s", id)
	}
	c := db.DB(DatabaseName).C(StrategyCollection)
	objectId := bson.ObjectIdHex(id)
	err := c.RemoveId(objectId)
	return err
}

func NewOrder(remoteOrderId string, rate float64) (*Order, error) {
	if remoteOrderId == "" || rate == 0 {
		return nil, fmt.Errorf("Action, RemoteOrderId, StrategyId and Rate are required to create an Order")
	}
	order := &Order{
		RemoteOrderId: remoteOrderId,
		Status:        "pending",
		Rate:          rate,
	}
	return order, nil
}
