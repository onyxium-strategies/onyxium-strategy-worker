package models

import (
	"fmt"
	// "gopkg.in/mgo.v2/bson"
)

const OrderCollection = "order"

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

// CRUD IS NIET MEER NODIG ALS ORDER EEN ONDERDEEL IS VAN STRATEGY
// func (db *MGO) OrderAll() ([]Order, error) {
// 	c := db.DB(DatabaseName).C(OrderCollection)
// 	var orders []Order
// 	err := c.Find(bson.M{}).All(&orders)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }

// func (db *MGO) OrderCreate(order *Order) error {
// 	c := db.DB(DatabaseName).C(OrderCollection)

// 	err := c.Insert(order)
// 	return err
// }

// func (db *MGO) OrderGet(id string) (*Order, error) {
// 	ok := bson.IsObjectIdHex(id)
// 	if !ok {
// 		return nil, fmt.Errorf("Incorrect IdHex received: %s", id)
// 	}
// 	c := db.DB(DatabaseName).C(OrderCollection)
// 	order := &Order{}
// 	objectId := bson.ObjectIdHex(id)
// 	err := c.FindId(objectId).One(order)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error getting order with message: %s", err)
// 	}
// 	return order, nil
// }

// func (db *MGO) OrderUpdate(order *Order) error {
// 	c := db.DB(DatabaseName).C(OrderCollection)
// 	err := c.UpdateId(order.Id, order)
// 	return err
// }

// func (db *MGO) OrderDelete(id string) error {
// 	ok := bson.IsObjectIdHex(id)
// 	if !ok {
// 		return fmt.Errorf("Incorrect id hex received: %s", id)
// 	}
// 	c := db.DB(DatabaseName).C(OrderCollection)
// 	objectId := bson.ObjectIdHex(id)
// 	err := c.RemoveId(objectId)
// 	return err
// }

// func (db *MGO) OrdersGetPending() ([]Order, error) {
// 	var orders []Order
// 	c := db.DB(DatabaseName).C(OrderCollection)
// 	err := c.Find(bson.M{"status": "pending"}).All(&orders)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }
