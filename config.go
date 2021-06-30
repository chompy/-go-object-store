package main

import (
	"io/ioutil"
	"os"

	"github.com/philippgille/gokv/redis"

	"github.com/philippgille/gokv/syncmap"

	"github.com/philippgille/gokv"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configPath = "config.yaml"

type Config struct {
	HTTP struct {
		Port int16 `yaml:"port"`
	} `yaml:"http"`
	Storage struct {
		Type   string                 `yaml:"type"`
		Config map[string]interface{} `yaml:"config"`
	} `yaml:"storage"`
}

func loadConfig() (*Config, error) {
	// set default
	config := &Config{}
	config.HTTP.Port = 8081
	config.Storage.Type = "memory"
	// load config
	raw, err := ioutil.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return config, errors.WithStack(err)
	}
	// parse yaml
	if err := yaml.Unmarshal(raw, config); err != nil {
		return config, errors.WithStack(err)
	}
	return config, nil
}

// StorageClient returns the configured Gokv storage client.
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
	default:
		{
			// defaults to memory map
			return syncmap.NewStore(syncmap.DefaultOptions)
		}
	}
}
