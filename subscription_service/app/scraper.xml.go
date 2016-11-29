package main

import "errors"

type XmlScraper struct {
	Scraper

	Url string
	Tag string
}

// TODO: add implementation
func (s *XmlScraper) scrap() (string, error) {
	return "", errors.New("XML tag not found")
}
