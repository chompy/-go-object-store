package types

// APIResponse defines an API response.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Key     string      `json:"key,omitempty"`
	Objects []APIObject `json:"objects,omitempty"`
}
