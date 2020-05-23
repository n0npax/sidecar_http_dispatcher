package config

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Key         string          `yaml:"key"`
	Destination string          `yaml:"destination"`
	Rewrites    map[string]Rule `yaml:"rewrites"`
}

type Rule struct {
	Destination string `yaml:"destination"`
	Patch       []struct {
		Key string `yaml:"key"`
		Val string `yaml:"val"`
	} `yaml:"patch"`
}

func GetConfig(data []byte) (Config, error) {
	t := Config{}

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		panic(err)
	}
	if t.Destination == "" {
		err = errors.New("config doesn't have any destination")
	}
	return t, err
}
