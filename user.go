package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/pkg/errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const passwordMinLength = 8

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
	Groups       []string      `json:"groups"`
	Sessions     []UserSession `json:"sessions"`
}

// UserSession is a single user sign on session.
type UserSession struct {
	Key     string    `json:"key"`
	Salt    string    `json:"salt"`
	IP      string    `json:"ip"`
	Created time.Time `json:"created"`
}

// NewUser creates a new user.
func NewUser() *User {
	uid, err := gonanoid.New()
	if err != nil {
		return nil
	}
	return &User{
		UID:          uid,
		Username:     uid,
		PasswordHash: "",
		PasswordSalt: "",
		Created:      time.Now(),
		Modified:     time.Now(),
		Accessed:     time.Now(),
		Active:       true,
		Groups:       make([]string, 0),
		Sessions:     make([]UserSession, 0),
	}
}

func (u *User) generatePasswordSalt() string {
	randBytes := make([]byte, 32)
	hashBytes := sha256.Sum256([]byte(
		fmt.Sprintf("%x%s", randBytes, u.UID),
	))
	return fmt.Sprintf("%x", hashBytes)
}

func (u *User) generatePasswordHash(password string, salt string) string {
	hashBytes := sha256.Sum256([]byte(
		fmt.Sprintf("-1-n-$-%s%s", salt, password),
	))
	return fmt.Sprintf("%x", hashBytes)
}

// SetPassword sets the user's password.
func (u *User) SetPassword(password string) error {
	// validate password
	// TODO use regexp to check for characters
	if len(password) < passwordMinLength {
		return errors.WithStack(ErrInvalidPassword)
	}
	u.PasswordSalt = u.generatePasswordSalt()
	u.PasswordHash = u.generatePasswordHash(password, u.PasswordSalt)
	return nil
}

// CheckPassword returns true if given password matches hash.
func (u *User) CheckPassword(password string) bool {
	return u.generatePasswordHash(password, u.PasswordSalt) == u.PasswordHash
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
	return &UserSession{
		Key:     u.generateSessionKey(keySalt, created, ip),
		Salt:    keySalt,
		Created: created,
		IP:      ip,
	}
}

func (u *User) generateSessionKey(salt string, created time.Time, ip string) string {
	key := sha256.Sum256([]byte(
		fmt.Sprintf("-!-#-b-%s%s%s%s%s%s%s", salt, u.UID, u.Username, u.PasswordSalt, u.Created.String(), created.String(), ip),
	))
	return fmt.Sprintf("%x", key)
}
