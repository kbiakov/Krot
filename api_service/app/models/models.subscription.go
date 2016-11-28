package dfdf

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/bamzi/jobrunner"
	"fmt"
)

const (
	SubsType_HTML = 0
	SubsType_JSON = 1
	SubsType_XML = 2
	SubsType_Unknown = -1
)

const (
	SubsStatus_Create = 0
	SubsStatus_Running = 1
	SubsStatus_Stopped = 2
	SubsStatus_Removed = -1
)

const ErrInvalidSubsStatus = error("Invalid subscription status")

type Subscription struct {
	ID	string	`json:"_id" bson:"_id,omitempty"`
	UserID	string	`json:"user_id" bson:"user_id"`
	Type	int8	`json:"type" bson:"type"`
	Url	string	`json:"url" bson:"url"`
	Tag	string	`json:"tag" bson:"tag"`
	PollMs	uint32	`json:"poll_ms" bson:"poll_ms"`
	Status	uint8	`json:"status" bson:"status"`
}

var mongoSubs = mongo.C("subscriptions")

// - DAO

func (s *Subscription) Subscribe() error {
	if s.Status != SubsStatus_Create {
		return ErrInvalidSubsStatus
	}

	if err := mongoSubs.Insert(&s); err != nil {
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

func (s *Subscription) UpdateSubscription() error {
	return s.CheckStatus(s.Status < SubsStatus_Running || s.Status > SubsStatus_Stopped, false, 0)
}

func (s Subscription) Unsubscribe() error {
	return s.CheckStatus(s.Status == SubsStatus_Removed, true, SubsStatus_Removed)
}

func GetSubscriptions(userID string) (*[]Subscription, error) {
	var subscriptions []Subscription
	err := mongoSubs.Find(bson.M{"user_id": userID}).All(&subscriptions)
	return &subscriptions, err
}

// - Helpers

func (s *Subscription) CheckStatus(condition bool, isChange bool, to uint8) error {
	if condition {
		return ErrInvalidSubsStatus
	}

	if isChange {
		s.Status = to
	}

	return mongoSubs.UpdateId(s.ID, &s)
}

func (s *Subscription) ChangeStatus(from uint8, to uint8) error {
	return s.CheckStatus(s.Status != from, true, to)
}

func (s *Subscription) StartSubscription() error {
	return s.ChangeStatus(SubsStatus_Create, SubsStatus_Running)
}

func (s Subscription) PollingInterval() string {
	return fmt.Sprintf("@every %d ms", s.PollMs)
}

// - Job processing

func (s *Subscription) Observe() {
	jobrunner.In(s.PollingInterval, &s)
}

func (s Subscription) Dismiss() {
	if err := mongoSubs.RemoveId(s.ID); err != nil {
		Log(s.ID, LogLevel_Error, err.Error())
	}
}

func (s *Subscription) Run() {
	status := s.Status

	if err := mongoSubs.FindId(s.ID).One(&s); err != nil {
		Log(s.ID, LogLevel_Error, err.Error())
	} else {
		if status != s.Status {
			Log(s.ID, LogLevel_Warning,
				"Inconsistent status: " +
				"expected " + status + ", " +
				"found "+ s.Status)
		}

		switch s.Status {

		case SubsStatus_Create:
			Log(s.ID, LogLevel_Info, "Observer running...")
			s.StartSubscription()
			s.Observe()

		case SubsStatus_Running:
			s.Lookup()
			s.Observe()

		case SubsStatus_Stopped:
			Log(s.ID, LogLevel_Info, "Observer stopped")

		case SubsStatus_Removed:
			Log(s.ID, LogLevel_Info, "Observer removed")
			s.Dismiss()
		}
	}
}

func (s Subscription) Lookup() {
	results := make(chan string)

	go func(results chan string) {
		job := Job{s, results}
		job.Process()
		close(results)
	}(results)

	awaitResults(results, func(result string, err error) {
		if err != nil {
			Log(s.ID, LogLevel_Error, err.Error())
		} else {
			NotifyUser(s.UserID, s.ID, result)
		}
	})
}
