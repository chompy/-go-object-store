package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"gitlab.com/contextualcode/go-object-store/store"
	"gitlab.com/contextualcode/go-object-store/types"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var objSubCmd = &cobra.Command{
	Use:     "object [-u uid] [--user]",
	Aliases: []string{"obj"},
	Short:   "Object store commands.",
}

func getUserFromObjectCommand(client *store.Client) (*types.User, error) {
	// fetch flag value
	username := objSubCmd.PersistentFlags().Lookup("user").Value.String()
	if username == "" {
		return nil, nil
	}
	// load user
	user, err := client.GetUserByUsername(username)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, errors.WithStack(err)
	}
	if user == nil {
		user, err = client.GetUser(username)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return user, nil
}

func getObjectUidsFromCommand() []string {
	uids := objSubCmd.PersistentFlags().Lookup("uid").Value.(pflag.SliceValue).GetSlice()
	// remove duplicates
	out := make([]string, 0)
	for _, uid := range uids {
		hasUid := false
		for _, outUid := range out {
			if outUid == uid {
				hasUid = true
				break
			}
		}
		if !hasUid {
			out = append(out, uid)
		}
	}
	return out
}

var objSetCmd = &cobra.Command{
	Use:   "set [--data]",
	Short: "Set one or more objects.",
	Run: func(cmd *cobra.Command, args []string) {
		// get config + store
		config, err := loadConfigFromCommand()
		cliHandleError(err)
		client := store.NewClient(config)
		// get user to set as
		user, err := getUserFromObjectCommand(client)
		cliHandleError(err)
		// get object data
		data := []byte(cmd.Flags().Lookup("data").Value.String())
		if len(data) == 0 {
			// read data from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				data, err = ioutil.ReadAll(os.Stdin)
				cliHandleError(err)
			}
			if len(data) == 0 {
				cliHandleError(store.ErrMissingObject)
			}
		}
		// parse object data
		objs := make([]types.APIObject, 0)
		cliHandleError(json.Unmarshal(data, &objs))
		// set objects
		out := make([]types.APIObject, 0)
		for _, obj := range objs {
			storeObj := obj.Object()
			cliHandleError(client.Set(storeObj, user))
			out = append(out, storeObj.API())
		}
		cliHandleError(client.Sync())
		cliResponse(out)
	},
}

var objDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del"},
	Short:   "Delete an object.",
	Run: func(cmd *cobra.Command, args []string) {
		// get config + store
		config, err := loadConfigFromCommand()
		cliHandleError(err)
		client := store.NewClient(config)
		// get user to set as
		user, err := getUserFromObjectCommand(client)
		cliHandleError(err)
		// get uids
		uids := getObjectUidsFromCommand()
		if len(uids) == 0 {
			cliHandleError(store.ErrMissingUID)
		}
		// itterate and delete
		out := make([]types.APIObject, 0)
		for _, uid := range uids {
			cliHandleError(client.Delete(&types.Object{UID: uid}, user))
			out = append(out, types.APIObject{"_uid": uid})
		}
		cliResponse(out)
	},
}

var objGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get one or more objects.",
	Run: func(cmd *cobra.Command, args []string) {
		// get config + store
		config, err := loadConfigFromCommand()
		cliHandleError(err)
		client := store.NewClient(config)
		// get user to set as
		user, err := getUserFromObjectCommand(client)
		cliHandleError(err)
		// get uids
		uids := getObjectUidsFromCommand()
		if len(uids) == 0 {
			// use arg 1 if uid flag not set
			if len(args) == 0 {
				cliHandleError(store.ErrMissingUID)
			}
			uids = []string{args[0]}
		}
		// get object
		out := make([]types.APIObject, 0)
		for _, uid := range uids {
			obj, err := client.Get(uid, user)
			cliHandleError(err)
			out = append(out, obj.API())
		}
		cliResponse(out)
	},
}

var objQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "Run a query.",
	Run: func(cmd *cobra.Command, args []string) {
		// get config + store
		config, err := loadConfigFromCommand()
		cliHandleError(err)
		client := store.NewClient(config)
		// get user to set as
		user, err := getUserFromObjectCommand(client)
		cliHandleError(err)
		// get query
		query := strings.Join(args, " ")
		if query == "" {
			cliHandleError(store.ErrInvalidArg)
		}
		// perform query
		cliHandleError(client.Sync())
		res, err := client.Query(query, user)
		cliHandleError(err)
		// get object
		out := make([]types.APIObject, 0)
		for _, obj := range res {
			out = append(out, obj.API())
		}
		cliResponse(out)
	},
}

func init() {
	objSubCmd.PersistentFlags().StringArrayP("uid", "u", []string{}, "UID of object.")
	objSubCmd.PersistentFlags().String("user", "", "User to access object as.")
	objSetCmd.Flags().String("data", "", "JSON object data.")
	objSubCmd.AddCommand(objSetCmd)
	objSubCmd.AddCommand(objDeleteCmd)
	objSubCmd.AddCommand(objGetCmd)
	objSubCmd.AddCommand(objQueryCmd)
}
