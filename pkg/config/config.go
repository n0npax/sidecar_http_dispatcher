package config

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/n0npax/sidecar_http_dispatcher/pkg/utils"

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

func buildConfig(data []byte) (Config, error) {
	t := Config{}

	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return Config{}, err
	}

	if t.Destination == "" {
		return Config{}, errors.New("config doesn't have any destination")
	}

	return t, err
}

func readConfigFile() ([]byte, error) {
	return ioutil.ReadFile(utils.GetEnv("SIDECAR_CONFIG", "config.yaml"))
}

// mapping functions to vars to provide testing possibility
var (
	readConfigFileF = readConfigFile // nolint
	buildConfigF    = buildConfig    // nolint
	logFatalfF      = log.Fatalf     // nolint
)

func GetConfig() Config {
	fileContent, err := readConfigFileF()
	if err != nil {
		logFatalfF("Critical error when reading config file: %v", err)
	}

	config, err := buildConfigF(fileContent)
	if err != nil {
		logFatalfF("Critical error when parcing config file: %v", err)
	}

	log.Printf("Using config: %v", config)

	return config
}
