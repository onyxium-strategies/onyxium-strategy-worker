package models

import (
	"fmt"
	"github.com/goware/emailx"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const UserCollection = "user"

func NewUser(email, password string) (*User, error) {
	err := emailx.Validate(email)
	if err != nil {
		return nil, err
	}
	hashPassword, err := HashAndSalt(password)
	if err != nil {
		return nil, err
	}
	user := &User{
		Id:          bson.NewObjectId(),
		Email:       email,
		Password:    hashPassword,
		IsActivated: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return user, nil
}

func (db *MGO) UserActivate(id string, token string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return fmt.Errorf("Incorrect id hex received: %s", id)
	}
	c := db.DB(DatabaseName).C(UserCollection)
	user := &User{}
	objectId := bson.ObjectIdHex(id)
	err := c.FindId(objectId).One(user)
	if err != nil {
		return fmt.Errorf("Error getting user with message: %s", err)
	}
	if user.IsActivated {
		return fmt.Errorf("User is already activated")
	}
	if ok, err := ComparePasswords(token, []byte(user.Email)); ok && err == nil {
		user.IsActivated = true
		user.ActivatedAt = time.Now()
		err = c.UpdateId(objectId, user)
		return err
	} else {
		return err
	}
}

func (db *MGO) UserAll() ([]User, error) {
	c := db.DB(DatabaseName).C(UserCollection)
	var users []User
	err := c.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (db *MGO) UserCreate(user *User) error {
	c := db.DB(DatabaseName).C(UserCollection)

	index := mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		return err
	}

	err = c.Insert(user)
	return err
}

func (db *MGO) UserGet(id string) (*User, error) {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return nil, fmt.Errorf("Incorrect IdHex received: %s", id)
	}
	c := db.DB(DatabaseName).C(UserCollection)
	user := &User{}
	objectId := bson.ObjectIdHex(id)
	err := c.FindId(objectId).One(user)
	if err != nil {
		return nil, fmt.Errorf("Error getting user with message: %s", err)
	}
	return user, nil
}

func (db *MGO) UserUpdate(user *User) error {
	c := db.DB(DatabaseName).C(UserCollection)
	user.UpdatedAt = time.Now()
	err := c.UpdateId(user.Id, user)
	return err
}

func (db *MGO) UserDelete(id string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return fmt.Errorf("Incorrect id hex received: %s", id)
	}
	c := db.DB(DatabaseName).C(UserCollection)
	objectId := bson.ObjectIdHex(id)
	err := c.RemoveId(objectId)
	return err
}

func HashAndSalt(pwd string) (string, error) {
	bytePwd := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, plainPwd []byte) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false, err
	}

	return true, nil
}
