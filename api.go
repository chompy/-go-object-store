package main

import (
	"fmt"
	"strings"
)

// APIResource defines an action to perform via the API.
type APIResource int

const (
	// APILogin defines login action.
	APILogin APIResource = 1
	// APILogout defines logout action.
	APILogout APIResource = 2
	// APIGet defines get
	APIGet    APIResource = 3
	APISet    APIResource = 4
	APIDelete APIResource = 5
	APIQuery  APIResource = 6
)

// Name returns string name for API resource.
func (r APIResource) Name() string {
	switch r {
	case APILogin:
		{
			return "LOGIN"
		}
	case APILogout:
		{
			return "LOGOUT"
		}
	case APIGet:
		{
			return "GET"
		}
	case APISet:
		{
			return "SET"
		}
	case APIDelete:
		{
			return "DELETE"
		}
	case APIQuery:
		{
			return "QUERY"
		}
	}
	return ""
}

// APIRequest defines an API request.
type APIRequest struct {
	IP         string      `json:"-"`
	SessionKey string      `json:"key,omitempty"`
	Username   string      `json:"username,omitempty"`
	Password   string      `json:"password,omitempty"`
	Objects    []APIObject `json:"objects,omitempty"`
	Query      string      `json:"query,omitempty"`
}

func (a *APIRequest) sanitizeValues() {
	a.Username = strings.ToLower(strings.TrimSpace(a.Username))
	a.Password = strings.TrimSpace(a.Password)
}

// ObjectUIDs return list of object uids in api request.
func (a APIRequest) ObjectUIDs() []string {
	out := make([]string, 0)
	for _, o := range a.Objects {
		uid := o.Object().UID
		if uid != "" {
			out = append(out, uid)
		}
	}
	return out
}

// Log logs the request.
func (a APIRequest) Log(res APIResource) {
	userIdentity := anonymousUser
	if a.Username != "" {
		userIdentity = a.Username
	} else if a.SessionKey != "" {
		userIdentity = a.SessionKey
	}
	objString := ""
	if len(a.Objects) > 0 {
		objList := a.ObjectUIDs()
		newObjCt := len(a.Objects) - len(objList)
		objString = fmt.Sprintf(" %s", strings.Join(objList, ","))
		if newObjCt > 0 {
			objString += fmt.Sprintf(" + %d new", newObjCt)
		}
	} else if a.Query != "" {
		objString += " " + a.Query
	}
	logInfo(
		fmt.Sprintf("@%s - %s%s", userIdentity, res.Name(), objString),
	)
}

// APIResponse defines an API response.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Key     string      `json:"key,omitempty"`
	Objects []APIObject `json:"objects,omitempty"`
}
