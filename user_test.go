package main

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/go-object-store/types"
)

func TestPassword(t *testing.T) {
	u := &types.User{
		Username: "testuser",
	}
	if setPassword("test123", u) == nil {
		t.Error("expected password error")
		return
	}
	if err := setPassword("test1234", u); err != nil {
		t.Error(err)
		return
	}
	if !checkPassword("test1234", u.PasswordHash) {
		t.Error("expected password check valid")
		return
	}
	if checkPassword("test5678", u.PasswordHash) {
		t.Error("expected password check invalid")
		return
	}
}

func TestUserStore(t *testing.T) {

	s := NewStore(nil)

	// create user
	u := &types.User{
		Username: "testuser",
	}
	setPassword("test1234", u)

	// store user
	if err := s.SetUser(u); err != nil {
		t.Error(err)
		return
	}

	// fetch user by uid
	sUser, err := s.GetUser(u.UID)
	if err != nil {
		t.Error(err)
		return
	}
	if sUser.UID != u.UID || sUser.Username != u.Username {
		t.Error("unexpected stored user")
		return
	}

	// fetch user by username
	sUser, err = s.GetUserByUsername(u.Username)
	if err != nil {
		t.Error(err)
		return
	}
	if sUser.UID != u.UID || sUser.Username != u.Username {
		t.Error("unexpected stored user")
		return
	}

	// delete user
	if err := s.DeleteUser(u); err != nil {
		t.Error(err)
		return
	}
	_, err = s.GetUser(u.UID)
	if !errors.Is(err, ErrNotFound) {
		t.Error("expected user not found error")
		return
	}

}

func TestUserPermission(t *testing.T) {

	c := &Config{}
	c.UserGroups = map[string]UserGroup{
		"admin": UserGroup{
			Get:    true,
			Set:    true,
			Update: true,
			Delete: true,
		},
		"test1": UserGroup{
			Get: "type = 'page'",
		},
		"test2": UserGroup{
			Set: "type = 'food'",
		},
		"test3": UserGroup{
			Update: false,
		},
	}

	o1 := &types.Object{
		Data: map[string]interface{}{
			"type":    "page",
			"content": "Test page 1",
		},
	}

	m, err := c.UserGroups["test1"].CanGet(o1.Index())
	if err != nil {
		t.Error(err)
		return
	}
	if !m {
		t.Error("expected match")
		return
	}

	m, err = c.UserGroups["test1"].CanSet(o1.Index())
	if err != nil {
		t.Error(err)
		return
	}
	if m {
		t.Error("expected no match")
		return
	}

	m, err = c.UserGroups["test2"].CanSet(o1.Index())
	if err != nil {
		t.Error(err)
		return
	}
	if m {
		t.Error("expected no match")
		return
	}

	m, err = c.UserGroups["test3"].CanUpdate(o1.Index())
	if err != nil {
		t.Error(err)
		return
	}
	if m {
		t.Error("expected no match")
		return
	}

}

func TestStoreWithUser(t *testing.T) {

	c := &Config{}
	c.UserGroups = map[string]UserGroup{
		"admin": UserGroup{
			Get:    true,
			Set:    true,
			Update: true,
			Delete: true,
		},
		"food_reader": UserGroup{
			Get: "type = 'food'",
		},
		"food_writer": UserGroup{
			Set: "type = 'food'",
		},
		"food_updater": UserGroup{
			Update: "type = 'food'",
		},
		"food_deleter": UserGroup{
			Delete: "type = 'food'",
		},
	}

	s := NewStore(c)
	u := &types.User{
		UID: "sample_user",
	}
	o := &types.Object{
		Data: map[string]interface{}{
			"type": "food",
		},
	}

	// no user group, expect permission error
	if err := s.Set(o, u); !errors.Is(err, ErrPermission) {
		t.Error("expected permission error")
		return
	}

	// add user group, expect pass
	u.Groups = []string{"food_writer"}
	if err := s.Set(o, u); err != nil {
		t.Error(err)
		return
	}

	// change type, user only has permission to set type 'food', expect permission error
	o.Data["type"] = "page"
	if err := s.Set(o, u); !errors.Is(err, ErrPermission) {
		t.Error("expected permission error")
		return
	}

	// add a new key+value (with type food), expect pass
	o.Data["type"] = "food"
	o.Data["name"] = "Pizza"
	if err := s.Set(o, u); err != nil {
		t.Error(err)
		return
	}

	// try to read, user owns object, expect pass
	if _, err := s.Get(o.UID, u); err != nil {
		t.Error(err)
		return
	}

	// reset permissions, try to update, user owns object, expect fail
	u.Groups = []string{}
	if err := s.Set(o, u); !errors.Is(err, ErrPermission) {
		t.Error("expected permission error")
		return
	}

	// add updater permission, expect pass
	u.Groups = []string{"food_updater"}
	if err := s.Set(o, u); err != nil {
		t.Error(err)
		return
	}

	// replace with permission, expect pass
	u.Groups = []string{"food_writer"}
	if err := s.Set(o, u); err != nil {
		t.Error(err)
		return
	}

	// create a new object not owned by user, try to read, expect fail
	u.Groups = []string{"food_writer"}
	o2 := &types.Object{
		Data: map[string]interface{}{
			"type": "food",
			"name": "Salad",
		},
	}
	s.Set(o2, nil)
	if _, err := s.Get(o2.UID, u); !errors.Is(err, ErrPermission) {
		t.Error("expected permission error")
		return
	}

	// grant user food reader access and try again to read, expect pass
	u.Groups = []string{"food_reader"}
	if _, err := s.Get(o2.UID, u); err != nil {
		t.Error(err)
		return
	}

	// try to update new object, expect fail
	o2.Data["name"] = "Pie"
	if err := s.Set(o2, u); !errors.Is(err, ErrPermission) {
		t.Error("expected permission error")
		return
	}

	// grant permission food updater and try again, expect pass
	u.Groups = []string{"food_reader", "food_writer", "food_updater"}
	if err := s.Set(o2, u); err != nil {
		t.Error(err)
		return
	}

}

func TestUserSession(t *testing.T) {

	// create test user
	u := &types.User{
		Username: "testuser1",
	}
	setPassword("test1234", u)

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
	sess = getSessionFromKey(key)
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
