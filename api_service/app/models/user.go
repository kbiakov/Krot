package krot

import (
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

var mongoUsers = mongo.C("users")

type User struct {
	ID		string		`json:"_id" bson:"_id,omitempty"`
	Email		string		`json:"email" bson:"email"`
	Password	string		`json:"password" bson:"password"`
	Fullname	string		`json:"fullname" bson:"fullname"`
	Receivers	[]Receiver	`json:"receivers" bson:"receivers,inline"`
}

// - DAO

func (u *User) CreateUser() error {
	hashedPassword, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = hashedPassword

	return mongoUsers.Insert(&u);
}

func (u User) IsValidPassword(hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, u.Password)
	return err == nil
}

func (u User) DeleteUser() error {
	return mongoUsers.RemoveId(u.ID)
}

func GetUser(email string) (*User, error) {
	user := User{}
	err := mongoUsers.Find(bson.M{"email": email}).One(&user)
	return &user, err
}

func GetUserByID(userID string) (*User, error) {
	user := User{}
	err := mongoUsers.FindId(userID).One(&user)
	return &user, err
}
