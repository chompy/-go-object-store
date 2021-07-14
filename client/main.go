package client

import (
	"gitlab.com/contextualcode/go-object-store/types"
)

// Login authenticates given user with object store API and returns session key.
func Login(username string, password string) (string, error) {
	req := types.APIRequest{
		Username: username,
		Password: password,
	}
	resp, err := request(types.APILogin, req)
	if err != nil {
		return "", err
	}
	return resp.Key, nil
}

// Get fetches object of given UID from API.
func Get(uid string, key string) (*types.Object, error) {
	req := types.APIRequest{
		SessionKey: key,
		Objects:    []types.APIObject{types.APIObject{"_uid": uid}},
	}
	resp, err := request(types.APIGet, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Objects) > 0 {
		return resp.Objects[0].Object(), nil
	}
	return nil, ErrNotFound
}

// Set stores given object to API.
func Set(obj *types.Object, key string) error {
	return nil
}

// Delete deletes object of given UID from API.
func Delete(uid string, key string) error {
	return nil
}

// Query queries
func Query(query string, key string) ([]*types.IndexObject, error) {
	return nil, nil
}
