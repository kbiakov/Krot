package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"

	service "../lib"
)

const DefaultChannel = "ch"

type Notification struct {
	SubscriptionID string   `json:"subscription_id"`
	Endpoints      []string `json:"endpoints"`
	Text           string   `json:"text"`
}

func main() {
	config := nsq.NewConfig()

	bindHandler("email", config, func(n *Notification) error {
		return service.SendEmail(n.SubscriptionID, n.Endpoints[0], n.Text)
	})

	bindHandler("apns", config, func(n *Notification) error {
		return service.SendApnsMessage(n.SubscriptionID, n.Endpoints[0], n.Text)
	})

	bindHandler("gcm", config, func(n *Notification) error {
		return service.SendGcmMessage(n.SubscriptionID, n.Endpoints, n.Text)
	})

	fmt.Scanln()
}

func bindHandler(topic string, config *nsq.Config, handler func(*Notification) error) {
	q, _ := nsq.NewConsumer(topic, DefaultChannel, config)

	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		notification := Notification{}
		if err := json.Unmarshal(message.Body, &notification); err != nil {
			return err
		}
		if err := handler(&notification); err != nil {
			return err
		}
		return nil
	}))

	if err := q.ConnectToNSQLookupd("lookupd:4161"); err != nil {
		panic(fmt.Sprintf("Cannot connect to NSQ server for %s-topic", topic))
	}
}
