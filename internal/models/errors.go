package models

import "errors"

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
	ErrWrongLinkFormat   = errors.New("wrong link format")
	ErrCacheMiss         = errors.New("cache miss")
)
