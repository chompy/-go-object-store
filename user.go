package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const passwordMinLength = 8
const sessionTimeout = 3600

// User defines an user accessing the key/value store.
type User struct {
	UID      string    `json:"uid"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Accessed time.Time `json:"accessed"`
	Active   bool      `json:"active"`
	Groups   []string  `json:"groups"`
}

func generateObjectUID() string {
	uid, err := gonanoid.New()
	if err != nil {
		return ""
	}
	return uid
}

// NewUser creates a new user.
func NewUser() *User {
	uid, err := gonanoid.New()
	if err != nil {
		return nil
	}
	return &User{
		UID:      uid,
		Username: uid,
		Password: "",
		Created:  time.Now(),
		Modified: time.Now(),
		Accessed: time.Now(),
		Active:   true,
		Groups:   make([]string, 0),
	}
}

func (u *User) generatePasswordHash(password string) (string, error) {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(hashedPasswordBytes), nil
}

// SetPassword sets the user's password.
func (u *User) SetPassword(password string) error {
	// validate password
	// TODO use regexp to check for characters
	if len(password) < passwordMinLength {
		return errors.WithStack(ErrInvalidPassword)
	}
	var err error
	u.Password, err = u.generatePasswordHash(password)
	return errors.WithStack(err)
}

// CheckPassword returns true if given password matches hash.
func (u *User) CheckPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func (u *User) generateSessionKey() string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	key := sha256.Sum256([]byte(
		fmt.Sprintf("%x%s%s", randBytes, u.UID, u.Username),
	))
	return fmt.Sprintf("%x", key)
}
