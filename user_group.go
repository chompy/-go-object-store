package main

import (
	"github.com/caibirdme/yql"
	"github.com/pkg/errors"
	"gitlab.com/contextualcode/go-object-store/types"
)

const (
	permGet    = "get"
	permSet    = "set"
	permUpdate = "update"
	permDelete = "delete"
)

// UserGroup defines access parameters for a user group.
type UserGroup struct {
	Get      interface{}          `yaml:"get"`    // read
	Set      interface{}          `yaml:"set"`    // create new
	Update   interface{}          `yaml:"update"` // update existing (that user is not author of)
	Delete   interface{}          `yaml:"delete"` // delete
	compiled map[string]yql.Ruler `yaml:"-"`
}

func (g *UserGroup) getPerm(permType string) interface{} {
	switch permType {
	case permGet:
		{
			return g.Get
		}
	case permSet:
		{
			return g.Set
		}
	case permUpdate:
		{
			return g.Update
		}
	case permDelete:
		{
			return g.Delete
		}
	}
	return nil
}

func (g *UserGroup) compile() error {
	g.compiled = make(map[string]yql.Ruler)
	for _, permType := range []string{permGet, permSet, permUpdate, permDelete} {
		perm := g.getPerm(permType)
		var err error
		switch v := perm.(type) {
		case string:
			{
				g.compiled[permType], err = yql.Rule(v)
				if err != nil {
					return errors.WithStack(err)
				}
				break
			}
		}
	}
	return nil
}

func (g UserGroup) check(permType string, o *types.IndexObject) (bool, error) {
	if o == nil {
		return false, errors.WithStack(ErrMissingObject)
	}
	perm := g.getPerm(permType)
	switch perm := perm.(type) {
	case string:
		{
			if g.compiled[permType] == nil {
				if err := g.compile(); err != nil {
					return false, errors.WithStack(err)
				}
			}
			match, err := g.compiled[permType].Match(o.QueryMap())
			return match, errors.WithStack(err)
		}
	case bool:
		{
			return perm, nil
		}
	default:
		{
			return false, nil
		}
	}
}

// CanGet returns true if group permission allows reading given object.
func (g UserGroup) CanGet(o *types.IndexObject) (bool, error) {
	return g.check(permGet, o)
}

// CanSet returns true if group permission allow creation of given object.
func (g UserGroup) CanSet(o *types.IndexObject) (bool, error) {
	return g.check(permSet, o)
}

// CanUpdate returns true if group permission allows updating the given object.
func (g UserGroup) CanUpdate(o *types.IndexObject) (bool, error) {
	return g.check(permUpdate, o)
}

// CanDelete returns true if group permission allows deleting the given object.
func (g UserGroup) CanDelete(o *types.IndexObject) (bool, error) {
	return g.check(permDelete, o)
}
