package models

import (
	"gopkg.in/mgo.v2"
	// "log"
)

const DatabaseName = "onyxium"

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
