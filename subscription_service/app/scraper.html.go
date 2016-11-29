package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type HtmlScraper struct {
	Scraper

	Url       string
	ClassName string
}

func (s *HtmlScraper) scrap() (string, error) {
	utfBody, err := GetUtfContentByUrl(s.Url)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return "", err
	}

	content := extractContent(doc, s.ClassName)

	return content, nil
}

func extractContent(doc *goquery.Document, className string) string {
	var buffer bytes.Buffer

	doc.Find(className).Each(func(i int, s *goquery.Selection) {
		newContent := strings.TrimSpace(s.Text())
		buffer.WriteString(newContent)
		buffer.WriteRune('\n')
	})

	return buffer.String()
}
