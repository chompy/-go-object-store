package main

import (
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func getUserFromCommand(store *Store) (*User, error) {
	// fetch flag values
	username := userSubCmd.PersistentFlags().Lookup("username").Value.String()
	uid := userSubCmd.PersistentFlags().Lookup("uid").Value.String()
	if store == nil {
		// get config + store
		config, err := loadConfigFromCommand()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		store = NewStore(config)
	}
	// load existing user
	var err error
	user := &User{}
	if username != "" {
		user, err = store.GetUserByUsername(username)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else if uid != "" {
		user, err = store.GetUser(uid)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if user == nil || user.UID == "" {
		return nil, errors.WithStack(ErrMissingUID)
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
		store := NewStore(config)
		// flag values
		username := userSubCmd.PersistentFlags().Lookup("username").Value.String()
		password := cmd.Flags().Lookup("password").Value.String()
		groups := cmd.Flags().Lookup("groups").Value.(pflag.SliceValue).GetSlice()
		disable := cmd.Flags().Lookup("disable").Value.String()
		// load existing user
		user, err := getUserFromCommand(store)
		if err != nil && !errors.Is(err, ErrNotFound) {
			cliHandleError(err)
		}
		// create new if not exist
		if user == nil {
			user = NewUser()
			if username == "" || password == "" {
				cliHandleError(ErrInvalidArg)
			}
		}
		// update user
		user.Username = username
		if err := user.SetPassword(password); err != nil {
			cliHandleError(err)
		}
		user.Groups = groups
		user.Active = disable == "false"
		// store
		if err := store.SetUser(user); err != nil {
			cliHandleError(err)
		}
		user.Password = "**redacted**"
		cliResponse([]APIObject{user.API()})
	},
}

var userGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get user data.",
	Run: func(cmd *cobra.Command, args []string) {
		user, err := getUserFromCommand(nil)
		cliHandleError(err)
		user.Password = "**redacted**"
		cliResponse([]APIObject{user.API()})
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
		store := NewStore(config)
		cliHandleError(store.DeleteUser(user))
		cliResponse([]APIObject{APIObject{"_uid": user.UID}})
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
