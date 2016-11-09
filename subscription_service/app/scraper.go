package main

import "errors"

type Scraper interface {
	scrap() (string, error)
}

func CreateScraper(sType uint8, url string, tag string) (Scraper, error) {
	switch sType {

	case SubsType_HTML:
		return &PageScraper{
			Url: url,
			ClassName: tag,
		}, nil
	case SubsType_JSON:
		return &RestScraper{
			Url: url,
			Key: tag,
		}, nil
	case SubsType_XML:
		return &SoapScraper{
			Url: url,
			Tag: tag,
		}, nil
	}

	return nil, errors.New("Invalid scraper type")
}
