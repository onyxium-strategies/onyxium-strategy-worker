package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

const OrderCollection = "order"

func NewOrder(orderType, remoteOrderId string, strategyId bson.ObjectId, nodeId int, rate float64) (*Order, error) {
	if orderType == "" || remoteOrderId == "" || strategyId == "" || rate == 0 {
		return nil, fmt.Errorf("OrderType, RemoteOrderId, StrategyId and Rate are required to create an Order")
	}
	order := &Order{
		RemoteOrderId: remoteOrderId,
		StrategyId:    strategyId,
		NodeId:        nodeId,
		Status:        "pending",
		Rate:          rate,
		OrderType:     orderType,
	}
	return order, nil
}

func (db *MGO) OrderAll() ([]Order, error) {
	c := db.DB(DatabaseName).C(OrderCollection)
	var orders []Order
	err := c.Find(bson.M{}).All(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (db *MGO) OrderCreate(order *Order) error {
	c := db.DB(DatabaseName).C(OrderCollection)

	err := c.Insert(order)
	return err
}

func (db *MGO) OrderGet(id string) (*Order, error) {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return nil, fmt.Errorf("Incorrect IdHex received: %s", id)
	}
	c := db.DB(DatabaseName).C(OrderCollection)
	order := &Order{}
	objectId := bson.ObjectIdHex(id)
	err := c.FindId(objectId).One(order)
	if err != nil {
		return nil, fmt.Errorf("Error getting order with message: %s", err)
	}
	return order, nil
}

func (db *MGO) OrderUpdate(order *Order) error {
	c := db.DB(DatabaseName).C(OrderCollection)
	err := c.UpdateId(order.Id, order)
	return err
}

func (db *MGO) OrderDelete(id string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return fmt.Errorf("Incorrect id hex received: %s", id)
	}
	c := db.DB(DatabaseName).C(OrderCollection)
	objectId := bson.ObjectIdHex(id)
	err := c.RemoveId(objectId)
	return err
}
