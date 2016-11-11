package main

import (
	"errors"
	"github.com/djimenez/iconv-go"
	"io"
	"net/http"
)

type Scraper interface {
	scrap() (string, error)
}

func CreateScraper(sType uint8, url string, tag string) (Scraper, error) {
	switch sType {

	case SubsType_HTML:
		return &HtmlScraper{
			Url:       url,
			ClassName: tag,
		}, nil
	case SubsType_JSON:
		return &JsonScraper{
			Url: url,
			Key: tag,
		}, nil
	case SubsType_XML:
		return &XmlScraper{
			Url: url,
			Tag: tag,
		}, nil
	}

	return nil, errors.New("Invalid scraper type")
}

func GetUtfContentByUrl(url string) (io.Reader, error) {
	// Get content by url
	// TODO: proxy request
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode content
	utfBody, err := iconv.NewReader(res.Body, "utf-8", "windows-1252")
	if err != nil {
		return nil, err
	}

	return utfBody, nil
}
