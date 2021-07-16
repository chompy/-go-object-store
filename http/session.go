package http

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"gitlab.com/contextualcode/go-object-store/types"
	"golang.org/x/crypto/bcrypt"
)

const sessionTimeout = 3600

var sessions []*UserSession

// UserSession is a single user sign on session.
type UserSession struct {
	UserUID string
	Key     string
	IP      string
	Created time.Time
}

// Expires returns the session expiration time.
func (u *UserSession) Expires() time.Time {
	return u.Created.Add(time.Second * sessionTimeout)
}

func newSession(u *types.User, ip string) (*UserSession, *types.SessionKey) {
	key := generateSessionKey(u)
	keyHash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		logWarnErr(err, "")
		return nil, nil
	}
	created := time.Now()
	return &UserSession{
			UserUID: u.UID,
			Key:     string(keyHash),
			Created: time.Now(),
			IP:      ip,
		}, &types.SessionKey{
			Key:     key,
			Expires: created.Add(time.Second * sessionTimeout),
		}
}

func checkSessions() {
	sessEnd := time.Second * sessionTimeout
	for i, sess := range sessions {
		if sess.Created.Add(sessEnd).Before(time.Now()) {
			sessions = append(sessions[:i], sessions[i+1:]...)
			checkSessions()
			return
		}
	}
}

func getSessionFromKey(key string) *UserSession {
	for _, sess := range sessions {
		if sess.checkKey(key) {
			return sess
		}
	}
	return nil
}

func (s *UserSession) checkKey(key string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(s.Key), []byte(key)); err != nil {
		return false
	}
	return true
}

func generateSessionKey(u *types.User) string {
	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	key := sha256.Sum256([]byte(
		fmt.Sprintf("%x%s%s", randBytes, u.UID, u.Username),
	))
	return fmt.Sprintf("%x", key)
}
