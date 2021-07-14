package client

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrResponse = errors.New("error reponse")
)
