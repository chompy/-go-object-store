package main

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrNotFound            = errors.New("object not found")
	ErrInvalidResource     = errors.New("invalid resource")
	ErrInvalidArg          = errors.New("invalid argument")
	ErrInvalidCreds        = errors.New("user provided invalid login credientials")
	ErrPermission          = errors.New("user does not have permission to access this resource")
	ErrEmptyReponse        = errors.New("empty response")
	ErrMissingUID          = errors.New("object missing uid")
	ErrMissingUsername     = errors.New("user must have an username")
	ErrMissingObject       = errors.New("object can't be nil")
	ErrInvalidCredientials = errors.New("invalid credientials")
	ErrInvalidPassword     = errors.New("invalid password, cannot be empty or less than eight characters")
	ErrUnknown             = errors.New("an unknown error has occured")
)

func errHTTPResponseCode(err error) int {
	switch err {
	case ErrNotFound, ErrInvalidResource:
		{
			return http.StatusNotFound
		}
	case ErrInvalidArg:
		{
			return http.StatusBadRequest
		}
	case ErrInvalidCreds:
		{
			return http.StatusUnauthorized
		}
	case ErrPermission:
		{
			return http.StatusUnauthorized
		}
	}

	return http.StatusInternalServerError
}
