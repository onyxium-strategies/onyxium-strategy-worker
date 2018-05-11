package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id          bson.ObjectId     `json:"id" bson:"_id,omitempty"`
	Email       string            `json:"email" bson:"email"`
	Password    string            `json:"password" bson:"password"`
	IsActivated bool              `bson:"isActivated"`
	ActivatedAt time.Time         `bson:"activatedAt"`
	LastLogin   time.Time         `bson:"lastLogin"`
	CreatedAt   time.Time         `bson:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt"`
	ApiKeys     map[string]string `bson:"apiKeys"`
	StrategyIds []int             `bson:"strategyIds"`
}

func (db *MGO) UserActivate(id string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return fmt.Errorf("Incorrect IdHex received: %s", id)
	}
	c := db.DB("coinflow").C("user")
	user := &User{}
	objectId := bson.ObjectIdHex(id)
	err := c.FindId(objectId).One(user)
	if err != nil {
		return fmt.Errorf("Error getting user with message: %s", err)
	}
	user.IsActivated = true
	user.ActivatedAt = time.Now()
	err = c.UpdateId(objectId, user)
	return err
}

func (db *MGO) UserAll() ([]User, error) {
	c := db.DB("coinflow").C("user")
	var users []User
	err := c.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *MGO) UserCreate(user *User) (*User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Id = bson.NewObjectId()
	pwd, err := hashAndSalt(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = pwd
	c := db.DB("coinflow").C("user")
	err = c.Insert(user)
	if err != nil {
		return nil, fmt.Errorf("Error creating user with message: %s", err)
	}
	return user, nil
}

func (db *MGO) UserGet(id string) (*User, error) {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return nil, fmt.Errorf("Incorrect IdHex received: %s", id)
	}
	c := db.DB("coinflow").C("user")
	user := &User{}
	objectId := bson.ObjectIdHex(id)
	err := c.FindId(objectId).One(user)
	if err != nil {
		return nil, fmt.Errorf("Error getting user with message: %s", err)
	}
	return user, nil
}

func (db *MGO) UserUpdate(user *User) error {
	user.UpdatedAt = time.Now()
	c := db.DB("coinflow").C("user")
	err := c.UpdateId(user.Id, user)
	return err
}

func (db *MGO) UserDelete(id string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return fmt.Errorf("Incorrect IdHex received: %s", id)
	}
	c := db.DB("coinflow").C("user")
	objectId := bson.ObjectIdHex(id)
	err := c.RemoveId(objectId)
	return err
}

func hashAndSalt(pwd string) (string, error) {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd []byte) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false, err
	}

	return true, nil
}
