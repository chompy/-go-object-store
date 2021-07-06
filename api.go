package main

import (
	"strings"
)

// APIResource defines an action to perform via the API.
type APIResource int

const (
	// APILogin defines login action.
	APILogin APIResource = 1
	// APILogout defines logout action.
	APILogout APIResource = 2
	APIGet    APIResource = 3
	APISet    APIResource = 4
	APIDelete APIResource = 5
)

// APIRequest defines an API request.
type APIRequest struct {
	SessionKey string      `json:"key,omitempty"`
	Username   string      `json:"username,omitempty"`
	Password   string      `json:"password,omitempty"`
	Objects    []APIObject `json:"objects,omitempty"`
}

func (a *APIRequest) sanitizeValues() {
	a.Username = strings.ToLower(strings.TrimSpace(a.Username))
	a.Password = strings.TrimSpace(a.Password)
}

// APIResponse defines an API response.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Key     string      `json:"key,omitempty"`
	Objects []APIObject `json:"objects,omitempty"`
}
