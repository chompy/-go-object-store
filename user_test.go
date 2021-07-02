package main

import (
	"testing"

	"github.com/pkg/errors"
)

func TestPassword(t *testing.T) {
	u := NewUser()
	u.Username = "testuser"
	if u.SetPassword("test123") == nil {
		t.Error("expected password error")
		return
	}
	if err := u.SetPassword("test1234"); err != nil {
		t.Error(err)
		return
	}
	if !u.CheckPassword("test1234") {
		t.Error("expected password check valid")
		return
	}
	if u.CheckPassword("test5678") {
		t.Error("expected password check invalid")
		return
	}
}

func TestUserStore(t *testing.T) {

	s := NewStore(nil)

	// create user
	u := NewUser()
	u.Username = "testuser"
	u.SetPassword("test1234")

	// stoer user
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

	u := NewUser()
	u.Username = "testuser"
	u.Groups = []string{"test1", "test2"}

	o1 := NewObject(nil)
	o1.Data["type"] = "page"
	o1.Data["content"] = "Test page 1"

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
	u := NewUser()

	o := NewObject(u)
	o.Data["type"] = "food"

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
	o2 := NewObject(nil)
	o2.Data["type"] = "food"
	o2.Data["name"] = "Salad"
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
