package main

import (
	"errors"
	"time"
)

const (
	SubsType_HTML = iota
	SubsType_JSON
	SubsType_XML
)

const (
	SubsStatus_Create = iota
	SubsStatus_Running
	SubsStatus_Stopped
	SubsStatus_Remove
)

type Subscription struct {
	ID     string `json:"_id" bson:"_id,omitempty"`
	UserId string `json:"user_id" bson:"user_id"`
	Type   uint8  `json:"type" bson:"type"`
	Url    string `json:"url" bson:"url"`
	Tag    string `json:"tag" bson:"tag"`
	PollMs uint32 `json:"poll_ms" bson:"poll_ms"`
	Status uint8  `json:"status" bson:"status"`
}

// - Repository

var data = mongo.C("subscriptions")

// - API

func (s *Subscription) Subscribe() error {
	if s.Status != SubsStatus_Create {
		return NewInvalidStatusError()
	}

	if err := data.Insert(&s); err != nil {
		return err
	}

	s.Schedule()

	return nil
}

func (s *Subscription) ResumeSubscription() error {
	return s.ChangeStatus(SubsStatus_Stopped, SubsStatus_Running)
}

func (s *Subscription) StopSubscription() error {
	return s.ChangeStatus(SubsStatus_Running, SubsStatus_Stopped)
}

func (s Subscription) Unsubscribe() error {
	if s.Status != SubsStatus_Running {
		return s.RemoveSubscription()
	}
	return s.ChangeStatus(SubsStatus_Running, SubsStatus_Remove)
}

// - Internal API

func GetSubscription(ID string) (*Subscription, error) {
	s := Subscription{}
	err := data.FindId(ID).One(&s)
	return &s, err
}

func (s *Subscription) StartSubscription() error {
	if err := s.ChangeStatus(SubsStatus_Create, SubsStatus_Running); err != nil {
		return err
	}

	s.Schedule()

	return nil
}

func (s Subscription) RemoveSubscription() error {
	return data.RemoveId(s.ID)
}

// - Helpers

func (s *Subscription) ChangeStatus(from uint8, to uint8) error {
	if s.Status != from {
		return NewInvalidStatusError()
	}

	s.Status = to

	return data.UpdateId(s.ID, &s)
}

func StatusToString(status uint8) string {
	switch status {
	case SubsStatus_Create:
		return "create"
	case SubsStatus_Running:
		return "running"
	case SubsStatus_Stopped:
		return "stopped"
	case SubsStatus_Remove:
		return "remove"
	default:
		return "undefined"
	}
}

func (s Subscription) PollingInterval() time.Duration {
	return time.Duration(s.PollMs) * time.Millisecond
}

func NewInvalidStatusError() error {
	return errors.New("Invalid subscription status")
}
