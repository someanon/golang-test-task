package main

import (
	"bytes"
	"errors"
	"io"

	"golang.org/x/net/html"
)

type page struct {
	URL      string    `json:"url"`
	Meta     meta      `json:"meta"`
	Elements []element `json:"elements,omitempty"`
}

type meta struct {
	Status        int    `json:"status"`
	ContentType   string `json:"content-type,omitempty"`
	ContentLength int    `json:"content-length,omitempty"`
}

type element struct {
	TagName string `json:"tag-name"`
	Count   int    `json:"count"`
}

func parse(body []byte) ([]element, error) {
	t := html.NewTokenizer(bytes.NewReader(body))
	counts := map[string]int{}
loop:
	for {
		switch t.Next() {
		case html.ErrorToken:
			if t.Err() == io.EOF {
				break loop
			}
			return nil, errors.New("error token found: " + t.Err().Error())
		case html.DoctypeToken:
			counts["!doctype"]++
		case html.StartTagToken, html.SelfClosingTagToken:
			n, _ := t.TagName()
			counts[string(n)]++
		}
	}
	var es []element
	for tag, count := range counts {
		es = append(es, element{
			TagName: tag,
			Count:   count,
		})
	}
	return es, nil
}
