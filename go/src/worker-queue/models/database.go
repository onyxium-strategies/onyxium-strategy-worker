package models

import (
	"gopkg.in/mgo.v2"
	"log"
)

var (
	DBCon *mgo.Session
)

// https://stackoverflow.com/questions/31218008/sharing-a-globally-defined-db-conn-with-multiple-packages-in-golang
// http://www.alexedwards.net/blog/organising-database-access
// https://hackernoon.com/how-to-work-with-databases-in-golang-33b002aa8c47
func InitDB(dataSourceName string) {
	var err error
	DBCon, err = mgo.Dial(dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	// defer DBCon.Close()
	DBCon.SetMode(mgo.Monotonic, true)

	if err = DBCon.Ping(); err != nil {
		log.Fatal(err)
	}
}
