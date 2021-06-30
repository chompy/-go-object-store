package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// User defines an user accessing the key/value store.
type User struct {
	UID          string        `json:"uid"`
	Username     string        `json:"username"`
	PasswordHash string        `json:"password_hash"`
	PasswordSalt string        `json:"password_salt"`
	Created      time.Time     `json:"created"`
	Modified     time.Time     `json:"modified"`
	Accessed     time.Time     `json:"accessed"`
	Active       bool          `json:"active"`
	Sessions     []UserSession `json:"sessions"`
}

// UserSession is a single user sign on session.
type UserSession struct {
	Key     string    `json:"key"`
	Salt    string    `json:"salt"`
	IP      string    `json:"ip"`
	Created time.Time `json:"created"`
}

// NewSession creates a new user session.
func (u *User) NewSession(ip string) *UserSession {
	keySalt, err := gonanoid.New()
	if err != nil {
		return nil
	}
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	keySalt = fmt.Sprintf("%s%x", keySalt, randBytes)
	created := time.Now()
	key := sha256.Sum256([]byte(
		fmt.Sprintf("-!-#-b-%s%s%s%s%s%s%s", keySalt, u.UID, u.Username, u.PasswordSalt, u.Created.String(), created.String(), ip),
	))
	return &UserSession{
		Key:     fmt.Sprintf("%x", key),
		Salt:    keySalt,
		Created: created,
		IP:      ip,
	}
}
