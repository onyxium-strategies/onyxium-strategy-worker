package main

type Condition struct {
  ConditionType string
  BaseCurrency string
  QuoteCurrency string
  TimeframeInMS int
  BaseMetric string
  Value float64
}
// example condition json:
// {
//  "conditionType":"percentage-increase",
//  "baseCurrency":"ETH",
//  "quoteCurrency":"OMG",
//  "timeframeInMS":3600000,
//  "baseMetric":"price",
//  "value":20
// }

type Action struct {
  Type0 string
  OrderType string
  OrderValueType string
  BaseCurrency string
  QuoteCurrency string
  Quantity float32
  Value float64
}
// example action json:
// "action":{
//  "type":"order",
//  "orderType":"limit-buy",
//  "orderValueType":"absolute",
//  "baseCurrency":"ETH",
//  "quoteCurrency":"OMG",
//  "quantity":100,
//  "value":0.012
// }


type Tree struct {
  // Parent *Tree
  Left  *Tree
  Right *Tree
  Conditions []Condition
  Action Action
}

type WorkRequest struct {
  Id int
  Tree *Tree
}

