package models

import (
	"errors"
	"net/http"
)

type Error struct {
	Err        error
	StatusCode int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func NewError(err error, statusCode int) Error {
	return Error{
		Err:        err,
		StatusCode: statusCode,
	}
}

var (
	ErrLinkNotFound      = NewError(errors.New("link not found"), http.StatusNotFound)
	ErrLinkAlreadyExists = NewError(errors.New("link already exists"), http.StatusConflict)
	ErrLinkExpired       = NewError(errors.New("link expired"), http.StatusGone)
	ErrWrongLinkFormat   = NewError(errors.New("wrong link format"), http.StatusBadRequest)
	ErrCacheMiss         = NewError(errors.New("cache miss"), http.StatusNotFound)
	ErrForbidden         = NewError(errors.New("forbidden"), http.StatusForbidden)

	ErrUserExists      = NewError(errors.New("user already exists"), http.StatusConflict)
	ErrUserNotFound    = NewError(errors.New("user not found"), http.StatusNotFound)
	ErrInvalidPassword = NewError(errors.New("invalid password"), http.StatusUnauthorized)

	ErrTokenExpired = NewError(errors.New("token expired"), http.StatusUnauthorized)
	ErrInvalidToken = NewError(errors.New("invalid token"), http.StatusUnauthorized)
)
