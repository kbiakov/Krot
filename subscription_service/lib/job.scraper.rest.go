package krot

import (
	"github.com/djimenez/iconv-go"
	"net/http"
	"encoding/json"
)

type RestScraper struct {
	Scraper

	Url string
	Key string
}

const ErrJsonKeyNotFound = error("JSON key not found.")

func (s *RestScraper) scrap() (string, error) {
	// Get content by url
	res, err := http.Get(s.Url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Decode content
	utfBody, err := iconv.NewReader(res.Body, "utf-8", "windows-1252")
	if err != nil {
		return nil, err
	}

	var j interface{}
	err = json.Unmarshal(utfBody, &j)
	if err != nil {
		return nil, err
	}

	if content, ok := searchJsonKey(j, s.Key); ok {
		return content, nil
	}

	return nil, ErrJsonKeyNotFound
}

func searchJsonKey(obj interface{}, key string) (interface{}, bool) {
	switch t := obj.(type) {

	case map[string]interface{}:
		if v, ok := t[key]; ok {
			return v, ok
		}
		for _, v := range t {
			if result, ok := searchJsonKey(v, key); ok {
				return result, ok
			}
		}

	case []interface{}:
		for _, v := range t {
			if result, ok := searchJsonKey(v, key); ok {
				return result, ok
			}
		}
	}

	// Key not found
	return nil, false
}
