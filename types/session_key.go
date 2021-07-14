package types

import (
	"time"
)

// SessionKey defines a session key.
type SessionKey struct {
	Key     string    `json:"key"`
	Expires time.Time `json:"expires"`
}
