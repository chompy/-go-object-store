package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const configPath = "config.yaml"

// Config defines configuration.
type Config struct {
	Storage []struct {
		Name   string                 `yaml:"name"`
		Type   string                 `yaml:"type"`
		Config map[string]interface{} `yaml:"config"`
	} `yaml:"storage"`
}

var config *Config

// Get returns configuration.
func Get() *Config {
	if config == nil {
		configRaw, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(configRaw, config); err != nil {
			panic(err)
		}
	}
	return config
}
