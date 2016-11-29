package main

import (
	"encoding/json"
	"fmt"
)

type Notification struct {
	SubscriptionID string    `json:"subscription_id"`
	Endpoints      *[]string `json:"endpoints"`
	Text           string    `json:"text"`
}

func CreateNotification(subsID string, to string, text string) *Notification {
	return &Notification{subsID, &[]string{to}, text}
}

func PushToQueue(topic string, n *Notification) {
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
