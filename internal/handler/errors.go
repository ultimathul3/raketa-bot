package handler

import "errors"

var (
	errInvalidUrlInput       = errors.New("invalid URL")
	errInvalidUserIdInput    = errors.New("invalid user ID")
	errGettingUrlFromStorage = errors.New("error getting url from storage")
)
