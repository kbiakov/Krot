package main

import (
	"time"
	"errors"
)

const (
	ErrInvalidEmitterCast = errors.New("Invalid emitter cast")
	ErrJobTimeout = errors.New("Job exited by timeout")
)

type Job struct {
	Subscription *Subscription
	results chan<- string
}

func (j *Job) Process() {
	s := j.Subscription

	// TODO: switch by kind & create typed emitter
	emitter, _ := CreateScraper(s.Type, s.Url, s.Tag)

	if result, err := getResult(emitter); err != nil {
		Log(j.Subscription.ID, LogLevel_Error, err.Error())
	} else {
		j.results <- result
	}
}

func getResult(e interface{}) (string, error) {
	switch obj := e.(type) {

	case Scraper:
		return obj.scrap()
	default:
		return nil, ErrInvalidEmitterCast
	}
}

func awaitResults(results <-chan string, callback func(string, error)) {
	timeout := time.After(time.Minute)

	for {
		select { // non-blocking

		case result := <-results:
			callback(result, nil)
		case <-timeout:
			callback(nil, ErrJobTimeout)
			return
		default:
			return
		}
	}
}
