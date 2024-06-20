package models

import "errors"

var (
	ErrLinkNotFound      = errors.New("link not found")
	ErrLinkAlreadyExists = errors.New("link already exists")
	ErrLinkExpired       = errors.New("link expired")
	ErrWrongLinkFormat   = errors.New("wrong link format")
	ErrCacheMiss         = errors.New("cache miss")
	ErrForbidden         = errors.New("forbidden")

	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)
