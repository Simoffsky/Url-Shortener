package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLink_Validate(t *testing.T) {
	tests := []struct {
		name    string
		link    Link
		wantErr bool
	}{
		{
			name: "valid link",
			link: Link{
				Url:       "http://example.com",
				ShortUrl:  "exmpl",
				ExpiredAt: 3600,
			},
			wantErr: false,
		},
		{
			name: "invalid URL",
			link: Link{
				Url:       "http//invalid",
				ShortUrl:  "exmpl",
				ExpiredAt: 3600,
			},
			wantErr: true,
		},
		{
			name: "empty URL",
			link: Link{
				Url:       "",
				ShortUrl:  "exmpl",
				ExpiredAt: 3600,
			},
			wantErr: true,
		},
		{
			name: "empty ShortUrl",
			link: Link{
				Url:       "http://example.com",
				ShortUrl:  "",
				ExpiredAt: 3600,
			},
			wantErr: true,
		},
		{
			name: "negative TTL",
			link: Link{
				Url:       "http://example.com",
				ShortUrl:  "exmpl",
				ExpiredAt: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.link.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
