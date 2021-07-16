package main

import (
	"github.com/pkg/errors"
	"gitlab.com/contextualcode/go-object-store/store"
	"gitlab.com/contextualcode/go-object-store/types"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func getUserFromCommand(client *store.Client) (*types.User, error) {
	// fetch flag values
	username := userSubCmd.PersistentFlags().Lookup("username").Value.String()
	uid := userSubCmd.PersistentFlags().Lookup("uid").Value.String()
	if client == nil {
		// get config + client
		config, err := loadConfigFromCommand()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		client = store.NewClient(config)
	}
	// load existing user
	var err error
	user := &types.User{}
	if username != "" {
		user, err = client.GetUserByUsername(username)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else if uid != "" {
		user, err = client.GetUser(uid)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if user == nil || user.UID == "" {
		return nil, errors.WithStack(store.ErrMissingUID)
	}
	return user, nil
}

var userSubCmd = &cobra.Command{
	Use:   "user [-u username] [--uid]",
	Short: "User commands.",
}

var userSetCmd = &cobra.Command{
	Use:     "set [-p password] [-g groups] [--disable]",
	Aliases: []string{"create", "make", "new"},
	Short:   "Create or update user.",
	Run: func(cmd *cobra.Command, args []string) {
		// load config + store
		config, err := loadConfigFromCommand()
		if err != nil {
			cliHandleError(err)
		}
		client := store.NewClient(config)
		// flag values
		username := userSubCmd.PersistentFlags().Lookup("username").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()
		groups := cmd.Flags().Lookup("groups").Value.(pflag.SliceValue).GetSlice()
		disable := cmd.Flags().Lookup("disable").Value.String()
		// load existing user
		user, err := getUserFromCommand(client)
		if err != nil && !errors.Is(err, store.ErrNotFound) {
			cliHandleError(err)
		}
		// create new if not exist
		if user == nil {
			user = &types.User{}
			if username == "" || password == "" {
				cliHandleError(store.ErrInvalidArg)
			}
		}
		// update user
		user.Username = username
		store.SetPassword(password, user)
		if err != nil {
			cliHandleError(err)
		}
		user.Groups = groups
		user.Active = disable == "false"
		// store user
		if err := client.SetUser(user); err != nil {
			cliHandleError(err)
		}
		user.PasswordHash = "**redacted**"
		cliResponse([]types.APIObject{user.API()})
	},
}

var userGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get user data.",
	Run: func(cmd *cobra.Command, args []string) {
		user, err := getUserFromCommand(nil)
		cliHandleError(err)
		user.PasswordHash = "**redacted**"
		cliResponse([]types.APIObject{user.API()})
	},
}

var userDeleteCmd = &cobra.Command{
	Use:     "delete [-u username] [--uid]",
	Aliases: []string{"del"},
	Short:   "Delete user.",
	Run: func(cmd *cobra.Command, args []string) {
		user, err := getUserFromCommand(nil)
		cliHandleError(err)
		config, err := loadConfigFromCommand()
		cliHandleError(err)
		client := store.NewClient(config)
		cliHandleError(client.DeleteUser(user))
		cliResponse([]types.APIObject{types.APIObject{"_uid": user.UID}})
	},
}

func init() {
	userSubCmd.PersistentFlags().StringP("username", "u", "", "Username of user to get.")
	userSubCmd.PersistentFlags().String("uid", "", "UID of user to get.")

	userSetCmd.Flags().StringP("password", "p", "", "Set user password.")
	userSetCmd.Flags().StringArrayP("groups", "g", []string{}, "Groups to set user to.")
	userSetCmd.Flags().Bool("disable", false, "Disable user.")

	userSubCmd.AddCommand(userSetCmd)
	userSubCmd.AddCommand(userGetCmd)
	userSubCmd.AddCommand(userDeleteCmd)
	rootCmd.AddCommand(userSubCmd)

}
