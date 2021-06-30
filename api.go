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
	AuthKey  string `json:"auth"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *APIRequest) sanitizeValues() {
	a.Username = strings.ToLower(strings.TrimSpace(a.Username))
	a.Password = strings.TrimSpace(a.Password)
}

// APIResponse defines an API response.
type APIResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
	Objects []*Object `json:"objects,omitempty"`
	User    *User     `json:"user,omitempty"`
}
