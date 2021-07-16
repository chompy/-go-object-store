package http

import (
	"testing"
	"time"

	"gitlab.com/contextualcode/go-object-store/store"
	"gitlab.com/contextualcode/go-object-store/types"
)

func TestUserSession(t *testing.T) {

	// create test user
	u := &types.User{
		Username: "testuser1",
	}
	store.SetPassword("test1234", u)

	// init session
	sessions = make([]*UserSession, 0)
	sess, key := newSession(u, "127.0.0.1")
	sessions = append(sessions, sess)

	// ensure session check doesn't remove session prematurely
	if len(sessions) != 1 {
		t.Error("unexpected session size")
		return
	}
	checkSessions()
	if len(sessions) != 1 {
		t.Error("unexpected session size")
		return
	}

	// ensure session can be obtained with session key
	sess = getSessionFromKey(key.Key)
	if sess == nil {
		t.Error("session key did not match")
		return
	}

	// check that session is removed when older than sessionTimeout
	sess.Created = time.Now().Add(-(time.Second * sessionTimeout * 2))
	checkSessions()
	if len(sessions) > 0 {
		t.Error("unexpected session size")
		return
	}

}
