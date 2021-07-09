package main

import (
	"strings"
	"sync"
	"time"

	"github.com/caibirdme/yql"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/syncmap"

	"github.com/pkg/errors"
)

const (
	userPrefix     = "user_"
	usernamePrefix = "username_"
	objectPrefix   = "obj_"
	indexName      = "index"
)

// Store is the key/value store interface.
type Store struct {
	client     gokv.Store
	sync       sync.Mutex
	index      []*IndexObject
	indexSync  sync.Mutex
	userGroups map[string]UserGroup
}

// NewStore creates a new object store from given configuration.
func NewStore(c *Config) *Store {
	if c == nil {
		// use memory store by default
		return &Store{
			client:     syncmap.NewStore(syncmap.DefaultOptions),
			userGroups: make(map[string]UserGroup),
		}
	}
	s := &Store{
		client:     c.storageClient(),
		userGroups: c.UserGroups,
	}
	return s
}

func (s *Store) getUserGroups(u *User) []UserGroup {
	out := make([]UserGroup, 0)
	if u == nil {
		return out
	}
	for k, v := range s.userGroups {
		for _, name := range u.Groups {
			if name == k {
				out = append(out, v)
			}
		}
	}
	return out
}

func (s *Store) getRaw(k string, o interface{}) error {
	found, err := s.client.Get(k, o)
	if !found {
		return errors.WithStack(ErrNotFound)
	} else if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Store) addIndex(o *IndexObject) {
	s.indexSync.Lock()
	defer s.indexSync.Unlock()
	for i := range s.index {
		if s.index[i].UID == o.UID {
			s.index[i] = o
			return
		}
	}
	s.index = append(s.index, o)
}

func (s *Store) deleteIndex(o *Object) {
	s.indexSync.Lock()
	defer s.indexSync.Unlock()
	for i := range s.index {
		if s.index[i].UID == o.UID {
			s.index = append(s.index[:i], s.index[i+1:]...)
			break
		}
	}
}

func (s *Store) commitIndex() error {
	s.indexSync.Lock()
	defer s.indexSync.Unlock()
	if err := s.client.Set(indexName, s.index); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Store) checkPermission(perm string, u *User, o *IndexObject) error {
	if o == nil {
		return errors.WithStack(ErrMissingObject)
	}
	if u == nil {
		return nil
	}
	// if user is author then they can 'get' the object
	if perm == permGet && o.Author == u.UID {
		return nil
	}
	// itterate groups and see if any allow permission
	userGroups := s.getUserGroups(u)
	for _, userGroup := range userGroups {
		match, err := userGroup.check(perm, o)
		if err != nil {
			return errors.WithStack(err)
		}
		if match {
			return nil
		}
	}
	// if user owns object then they are allowed to update it provided they have 'set' permission
	if perm == permUpdate && o.Author == u.UID {
		return s.checkPermission(permSet, u, o)
	}
	return errors.WithStack(ErrPermission)
}

// Sync syncs the local memory index with the remote store index.
func (s *Store) Sync() error {
	s.indexSync.Lock()
	defer s.indexSync.Unlock()
	remoteIndex := make([]*IndexObject, 0)
	if err := s.getRaw(indexName, &remoteIndex); err != nil {
		if !errors.Is(err, ErrNotFound) {
			return errors.WithStack(err)
		}
	}
	hasChange := false
	// check if remote index has items that local does not and check if matching items have been modified
	for _, remoteIndexItem := range remoteIndex {
		hasLocal := false
		for i, localIndexItem := range s.index {
			if localIndexItem.UID == remoteIndexItem.UID {
				hasLocal = true
				if remoteIndexItem.Modified.After(localIndexItem.Modified) {
					s.index[i] = remoteIndexItem
					continue
				}
				hasChange = true
			}
		}
		if !hasLocal {
			hasChange = true
			s.index = append(s.index, remoteIndexItem)
		}
	}
	// check if local index has items that remote does not
	if !hasChange {
		for _, localIndexItem := range s.index {
			hasItem := false
			for _, remoteIndexItem := range remoteIndex {
				if localIndexItem.UID == remoteIndexItem.UID {
					hasItem = true
					break
				}
			}
			if !hasItem {
				hasChange = true
				break
			}
		}
	}
	// update remote only if local has changes
	if hasChange {
		s.indexSync.Unlock()
		err := s.commitIndex()
		s.indexSync.Lock()
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Index returns index data.
func (s *Store) Index() ([]IndexObject, error) {
	s.indexSync.Lock()
	defer s.indexSync.Unlock()
	out := make([]IndexObject, 0)
	for _, o := range s.index {
		out = append(out, *o)
	}
	return out, nil
}

// Get retrieves object from store.
func (s *Store) Get(uid string, u *User) (*Object, error) {
	o := &Object{}
	if err := s.getRaw(objectPrefix+uid, o); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := s.checkPermission(permGet, u, o.Index()); err != nil {
		return nil, errors.WithStack(err)
	}
	return o, nil
}

// Set stores object.
func (s *Store) Set(o *Object, u *User) error {
	if o == nil {
		return errors.WithStack(ErrMissingObject)
	}
	// new object
	isNew := false
	if o.UID == "" {
		isNew = true
		o.UID = generateObjectUID()
		o.Author = ""
		if u != nil {
			o.Author = u.UID
		}
		o.Created = time.Now()
	}
	defer s.sync.Unlock()
	s.sync.Lock()
	// check against previous existing object
	if u != nil {
		var existingObj *Object
		if !isNew {
			var err error
			existingObj, err = s.Get(o.UID, nil)
			if err != nil && !errors.Is(err, ErrNotFound) {
				return errors.WithStack(err)
			}
		}
		if existingObj == nil || isNew {
			// if no existing object then use 'set' permission
			if err := s.checkPermission(permSet, u, o.Index()); err != nil {
				return errors.WithStack(err)
			}
		} else if existingObj != nil {
			// if existing object then use 'update' permission
			if err := s.checkPermission(permUpdate, u, existingObj.Index()); err != nil {
				return errors.WithStack(err)
			}
			if err := s.checkPermission(permUpdate, u, o.Index()); err != nil {
				return errors.WithStack(err)
			}
			// author and created aren't allowed to be changed
			o.Author = existingObj.Author
			o.Created = existingObj.Created
		}
	}
	o.Modified = time.Now()
	o.Modifier = ""
	if u != nil {
		o.Modifier = u.UID
	}
	if err := s.client.Set(objectPrefix+o.UID, o); err != nil {
		return errors.WithStack(err)
	}
	s.addIndex(o.Index())
	return nil
}

// Delete deletes object from store.
func (s *Store) Delete(o *Object, u *User) error {
	if o.UID == "" {
		return errors.WithStack(ErrMissingUID)
	}
	if err := s.checkPermission(permDelete, u, o.Index()); err != nil {
		return errors.WithStack(err)
	}
	defer s.sync.Unlock()
	s.sync.Lock()
	if err := s.client.Delete(objectPrefix + o.UID); err != nil {
		return errors.WithStack(err)
	}
	s.deleteIndex(o)
	return nil
}

// Query returns indexed objects based on provided query match.
func (s *Store) Query(q string, u *User) ([]IndexObject, error) {
	ruler, err := yql.Rule(q)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	out := make([]IndexObject, 0)
	for _, obj := range s.index {
		match, err := ruler.Match(obj.QueryMap())
		if err != nil {
			if strings.Contains(err.Error(), "not provided") {
				continue
			}
			return nil, errors.WithStack(err)
		}
		if match {
			if err := s.checkPermission(permGet, u, obj); err != nil {
				if errors.Is(err, ErrPermission) {
					continue
				}
				return nil, errors.WithStack(err)
			}
			out = append(out, *obj)
		}
	}
	return out, nil
}

// GetUser retrieves user from store.
func (s *Store) GetUser(uid string) (*User, error) {
	u := &User{}
	if err := s.getRaw(userPrefix+uid, u); err != nil {
		return nil, errors.WithStack(err)
	}
	return u, nil
}

// GetUserByUsername retrieves user from store by their username.
func (s *Store) GetUserByUsername(username string) (*User, error) {
	username = sanitizeUsername(username)
	u := &User{}
	if err := s.getRaw(usernamePrefix+username, u); err != nil {
		return nil, errors.WithStack(err)
	}
	return u, nil
}

// SetUser stores given user.
func (s *Store) SetUser(u *User) error {
	if u.UID == "" {
		return errors.WithStack(ErrMissingUID)
	}
	if u.Username == "" {
		return errors.WithStack(ErrMissingUsername)
	}
	if u.Password == "" {
		return errors.WithStack(ErrInvalidPassword)
	}
	defer s.sync.Unlock()
	s.sync.Lock()
	if err := s.client.Set(userPrefix+u.UID, u); err != nil {
		return errors.WithStack(err)
	}
	if err := s.client.Set(usernamePrefix+u.Username, u); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteUser deletes given user from store.
func (s *Store) DeleteUser(u *User) error {
	defer s.sync.Unlock()
	s.sync.Lock()
	if err := s.client.Delete(userPrefix + u.UID); err != nil {
		return errors.WithStack(err)
	}
	if err := s.client.Delete(usernamePrefix + u.Username); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
