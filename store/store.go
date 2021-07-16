package store

import (
	"strings"
	"sync"
	"time"

	"github.com/caibirdme/yql"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/syncmap"
	"gitlab.com/contextualcode/go-object-store/types"

	"github.com/pkg/errors"
)

const (
	userPrefix     = "user_"
	usernamePrefix = "username_"
	objectPrefix   = "obj_"
	indexName      = "index"
)

// Client is the key/value store interface.
type Client struct {
	store      gokv.Store
	sync       sync.Mutex
	index      []*types.IndexObject
	indexSync  sync.Mutex
	userGroups map[string]UserGroup
}

// NewClient creates a new object store client from given configuration.
func NewClient(c *Config) *Client {
	if c == nil {
		// use memory store by default
		return &Client{
			store:      syncmap.NewStore(syncmap.DefaultOptions),
			userGroups: make(map[string]UserGroup),
		}
	}
	s := &Client{
		store:      c.storageClient(),
		userGroups: c.UserGroups,
	}
	return s
}

func (c *Client) getUserGroups(u *types.User) []UserGroup {
	out := make([]UserGroup, 0)
	if u == nil {
		return out
	}
	for k, v := range c.userGroups {
		for _, name := range u.Groups {
			if name == k {
				out = append(out, v)
			}
		}
	}
	return out
}

func (c *Client) getRaw(k string, o interface{}) error {
	found, err := c.store.Get(k, o)
	if !found {
		return errors.WithStack(ErrNotFound)
	} else if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Client) addIndex(o *types.IndexObject) {
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

func (s *Client) deleteIndex(o *types.Object) {
	s.indexSync.Lock()
	defer s.indexSync.Unlock()
	for i := range s.index {
		if s.index[i].UID == o.UID {
			s.index = append(s.index[:i], s.index[i+1:]...)
			break
		}
	}
}

func (s *Client) commitIndex() error {
	s.indexSync.Lock()
	defer s.indexSync.Unlock()
	if err := s.store.Set(indexName, s.index); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Client) checkPermission(perm string, u *types.User, o *types.IndexObject) error {
	if o == nil {
		return errors.WithStack(ErrMissingObject)
	}
	if u == nil {
		return nil
	}
	// if user is author then they can 'get' the object
	if perm == permGet && u.UID != "" && o.Author == u.UID {
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
func (s *Client) Sync() error {
	s.indexSync.Lock()
	defer s.indexSync.Unlock()
	remoteIndex := make([]*types.IndexObject, 0)
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
func (c *Client) Index() ([]types.IndexObject, error) {
	c.indexSync.Lock()
	defer c.indexSync.Unlock()
	out := make([]types.IndexObject, 0)
	for _, o := range c.index {
		out = append(out, *o)
	}
	return out, nil
}

// Get retrieves object from store.
func (c *Client) Get(uid string, u *types.User) (*types.Object, error) {
	o := &types.Object{}
	if err := c.getRaw(objectPrefix+uid, o); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := c.checkPermission(permGet, u, o.Index()); err != nil {
		return nil, errors.WithStack(err)
	}
	return o, nil
}

// Set stores object.
func (c *Client) Set(o *types.Object, u *types.User) error {
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
	defer c.sync.Unlock()
	c.sync.Lock()
	// check against previous existing object
	if u != nil {
		var existingObj *types.Object
		if !isNew {
			var err error
			existingObj, err = c.Get(o.UID, nil)
			if err != nil && !errors.Is(err, ErrNotFound) {
				return errors.WithStack(err)
			}
		}
		if existingObj == nil || isNew {
			// if no existing object then use 'set' permission
			if err := c.checkPermission(permSet, u, o.Index()); err != nil {
				return errors.WithStack(err)
			}
		} else if existingObj != nil {
			// if existing object then use 'update' permission
			if err := c.checkPermission(permUpdate, u, existingObj.Index()); err != nil {
				return errors.WithStack(err)
			}
			if err := c.checkPermission(permUpdate, u, o.Index()); err != nil {
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
	if err := c.store.Set(objectPrefix+o.UID, o); err != nil {
		return errors.WithStack(err)
	}
	c.addIndex(o.Index())
	return nil
}

// Delete deletes object from store.
func (c *Client) Delete(o *types.Object, u *types.User) error {
	if o.UID == "" {
		return errors.WithStack(ErrMissingUID)
	}
	if err := c.checkPermission(permDelete, u, o.Index()); err != nil {
		return errors.WithStack(err)
	}
	defer c.sync.Unlock()
	c.sync.Lock()
	if err := c.store.Delete(objectPrefix + o.UID); err != nil {
		return errors.WithStack(err)
	}
	c.deleteIndex(o)
	o.UID = ""
	return nil
}

// Query returns indexed objects based on provided query match.
func (c *Client) Query(q string, u *types.User) ([]types.IndexObject, error) {
	ruler, err := yql.Rule(q)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	out := make([]types.IndexObject, 0)
	for _, obj := range c.index {
		match, err := ruler.Match(obj.QueryMap())
		if err != nil {
			if strings.Contains(err.Error(), "not provided") {
				continue
			}
			return nil, errors.WithStack(err)
		}
		if match {
			if err := c.checkPermission(permGet, u, obj); err != nil {
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
func (c *Client) GetUser(uid string) (*types.User, error) {
	u := &types.User{}
	if err := c.getRaw(userPrefix+uid, u); err != nil {
		return nil, errors.WithStack(err)
	}
	return u, nil
}

// GetUserByUsername retrieves user from store by their username.
func (c *Client) GetUserByUsername(username string) (*types.User, error) {
	username = sanitizeUsername(username)
	u := &types.User{}
	if err := c.getRaw(usernamePrefix+username, u); err != nil {
		return nil, errors.WithStack(err)
	}
	return u, nil
}

// SetUser stores given user.
func (c *Client) SetUser(u *types.User) error {
	// require a username
	if u.Username == "" {
		return errors.WithStack(ErrMissingUsername)
	}
	// generate user id if not exists
	if u.UID == "" {
		u.UID = generateObjectUID()
		u.Created = time.Now()
		u.Modified = time.Now()
		u.Active = true
	}
	defer c.sync.Unlock()
	c.sync.Lock()
	if err := c.store.Set(userPrefix+u.UID, u); err != nil {
		return errors.WithStack(err)
	}
	if err := c.store.Set(usernamePrefix+u.Username, u); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DeleteUser deletes given user from store.
func (c *Client) DeleteUser(u *types.User) error {
	defer c.sync.Unlock()
	c.sync.Lock()
	if err := c.store.Delete(userPrefix + u.UID); err != nil {
		return errors.WithStack(err)
	}
	if err := c.store.Delete(usernamePrefix + u.Username); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
