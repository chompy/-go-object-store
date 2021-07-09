package main

import (
	"io/ioutil"
	"os"

	"github.com/philippgille/gokv/file"
	"github.com/philippgille/gokv/redis"

	"github.com/philippgille/gokv/syncmap"

	"github.com/philippgille/gokv"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configPath = "config.yaml"

// Config defines configuration values.
type Config struct {
	HTTP struct {
		Port int16 `yaml:"port"`
	} `yaml:"http"`
	Storage struct {
		Type   string                 `yaml:"type"`
		Config map[string]interface{} `yaml:"config"`
	} `yaml:"storage"`
	UserGroups map[string]UserGroup `yaml:"user_groups"`
}

// loadConfig loads config file.
func loadConfig(path string) (*Config, error) {
	// set default
	config := &Config{}
	config.HTTP.Port = 8081
	config.Storage.Type = "memory"
	// load config
	raw, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return config, errors.WithStack(err)
	}
	// parse yaml
	if err := yaml.Unmarshal(raw, config); err != nil {
		return config, errors.WithStack(err)
	}
	return config, nil
}

// storageClient returns the configured Gokv storage client.
func (c *Config) storageClient() gokv.Store {
	switch c.Storage.Type {
	case "redis":
		{
			opts := redis.DefaultOptions
			if c.Storage.Config["address"] != nil {
				opts.Address = c.Storage.Config["address"].(string)
			}
			if c.Storage.Config["password"] != nil {
				opts.Password = c.Storage.Config["password"].(string)
			}
			client, err := redis.NewClient(opts)
			if err != nil {
				logWarnErr(err, "redis client error")
				return syncmap.NewStore(syncmap.DefaultOptions)
			}
			return client
		}
	case "file":
		{
			opts := file.DefaultOptions
			opts.Directory = "data"
			if c.Storage.Config["path"] != nil {
				opts.Directory = c.Storage.Config["path"].(string)
			}
			client, err := file.NewStore(opts)
			if err != nil {
				logWarnErr(err, "file client error")
				return syncmap.NewStore(syncmap.DefaultOptions)
			}
			return client
		}
	default:
		{
			// defaults to memory map
			return syncmap.NewStore(syncmap.DefaultOptions)
		}
	}
}
