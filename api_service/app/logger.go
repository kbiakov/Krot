package main

import "gopkg.in/mgo.v2/bson"

const (
	LogLevel_Info    = 'i'
	LogLevel_Warning = 'w'
	LogLevel_Error   = 'e'
)

type Log struct {
	ID	string `json:"subscription_id" bson:"subscription_id"`
	Level	rune   `json:"level" bson:"level"`
	Time	string `json:"time" bson:"time"`
	Message	string `json:"message" bson:"message"`
}

// TODO: add paging, etc?
func GetLogsForUserID(userID string) (*[]Log, error) {
	ss, err := GetSubscriptionsForUserID(userID)
	if err != nil {
		return nil, err
	}

	var logs []Log
	for _, s := range *ss {
		var tmp []Log
		mongo.C("logs").
			Find(bson.M{"subscription_id": s.ID}).
			Sort([]string{"-time"}).
			All(&tmp)
		logs = append(logs, tmp)
	}

	return logs, err
}
