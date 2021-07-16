package store

import (
	"gitlab.com/contextualcode/go-object-store/types"
	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const passwordMinLength = 8

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

// SetPassword hashes given password and sets it to given user.
func SetPassword(password string, u *types.User) error {
	hash, err := generatePasswordHash(password)
	if err != nil {
		return errors.WithStack(err)
	}
	u.PasswordHash = hash
	return nil
}

// CheckPassword checks given password against given hash.
func CheckPassword(password string, hash string) bool {
	if hash == "" {
		return false
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}
