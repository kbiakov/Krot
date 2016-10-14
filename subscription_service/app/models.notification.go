package main

import (
	"fmt"
	"encoding/json"
)

type Notification struct {
	SubscriptionID	string		`json:"subscription_id"`
	Endpoints	*[]string	`json:"endpoints"`
	Text		string		`json:"text"`
}

func CreateNotification(subsID string, to string) *Notification {
	return &Notification{subsID, []string{to}}
}

func PushMessageToQueue(topic string, n *Notification)  {
	jsonStr, err := json.Marshal(&n)
	if err != nil {
		msg := fmt.Sprintf("Error marshalling notification: %s", err.Error())
		Log(n.SubscriptionID, LogLevel_Error, msg)
		return
	}

	if err = w.Publish(topic, []byte(jsonStr)); err != nil {
		msg := fmt.Sprintf("Error publishing to NSQ topic %s: %s", topic, err.Error())
		Log(n.SubscriptionID, LogLevel_Error, msg)
	}
}
