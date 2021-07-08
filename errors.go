package main

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrNotFound            = errors.New("object not found")
	ErrObjectNotSpecified  = errors.New("object not specified")
	ErrInvalidResource     = errors.New("invalid resource")
	ErrInvalidArg          = errors.New("invalid argument")
	ErrInvalidCreds        = errors.New("user provided invalid login credientials")
	ErrInvalidSession      = errors.New("invalid session key")
	ErrPermission          = errors.New("user does not have permission to access this resource")
	ErrEmptyReponse        = errors.New("empty response")
	ErrMissingUID          = errors.New("object missing uid")
	ErrMissingUsername     = errors.New("user must have an username")
	ErrMissingObject       = errors.New("object can't be nil")
	ErrInvalidCredientials = errors.New("invalid credientials")
	ErrInvalidPassword     = errors.New("invalid password, cannot be empty or less than eight characters")
	ErrUnknown             = errors.New("an unknown error has occured")
	ErrAPIInvalidMethod    = errors.New("method not supported")
)

func errHTTPResponseCode(err error) int {
	switch errors.Cause(err) {
	case ErrNotFound, ErrInvalidResource:
		{
			return http.StatusNotFound
		}
	case ErrInvalidArg, ErrObjectNotSpecified:
		{
			return http.StatusBadRequest
		}
	case ErrInvalidCreds, ErrInvalidSession:
		{
			return http.StatusUnauthorized
		}
	case ErrPermission:
		{
			return http.StatusUnauthorized
		}
	case ErrAPIInvalidMethod:
		{
			return http.StatusMethodNotAllowed
		}
	}
	return http.StatusInternalServerError
}
