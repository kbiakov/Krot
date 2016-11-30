package main

import "gopkg.in/mgo.v2/bson"

type Subscription struct {
	ID     string `json:"_id" bson:"_id,omitempty"`
	UserId string `json:"user_id" bson:"user_id"`
	Type   uint8  `json:"type" bson:"type"`
	Url    string `json:"url" bson:"url"`
	Tag    string `json:"tag" bson:"tag"`
	PollMs uint32 `json:"poll_ms" bson:"poll_ms"`
	Status uint8  `json:"status" bson:"status"`
}

func GetSubscriptionsForUserID(userID string) (*[]Subscription, error) {
	var ss []Subscription
	query := bson.M{"user_id": userID}

	if err := mongo.C("subscriptions").Find(query).All(&ss); err != nil {
		return nil, err
	}

	return ss, nil
}
