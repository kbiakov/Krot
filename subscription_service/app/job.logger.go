package main

import "time"

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

func Log(id string, level rune, message string) {
	mongo.C("logs").Insert(&Log{
		ID:	 id,
		Level:	 level,
		Time:	 time.Now().Format(time.RFC3339),
		Message: message,
	})
}
