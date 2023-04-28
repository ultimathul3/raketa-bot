package handler

import "errors"

var (
	errInvalidUrlInput       = errors.New("invalid URL")
	errInvalidPriceInput     = errors.New("invalid price")
	errGettingUrlFromStorage = errors.New("error getting url from storage")
	errUnknownUserRole       = errors.New("unknown user role")
)
