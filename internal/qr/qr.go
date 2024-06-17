package qr

// generates a QR code for the given link

import (
	"net/url"
	"url-shorter/internal/models"

	qrcode "github.com/skip2/go-qrcode"
)

// returns raw png bytes of the QR code
// image has imgSize x imgSize pixels
func Generate(link string, imgSize int) ([]byte, error) {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return nil, models.ErrWrongLinkFormat
	}

	rawPng, err := qrcode.Encode(link, qrcode.Medium, imgSize)

	if err != nil {
		return nil, err
	}
	return rawPng, nil
}
