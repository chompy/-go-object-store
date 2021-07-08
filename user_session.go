package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

var sessions []*UserSession

// UserSession is a single user sign on session.
type UserSession struct {
	UserUID string
	Key     string
	IP      string
	Created time.Time
}

// NewSession creates a new session.
func (u *User) NewSession(ip string) (*UserSession, string) {
	key := u.generateSessionKey()
	keyHash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		logWarnErr(err, "")
		return nil, ""
	}
	s := &UserSession{
		UserUID: u.UID,
		Key:     string(keyHash),
		Created: time.Now(),
		IP:      ip,
	}
	return s, key
}

// CheckKey checks if given key matches session key.
func (s *UserSession) CheckKey(key string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(s.Key), []byte(key)); err != nil {
		return false
	}
	return true
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
		if sess.CheckKey(key) {
			return sess
		}
	}
	return nil
}
