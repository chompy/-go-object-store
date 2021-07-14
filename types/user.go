package types

import (
	"time"
)

// User defines an user accessing the key/value store.
type User struct {
	UID          string    `json:"uid"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Created      time.Time `json:"created"`
	Modified     time.Time `json:"modified"`
	Active       bool      `json:"active"`
	Groups       []string  `json:"groups"`
}

// API converts user to API object.
func (u *User) API() APIObject {
	out := make(APIObject)
	out["_uid"] = u.UID
	out["_username"] = u.Username
	out["_active"] = u.Active
	out["_groups"] = u.Groups
	out["_created"] = u.Created.Format(time.RFC3339)
	out["_modified"] = u.Modified.Format(time.RFC3339)
	return out
}
