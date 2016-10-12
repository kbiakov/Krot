package main

import (
	"github.com/bitly/go-nsq"
	"encoding/json"
	"sync"
	"fmt"

	service "../lib"
)

const DefaultChannel = "ch"

type Notification struct {
	SubscriptionID	string	 `json:"subscription_id"`
	Endpoints	[]string `json:"endpoints"`
	Text		string	 `json:"text"`
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	config := nsq.NewConfig()

	bindHandler("email", &config, &wg, func(n *Notification) {
		service.SendEmail(n.SubscriptionID, n.Endpoints[0], n.Text)
	})

	bindHandler("apns", &config, &wg, func(n *Notification) {
		service.SendApnsMessage(n.SubscriptionID, n.Endpoints[0], n.Text)
	})

	bindHandler("gcm", &config, &wg, func(n *Notification) {
		service.SendGcmMessage(n.SubscriptionID, n.Endpoints, n.Text)
	})

	wg.Wait()
}

func bindHandler(topic string, config *nsq.Config, wg *sync.WaitGroup, handler func(*Notification))  {
	q, _ := nsq.NewConsumer(topic, DefaultChannel, &config)

	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		notification := Notification{}
		if err := json.Unmarshal(message.Body, &notification); err != nil {
			return err
		}

		handler(&notification)

		wg.Done()

		return nil
	}))

	if err := q.ConnectToNSQLookupd("lookupd:4161"); err != nil {
		panic(fmt.Sprintf("Cannot connect to NSQ server for %s-topic", topic))
	}
}
