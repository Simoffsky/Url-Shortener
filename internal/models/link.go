package models

import (
	"errors"
	"net/url"
)

type Link struct {
	Url      string `json:"url"`
	ShortUrl string `json:"short_url"`
	TTL      int    `json:"ttl"`
}

func (l *Link) Validate() error {
	var textErr string

	if l.Url == "" {
		textErr += "Url is empty;"
	} else {
		_, err := url.ParseRequestURI(l.Url)
		if err != nil {
			textErr += "Url is invalid;"
		}
	}

	if l.ShortUrl == "" {
		textErr += "ShortUrl is empty;"
	}
	if l.TTL < 0 {
		textErr += "TTL is negative;"
	}
	if textErr != "" {
		return errors.New(textErr)
	}
	return nil
}
