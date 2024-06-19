package qr

import (
	"testing"
	"url-shorter/internal/models"
	"url-shorter/internal/qr/cache"
	"url-shorter/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetQRCode(t *testing.T) {
	testCases := []struct {
		name        string
		link        string
		size        int
		cacheHit    bool
		expectError bool
	}{
		{
			name:        "Cache miss success",
			link:        "http://example.com",
			size:        256,
			cacheHit:    false,
			expectError: false,
		},
		{
			name:        "Cache hit success",
			link:        "http://example.com",
			size:        256,
			cacheHit:    true,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := &cache.MockCache{}
			logger := &log.MockLogger{}

			if tc.cacheHit {
				cache.On("Get", mock.Anything).Return([]byte("cached QR code"), nil)
			} else {
				cache.On("Get", mock.Anything).Return([]byte{}, models.ErrCacheMiss)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}

			service := NewDefaultQRService(cache, logger)

			_, err := service.GetQRCode(tc.link, tc.size)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			cache.AssertExpectations(t)
		})
	}
}
