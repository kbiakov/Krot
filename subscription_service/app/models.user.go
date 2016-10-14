package main

type User struct {
	ID		string		`json:"_id" bson:"_id,omitempty"`
	Email		string		`json:"email" bson:"email"`
	Password	string		`json:"password" bson:"password"`
	Fullname	string		`json:"fullname" bson:"fullname"`
	Receivers	[]Receiver	`json:"receivers" bson:"receivers,inline"`
}

func GetUserByID(userID string) (*User, error) {
	user := User{}
	err := mongo.C("users").FindId(userID).One(&user)
	return &user, err
}
