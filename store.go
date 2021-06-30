package main

import (
	"log"
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
	client    gokv.Store
	sync      sync.Mutex
	index     []*IndexObject
	indexSync sync.Mutex
}

// NewStore creates a new object store from given configuration.
func NewStore(c *Config) *Store {
	if c == nil {
		// use memory store by default
		return &Store{
			client: syncmap.NewStore(syncmap.DefaultOptions),
		}
	}
	s := &Store{
		client: c.storageClient(),
	}
	return s
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
		log.Println("TESTdsfsd", hasChange)
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
	return o, nil
}

// Set stores object.
func (s *Store) Set(o *Object, u *User) error {
	defer s.sync.Unlock()
	s.sync.Lock()
	o.Modified = time.Now()
	if err := s.client.Set(objectPrefix+o.UID, o); err != nil {
		return errors.WithStack(err)
	}
	s.addIndex(o.Index())
	return nil
}

// Delete deletes object from store.
func (s *Store) Delete(o *Object, u *User) error {
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
		obj.Data["created"] = obj.Created.UTC().Unix()
		obj.Data["modified"] = obj.Modified.UTC().Unix()
		match, err := ruler.Match(obj.Data)
		if err != nil {
			if strings.Contains(err.Error(), "not provided") {
				continue
			}
			return nil, errors.WithStack(err)
		}
		if match {
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
