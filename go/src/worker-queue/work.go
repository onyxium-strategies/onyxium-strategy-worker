package main

type Condition struct {
	ConditionType string
	BaseCurrency  string
	QuoteCurrency string
	TimeframeInMS int
	BaseMetric    string
	Value         float64
}

type Action struct {
	Type           string
	OrderType      string
	OrderValueType string
	BaseCurrency   string
	QuoteCurrency  string
	Quantity       float32
	Value          float64
}

type Tree struct {
	Left       *Tree
	Right      *Tree
	Conditions []Condition
	Action     Action
}

type WorkRequest struct {
	ID   int
	Tree *Tree
}
