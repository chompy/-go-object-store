package main

import (
	"gitlab.com/contextualcode/go-object-store/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const passwordMinLength = 8
const sessionTimeout = 3600
const anonymousUser = "anonymous"

func generateObjectUID() string {
	uid, err := gonanoid.New()
	if err != nil {
		return ""
	}
	return uid
}

func generatePasswordHash(password string) (string, error) {
	if len(password) < passwordMinLength {
		return "", errors.WithStack(ErrInvalidPassword)
	}
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(hashedPasswordBytes), nil
}

func setPassword(password string, u *types.User) error {
	hash, err := generatePasswordHash(password)
	if err != nil {
		return errors.WithStack(err)
	}
	u.PasswordHash = hash
	return nil
}

func checkPassword(password string, hash string) bool {
	if hash == "" {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
