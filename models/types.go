package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type (
	DataStore interface {
		GetLatestMarket() (map[string]Market, error)
		GetHistoryMarket(TimeframeInMS int) (map[string]Market, error)
		UserActivate(id string, token string) error
		UserAll() ([]User, error)
		UserCreate(user *User) error
		UserGet(id string) (*User, error)
		UserUpdate(user *User) error
		UserDelete(id string) error
		StrategyCreate(strategy *Strategy) error
		StrategiesGetPaused() ([]Strategy, error)
		StrategyUpdate(strategy *Strategy) error
		OrderAll() ([]Order, error)
		OrderCreate(order *Order) error
		OrderGet(id string) (*Order, error)
		OrdersGetPending() ([]Order, error)
		OrderUpdate(order *Order) error
		OrderDelete(id string) error
	}

	MGO struct {
		*mgo.Session
	}

	Market struct {
		MarketName        string      `json:"MarketName" bson:"marketName"`
		High              float64     `json:"High" bson:"high"`
		Low               float64     `json:"Low" bson:"low"`
		Volume            float64     `json:"Volume" bson:"volume"`
		Last              float64     `json:"Last" bson:"last"`
		BaseVolume        float64     `json:"BaseVolume" bson:"baseVolume"`
		TimeStamp         string      `json:"TimeStamp" bson:"timeStamp"`
		Bid               float64     `json:"Bid" bson:"bid"`
		Ask               float64     `json:"Ask" bson:"ask"`
		OpenBuyOrders     int         `json:"OpenBuyOrders" bson:"openBuyOrders"`
		OpenSellOrders    int         `json:"OpenSellOrders" bson:"openSellOrders"`
		PrevDay           float64     `json:"PrevDay" bson:"prevDay"`
		Created           string      `json:"Created" bson:"created"`
		DisplayMarketName interface{} `json:"DisplayMarketName" bson:"displayMarketName"`
	}

	MarketRecord struct {
		Id     bson.ObjectId     `json:"id" bson:"_id,omitempty"`
		Market map[string]Market `bson:",inline"`
	}

	Strategy struct {
		Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Name     string        `json:"name" bson:"name"`
		JsonTree []interface{} `json:"jsonTree" bson:"jsonTree"`
		BsonTree *Tree         `json:"bsonTree" bson:"bsonTree"`
		Status   string        `json:"status" bson:"status"`
		State    int           `json:"state" bson:"state"`
		UserId   bson.ObjectId `json:"userId" bson:"userId"`
	}

	Tree struct {
		Id         int
		Left       *Tree
		Right      *Tree
		Conditions []Condition
		Action     Action
	}

	Condition struct {
		ConditionType string  `validate:"required,oneof=percentage-decrease percentage-increase greater-than-or-equal-to less-than-or-equal-to"`
		BaseCurrency  string  `validate:"required,nefield=QuoteCurrency"`
		QuoteCurrency string  `validate:"required",nefield=BaseCurrency`
		TimeframeInMS int     `validate:"omitempty,gt=0"`
		BaseMetric    string  `validate:"required,oneof=price-ask price-bid price-last volume"`
		Value         float64 `validate:"required,gte=0"`
	}

	Action struct {
		OrderType        string  `validate:"required,oneof=limit-buy limit-sell"`
		ValueType        string  `validate:"required,oneof=absolute relative-above relative-below percentage-above percentage-below"`
		ValueQuoteMetric string  `validate:"omitempty,oneof=price-ask price-bid price-last"`
		BaseCurrency     string  `validate:"required,nefield=QuoteCurrency"`
		QuoteCurrency    string  `validate:"required,nefield=BaseCurrency"`
		Quantity         float64 `validate:"required,gt=0"`
		Value            float64 `validate:"required"gt=0`
	}

	Order struct {
		Id            bson.ObjectId `bson:"_id,omitempty"`
		RemoteOrderId string        `bson:"remoteOrderId"`
		StrategyId    bson.ObjectId `bson:"strategyId"`
		NodeId        int           `bson:"nodeId"`
		Status        string        `bson:"status"`
		Rate          float64       `bson:"rate"`
		OrderType     string        `bson:"orderType"`
	}

	User struct {
		Id          bson.ObjectId     `json:"id" bson:"_id,omitempty"`
		Email       string            `json:"email" bson:"email"`
		Password    string            `json:"password" bson:"password"`
		IsActivated bool              `json:"isActivated" bson:"isActivated"`
		ActivatedAt time.Time         `json:"activatedAt" bson:"activatedAt"`
		LastLogin   time.Time         `json:"lastLogin" bson:"lastLogin"`
		CreatedAt   time.Time         `json:"createdAt" bson:"createdAt"`
		UpdatedAt   time.Time         `json:"updatedAt" bson:"updatedAt"`
		ApiKeys     map[string]string `json:"apiKeys" bson:"apiKeys"`
		// StrategyIds []int             `json:"strategyIds" bson:"strategyIds"`
	}
)
