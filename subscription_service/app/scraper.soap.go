package main

type SoapScraper struct {
	Scraper

	Url string
	Tag string
}

const ErrXmlKeyNotFound = error("XML tag not found")

func (s *SoapScraper) scrap() (string, error) {
	return nil, ErrXmlKeyNotFound
}
