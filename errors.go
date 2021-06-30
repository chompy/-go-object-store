package main

import (
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrNotFound        = errors.New("object not found")
	ErrInvalidResource = errors.New("invalid resource")
	ErrInvalidArg      = errors.New("invalid argument")
	ErrInvalidCreds    = errors.New("user provided invalid login credientials")
	ErrPermission      = errors.New("user does not have permission to access this resource")
	ErrEmptyReponse    = errors.New("empty response")
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
