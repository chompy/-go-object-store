package types

import "time"

// APIResponse defines an API response.
type APIResponse struct {
	Success bool        `json:"success"`           // indicates whether the request was successful
	Message string      `json:"message,omitempty"` // response message
	Key     string      `json:"key,omitempty"`     // session key
	Expires time.Time   `json:"expires,omitempty"` // key expiration time
	Objects []APIObject `json:"objects,omitempty"` // list of objects returned by the request
}
