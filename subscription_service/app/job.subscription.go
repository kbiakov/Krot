package main

import (
	"github.com/bamzi/jobrunner"
	"fmt"
)

// - Job processing

func (s *Subscription) Schedule() {
	jobrunner.In(s.PollingInterval, &s)
}

func (s *Subscription) Run() {
	status := s.Status

	if err := data.FindId(s.ID).One(&s); err != nil {
		Log(s.ID, LogLevel_Error, err.Error())
		return
	}

	if s.Status != status {
		msg := fmt.Sprintf("expected %s, found %s", status, s.Status)
		Log(s.ID, LogLevel_Warning, "Inconsistent status: " + msg)
	}

	switch s.Status {

	case SubsStatus_Create:
		if err := s.StartSubscription(); err != nil {
			Log(s.ID, LogLevel_Error, "Observer created, but not runned: " + err.Error())
		} else {
			Log(s.ID, LogLevel_Info, "Observer created, running...")
		}

	case SubsStatus_Running:
		s.Execute()
		s.Schedule()

	case SubsStatus_Stopped:
		Log(s.ID, LogLevel_Info, "Observer stopped")

	case SubsStatus_Remove:
		if err := s.RemoveSubscription(); err != nil {
			Log(s.ID, LogLevel_Error, err.Error())
		} else {
			Log(s.ID, LogLevel_Info, "Observer removed")
		}
	}
}

func (s Subscription) Execute() {
	results := make(chan string)

	go func(results chan string) {
		job := Job{s, results}
		job.Process()
		close(results)
	}(results)

	awaitResults(results, func(res string, err error) {
		if err != nil {
			Log(s.ID, LogLevel_Error, err.Error())
		} else {
			NotifyUser(s.UserId, s.ID, res)
		}
	})
}
