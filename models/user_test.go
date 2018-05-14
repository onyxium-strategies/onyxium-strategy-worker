package models_test

import (
	"bitbucket.org/visa-startups/coinflow-strategy-worker/models"
	// "gopkg.in/mgo.v2/bson"
	"testing"
)

func TestUser(t *testing.T) {
	db, err := models.InitDB("localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	user := &models.User{
		Email:    "test@email.com",
		Password: "pwddd",
	}

	// create user
	user, err = db.UserCreate(user)
	if err != nil {
		t.Fatal(err)
	}

	// get user
	hex := user.Id.Hex()
	user, err = db.UserGet(hex)
	if err != nil {
		t.Fatal(err)
	}

	// update user
	user.Email = "henk"
	_, err = db.UserUpdate(user)
	if err != nil {
		t.Fatal(err)
	}

	err = db.UserActivate("mijnid", hex)
	if err != nil {
		t.Fatal(err)
	}

	// remove user
	// err = db.UserDelete(hex)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	users, err := db.UserAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(users)
}
