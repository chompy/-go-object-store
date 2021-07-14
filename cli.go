package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/go-object-store/types"

	"github.com/spf13/cobra"
)

func loadConfigFromCommand() (*Config, error) {
	path := rootCmd.PersistentFlags().Lookup("config").Value.String()
	if path == "" {
		path = configPath
	}
	config, err := loadConfig(path)
	return config, errors.WithStack(err)
}

func cliHandleError(err error) {
	if err != nil {
		resp := types.APIResponse{
			Success: false,
			Message: err.Error(),
		}
		respJSON, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(respJSON))

		if rootCmd.PersistentFlags().Lookup("verbose").Value.String() == "true" {
			panic(err)
		}
		os.Exit(1)
	}
}

func cliResponse(objs []types.APIObject) {
	resp := types.APIResponse{
		Success: true,
		Objects: objs,
	}
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(respJSON))
}

var rootCmd = &cobra.Command{
	Use:     "cc_store [-c config]",
	Version: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TEST")
	},
}

var serveCmd = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"listen", "start"},
	Short:   "Start API web server.",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfigFromCommand()
		cliHandleError(err)
		cliHandleError(listen(config))
	},
}

func init() {
	rootCmd.AddCommand(objSubCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.PersistentFlags().StringP("config", "c", configPath, "Set config yaml path.")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show stack trace on error.")
}
