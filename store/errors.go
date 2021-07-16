package store

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound            = errors.New("object not found")
	ErrObjectNotSpecified  = errors.New("object not specified")
	ErrInvalidArg          = errors.New("invalid argument")
	ErrInvalidCreds        = errors.New("user provided invalid login credientials")
	ErrPermission          = errors.New("user does not have permission to access this resource")
	ErrMissingUID          = errors.New("object missing uid")
	ErrMissingUsername     = errors.New("user must have an username")
	ErrMissingObject       = errors.New("object can't be nil")
	ErrInvalidCredientials = errors.New("invalid credientials")
	ErrInvalidPassword     = errors.New("invalid password, cannot be empty or less than eight characters")
	ErrUnknown             = errors.New("an unknown error has occured")
	ErrInvalidUsername     = errors.New("invalid or missing username")
)
