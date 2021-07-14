package client

import (
	"github.com/pkg/errors"
	"gitlab.com/contextualcode/go-object-store/types"
)

// Login authenticates given user with object store API and returns session key.
func Login(username string, password string) (*types.SessionKey, error) {
	req := types.APIRequest{
		Username: username,
		Password: password,
	}
	resp, err := request(types.APILogin, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !resp.Success {
		return nil, errors.WithStack(errors.WithMessage(ErrResponse, resp.Message))
	}
	return &types.SessionKey{
		Key:     resp.Key,
		Expires: resp.Expires,
	}, nil
}

// Get fetches objects from store API.
func Get(uids []string, key string) ([]*types.Object, error) {
	objs, err := requestObj(types.APIGet, uids, key)
	return objs, errors.WithStack(err)
}

// Set stores given objects to store API.
func Set(objs []*types.Object, key string) error {
	apiObjs := make([]types.APIObject, 0)
	for _, obj := range objs {
		apiObjs = append(apiObjs, obj.API())
	}
	req := types.APIRequest{
		SessionKey: key,
		Objects:    apiObjs,
	}
	resp, err := request(types.APISet, req)
	if err != nil {
		return errors.WithStack(err)
	}
	if !resp.Success {
		return errors.WithStack(errors.WithMessage(ErrResponse, resp.Message))
	}
	return nil
}

// Delete deletes object of given UID from store API.
func Delete(uids []string, key string) error {
	_, err := requestObj(types.APIDelete, uids, key)
	return errors.WithStack(err)
}

// Query queries the store API.
func Query(query string, key string) ([]*types.IndexObject, error) {
	req := types.APIRequest{
		SessionKey: key,
		Query:      query,
	}
	resp, err := request(types.APIQuery, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !resp.Success {
		return nil, errors.WithStack(errors.WithMessage(ErrResponse, resp.Message))
	}
	objs := make([]*types.IndexObject, 0)
	for _, obj := range resp.Objects {
		objs = append(objs, obj.Object().Index())
	}
	return objs, nil
}
