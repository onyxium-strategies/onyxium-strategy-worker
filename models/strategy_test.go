package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Also tests Tree.Search
func TestSetIdsForBinarySearch(t *testing.T) {
	testCases := []struct {
		name              string
		t                 *Tree
		expectedInOrder   []int
		expectedPreOrder  []int
		expectedPostOrder []int
	}{
		{
			name: "Tree 1",
			t: &Tree{
				Conditions: []Condition{
					{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.072},
					{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.066},
				},
				Action: Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 1, Value: 0.08},
				Left:   nil,
				Right: &Tree{
					Conditions: []Condition{
						{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.075},
					},
					Action: Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 2, Value: 0.08},
					Left:   nil,
					Right: &Tree{
						Conditions: []Condition{
							{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "volume", Value: 8000},
						},
						Action: Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 3, Value: 0.08},
						Left: &Tree{
							Conditions: []Condition{
								{ConditionType: "greater-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.065},
							},
							Action: Action{OrderType: "limit-buy", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 4, Value: 0.08},
							Left:   nil,
							Right:  nil,
						},
						Right: nil,
					},
				},
			},
			expectedInOrder:   []int{0, 1, 2, 3},
			expectedPreOrder:  []int{0, 1, 3, 2},
			expectedPostOrder: []int{2, 3, 1, 0},
		},
		{
			name: "Tree 2",
			t: &Tree{
				Conditions: []Condition{
					{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "volume", Value: 9000},
				},
				Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 1, Value: 0.088},
				Left:   nil,
				Right: &Tree{
					Conditions: []Condition{
						{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.069},
						{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.07},
					},
					Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 2, Value: 0.088},
					Left:   nil,
					Right: &Tree{
						Conditions: []Condition{
							{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.071},
						},
						Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 3, Value: 0.088},
						Left: &Tree{
							Conditions: []Condition{
								{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.072},
							},
							Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 4, Value: 0.088},
							Left: &Tree{
								Conditions: []Condition{
									{ConditionType: "less-than-or-equal-to", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.075},
								},
								Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 5, Value: 0.088},
								Left:   nil,
								Right:  nil,
							},
							Right: nil,
						},
						Right: nil,
					},
				},
			},
			expectedInOrder:   []int{0, 1, 2, 3, 4},
			expectedPreOrder:  []int{0, 1, 4, 3, 2},
			expectedPostOrder: []int{2, 3, 4, 1, 0},
		},
		{
			name: "Tree 3",
			t: &Tree{
				Conditions: []Condition{
					{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.1, TimeframeInMS: 1},
				},
				Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 1, Value: 0.079},
				Left:   nil,
				Right: &Tree{
					Conditions: []Condition{
						{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.07, TimeframeInMS: 1},
					},
					Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 2, Value: 0.079},
					Left: &Tree{
						Conditions: []Condition{
							{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.065, TimeframeInMS: 1},
							{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "volume", Value: 0.3, TimeframeInMS: 1},
						},
						Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 3, Value: 0.079},
						Left:   nil,
						Right:  nil,
					},
					Right: &Tree{
						Conditions: []Condition{
							{ConditionType: "percentage-increase", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.07, TimeframeInMS: 1},
						},
						Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 4, Value: 0.079},
						Left:   nil,
						Right:  nil,
					},
				},
			},
			expectedInOrder:   []int{0, 1, 2, 3},
			expectedPreOrder:  []int{0, 2, 1, 3},
			expectedPostOrder: []int{1, 3, 2, 0},
		},
		{
			name: "Tree 4",
			t: &Tree{
				Conditions: []Condition{
					{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-last", Value: 0.045, TimeframeInMS: 0},
				},
				Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 1, Value: 0.079},
				Left: &Tree{
					Conditions: []Condition{
						{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-bid", Value: 0.065, TimeframeInMS: 0},
					},
					Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 2, Value: 0.079},
					Left: &Tree{
						Conditions: []Condition{
							{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "price-ask", Value: 0.058, TimeframeInMS: 0},
						},
						Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 3, Value: 0.079},
						Left: &Tree{
							Conditions: []Condition{
								{ConditionType: "percentage-decrease", BaseCurrency: "BTC", QuoteCurrency: "ETH", BaseMetric: "volume", Value: 0.1, TimeframeInMS: 0},
							},
							Action: Action{OrderType: "limit-sell", ValueType: "absolute", BaseCurrency: "BTC", QuoteCurrency: "ETH", Quantity: 4, Value: 0.079},
							Left:   nil,
							Right:  nil,
						},
						Right: nil,
					},
					Right: nil,
				},
				Right: nil,
			},
			expectedInOrder:   []int{0, 1, 2, 3},
			expectedPreOrder:  []int{3, 2, 1, 0},
			expectedPostOrder: []int{0, 1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.t.SetIdsForBinarySearch()
			tc.t.String()
			n := countNodes(tc.t)

			// Test if we can find all nodes
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			for _, id := range s {
				_, err := tc.t.Search(id)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Test in order traverse
			ch := make(chan int, n)
			inOrderTraverseWriteChan(tc.t, ch)
			close(ch)
			s = make([]int, 0)
			for i := range ch {
				s = append(s, i)
			}
			assert.Equal(t, tc.expectedInOrder, s)

			// Test pre order travers
			ch = make(chan int, n)
			preOrderTraverseWriteChan(tc.t, ch)
			close(ch)
			s = make([]int, 0)
			for i := range ch {
				s = append(s, i)
			}
			assert.Equal(t, tc.expectedPreOrder, s)

			// Test post order travers
			ch = make(chan int, n)
			postOrderTraverseWriteChan(tc.t, ch)
			close(ch)
			s = make([]int, 0)
			for i := range ch {
				s = append(s, i)
			}
			assert.Equal(t, tc.expectedPostOrder, s)
		})
	}
}

func inOrderTraverseWriteChan(t *Tree, ch chan int) {
	if t != nil {
		inOrderTraverseWriteChan(t.Left, ch)
		ch <- t.Id
		inOrderTraverseWriteChan(t.Right, ch)
	}
}

func preOrderTraverseWriteChan(t *Tree, ch chan int) {
	if t != nil {
		ch <- t.Id
		preOrderTraverseWriteChan(t.Left, ch)
		preOrderTraverseWriteChan(t.Right, ch)
	}
}

func postOrderTraverseWriteChan(t *Tree, ch chan int) {
	if t != nil {
		postOrderTraverseWriteChan(t.Left, ch)
		postOrderTraverseWriteChan(t.Right, ch)
		ch <- t.Id
	}
}
