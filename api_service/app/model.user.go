package main

import (
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID	  string     `json:"_id" bson:"_id,omitempty"`
	Email	  string     `json:"email" bson:"email"`
	Password  string     `json:"password" bson:"password"`
	Fullname  string     `json:"fullname" bson:"fullname"`
	Receivers []Receiver `json:"receivers" bson:"receivers,inline"`
}

// - Repository

var users = mongo.C("users")

func (u *User) CreateUser() error {
	hashedPassword, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = hashedPassword

	return users.Insert(&u);
}

func (u User) IsValidPassword(hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, u.Password) == nil
}

func (u User) DeleteUser() error {
	return users.RemoveId(u.ID)
}

func GetUser(email string) (*User, error) {
	u := User{}
	err := users.Find(bson.M{"email": email}).One(&u)
	return &u, err
}

func GetUserByID(userID string) (*User, error) {
	u := User{}
	err := users.FindId(userID).One(&u)
	return &u, err
}
