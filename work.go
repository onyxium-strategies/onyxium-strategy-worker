package main

type Condition struct {
	ConditionType string `validate:"required"`
	BaseCurrency  string `validate:"required"`
	QuoteCurrency string `validate:"required"`
	TimeframeInMS int
	BaseMetric    string  `validate:"required"`
	Value         float64 `validate:"required"`
}

type Action struct {
	OrderType        string `validate:"required"`
	ValueType        string `validate:"required"`
	ValueQuoteMetric string
	BaseCurrency     string  `validate:"required"`
	QuoteCurrency    string  `validate:"required"`
	Quantity         float64 `validate:"required"`
	Value            float64 `validate:"required"`
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
