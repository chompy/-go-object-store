package main

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var objSubCmd = &cobra.Command{
	Use:     "object [-u uid] [--user]",
	Aliases: []string{"obj"},
	Short:   "Object store commands.",
}

func getUserFromObjectCommand(store *Store) (*User, error) {
	// fetch flag value
	username := objSubCmd.PersistentFlags().Lookup("user").Value.String()
	if username == "" {
		return nil, nil
	}
	// load user
	user, err := store.GetUserByUsername(username)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, errors.WithStack(err)
	}
	if user == nil {
		user, err = store.GetUser(username)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return user, nil
}

var objSetCmd = &cobra.Command{
	Use:     "set [--data]",
	Aliases: []string{"del"},
	Short:   "Set an object.",
	Run: func(cmd *cobra.Command, args []string) {
		// get config + store
		config, err := loadConfigFromCommand()
		cliHandleError(err)
		store := NewStore(config)
		// get user to set as
		user, err := getUserFromObjectCommand(store)
		cliHandleError(err)
		// get uid
		uid := objSubCmd.PersistentFlags().Lookup("uid").Value.String()
		// set object
		obj := &Object{}
		if uid != "" {
			obj, err = store.Get(uid, user)
			cliHandleError(err)
		}
		// TODO set data
		obj.Data = make(map[string]interface{})
		cliHandleError(store.Set(obj, user))
		cliResponse([]APIObject{obj.API()})
	},
}

func init() {
	objSubCmd.PersistentFlags().StringP("uid", "u", "", "UID of object.")

	objSubCmd.AddCommand(objSetCmd)
}
