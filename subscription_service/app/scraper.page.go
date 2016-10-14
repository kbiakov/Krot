package main

import (
	"net/http"
	"strings"
	"github.com/djimenez/iconv-go"
	"github.com/PuerkitoBio/goquery"
)

type PageScraper struct {
	Scraper

	Url string
	ClassName string
}

func (s *PageScraper) scrap() (string, error) {
	// Get content by url
	// TODO: proxy request?
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

	// Pack in document
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, err
	}

	content := extractContents(&doc, s.ClassName)

	return content, nil
}

func extractContents(doc *goquery.Document, className string) *[]string {
	content := []string{}

	doc.Find(className).Each(func(i int, s *goquery.Selection) {
		newContent := strings.TrimSpace(s.Text())
		content = append(content, newContent)
	})

	return &content
}
