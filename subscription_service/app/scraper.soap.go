package main

import "errors"

type SoapScraper struct {
	Scraper

	Url string
	Tag string
}

// TODO: add implementation
func (s *SoapScraper) scrap() (string, error) {
	return "", errors.New("XML tag not found")
}
