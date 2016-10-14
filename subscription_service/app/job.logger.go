package main

import "time"

const (
	LogLevel_Info = iota
	LogLevel_Warning
	LogLevel_Error
)

type Log struct {
	SubscriptionID	string	`json:"subscription_id" bson:"subscription_id"`
	Level		rune	`json:"level" bson:"level"`
	Time		string	`json:"time" bson:"time"`
	Message		string	`json:"message" bson:"message"`
}

const logs = mongo.C("logs")

func Log(id string, level int, message string) {
	timeNow := time.Now().Format(time.RFC3339)
	logLevel := logLevelAsRune(level)

	logs.Insert(&Log{
		SubscriptionID: id,
		Level: logLevel,
		Time: timeNow,
		Message: message,
	})
}

func logLevelAsRune(level int) rune {
	switch level {

	case LogLevel_Info:
		return 'i'
	case LogLevel_Warning:
		return 'w'
	case LogLevel_Error:
		return 'e'
	default:
		return 'e'
	}
}
