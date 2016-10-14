package main

import (
	"github.com/bamzi/jobrunner"
	"fmt"
)

// - Job processing

func (s *Subscription) Observe() {
	jobrunner.In(s.PollingInterval, &s)
}

func (s Subscription) Dismiss() {
	if err := data.RemoveId(s.ID); err != nil {
		Log(s.ID, LogLevel_Error, err.Error())
	}
}

func (s *Subscription) Run() {
	status := s.Status

	if err := data.FindId(s.ID).One(&s); err != nil {
		Log(s.ID, LogLevel_Error, err.Error())
	} else {
		if status != s.Status {
			msg := fmt.Sprintf("expected %s, found %s", status, s.Status)
			Log(s.ID, LogLevel_Warning, "Inconsistent status: " + msg)
		}

		switch s.Status {

		case SubsStatus_Created:
			Log(s.ID, LogLevel_Info, "Observer created, running...")
			s.StartSubscription()
			s.Observe()

		case SubsStatus_Running:
			s.Lookup()
			s.Observe()

		case SubsStatus_Stopped:
			Log(s.ID, LogLevel_Info, "Observer stopped")

		case SubsStatus_NotExist:
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
			NotifyUser(s.UserId, s.ID, result)
		}
	})
}
