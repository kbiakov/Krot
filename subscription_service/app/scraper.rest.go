package main

import (
	"github.com/djimenez/iconv-go"
	"net/http"
	"encoding/json"
	"errors"
	"bytes"
)

type RestScraper struct {
	Scraper

	Url string
	Key string
}

func (s *RestScraper) scrap() (string, error) {
	// Get content by url
	// TODO: proxy request?
	res, err := http.Get(s.Url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Decode content
	utfBodyReader, err := iconv.NewReader(res.Body, "utf-8", "windows-1252")
	if err != nil {
		return "", err
	}
	utfBodyBuf := new(bytes.Buffer)
	utfBodyBuf.ReadFrom(utfBodyReader)
	utfBody := utfBodyBuf.Bytes()

	var j interface{}
	if err := json.Unmarshal(utfBody, &j); err != nil {
		return "", err
	}

	if content, ok := searchJsonKey(j, s.Key); ok {
		return content.(string), nil
	}

	return "", errors.New("JSON key not found")
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
