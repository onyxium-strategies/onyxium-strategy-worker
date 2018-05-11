package models

import (
	"gopkg.in/mgo.v2"
	// "log"
)

type DataStore interface {
	GetLatestMarket() (map[string]Market, error)
	GetHistoryMarket(TimeframeInMS int) (map[string]Market, error)
	UserActivate(id string) error
	UserAll() ([]User, error)
	UserCreate(user *User) (*User, error)
	UserGet(id string) (*User, error)
	UserUpdate(user *User) error
	UserDelete(id string) error
}

type MGO struct {
	*mgo.Session
}

// http://www.alexedwards.net/blog/organising-database-access
// https://hackernoon.com/how-to-work-with-databases-in-golang-33b002aa8c47
func InitDB(dataSourceName string) (*MGO, error) {
	var err error
	DBCon, err := mgo.Dial(dataSourceName)
	if err != nil {
		return nil, err
	}
	DBCon.SetMode(mgo.Monotonic, true)

	if err = DBCon.Ping(); err != nil {
		return nil, err
	}
	return &MGO{DBCon}, nil
}
