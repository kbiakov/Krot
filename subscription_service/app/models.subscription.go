package main

import (
	"fmt"
)

const (
	SubsType_HTML = iota
	SubsType_JSON
	SubsType_XML
)

const (
	SubsStatus_NotExist = iota
	SubsStatus_Created
	SubsStatus_Running
	SubsStatus_Stopped
)

const ErrInvalidSubsStatus = error("Invalid subscription status")

type Subscription struct {
	ID	string	`json:"_id" bson:"_id,omitempty"`
	UserId	string	`json:"user_id" bson:"user_id"`
	Type	int8	`json:"type" bson:"type"`
	Url	string	`json:"url" bson:"url"`
	Tag	string	`json:"tag" bson:"tag"`
	PollMs	uint32	`json:"poll_ms" bson:"poll_ms"`
	Status	uint8	`json:"status" bson:"status"`
}

const data = mongo.C("subscriptions")

// - Repository

func GetSubscription(ID string) (*Subscription, error) {
	s := Subscription{}
	err := data.FindId(ID).One(&s)
	return &s, err
}

func (s *Subscription) Subscribe() error {
	if s.Status != SubsStatus_Created {
		return ErrInvalidSubsStatus
	}

	if err := data.Insert(&s); err != nil {
		return err
	}

	s.Observe()

	return nil
}

func (s *Subscription) ResumeSubscription() error {
	return s.ChangeStatus(SubsStatus_Stopped, SubsStatus_Running)
}

func (s *Subscription) StopSubscription() error {
	return s.ChangeStatus(SubsStatus_Running, SubsStatus_Stopped)
}

func (s Subscription) Unsubscribe() error {
	return s.CheckStatus(s.Status == SubsStatus_NotExist, true, SubsStatus_NotExist)
}

// - Helpers

func (s *Subscription) CheckStatus(condition bool, isChange bool, to uint8) error {
	if condition {
		return ErrInvalidSubsStatus
	}

	if isChange {
		s.Status = to
	}

	return data.UpdateId(s.ID, &s)
}

func (s *Subscription) ChangeStatus(from uint8, to uint8) error {
	return s.CheckStatus(s.Status != from, true, to)
}

func (s *Subscription) StartSubscription() error {
	return s.ChangeStatus(SubsStatus_Created, SubsStatus_Running)
}

func (s Subscription) PollingInterval() string {
	return fmt.Sprintf("@every %d ms", s.PollMs)
}
