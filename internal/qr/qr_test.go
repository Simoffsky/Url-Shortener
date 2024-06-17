package qr

import (
	"testing"
	"url-shorter/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestGenerateValidURL(t *testing.T) {
	link := "https://example.com"
	qrCode, err := Generate(link, 256)
	assert.NoError(t, err, "Generate() should not return an error for a valid URL")
	assert.NotNil(t, qrCode, "QR code should not be nil")
	assert.NotEmpty(t, qrCode, "QR code should not be empty")
}

func TestGenerateInvalidURL(t *testing.T) {
	link := "not a valid link"
	_, err := Generate(link, 256)
	assert.Equal(t, models.ErrWrongLinkFormat, err, "Generate() должна возвращать ошибку ErrWrongLinkFormat для невалидного URL")
}
