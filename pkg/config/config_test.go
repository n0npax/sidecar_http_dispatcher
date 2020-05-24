package config

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	okYamlConfigData0 = `
destination: test
`
	okYamlConfigData1 = `
rewrites:
  dev:
    patch:
      - key: Host
        val: example.com
      - key: dispatched
        val: true
    destination: http://example.com
  qa:
    patch:
      - key: dispatched
        val: true
    destination: http://destination-app.default.svc.cluster.local
key: environment
destination: test
`
	conflictingYamlConfigData = `
rewrites: single_strin_instead_of_map
key: environment
destination: http://destination-app.default.svc.cluster.local
`

	notMatchingYamlConfigData = `
aaa: aaa
`
	randomJSONConfigData = `
{"aaa": "aaa"}
`
)

func TestBrokenConfig(t *testing.T) {
	for i, data := range []string{notMatchingYamlConfigData, randomJSONConfigData, conflictingYamlConfigData} {
		data := data

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			config, err := buildConfig([]byte(data))
			assert.Error(t, err, "expected error not present")
			assert.Equal(t, config, Config{}, "expected empty config")
		})
	}
}
func TestConfigokData(t *testing.T) {
	for i, data := range []string{okYamlConfigData0, okYamlConfigData1} {
		data := data

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			config, err := buildConfig([]byte(data))
			assert.NoError(t, err, "expected error not present")
			assert.NotEqual(t, config, Config{}, "expected empty config")
			assert.Contains(t, config.Destination, "test")
		})
	}
}

func TestReadConfigFile(t *testing.T) {
	orgEnv := os.Getenv("SIDECAR_CONFIG")

	defer func() { os.Setenv("SIDECAR_CONFIG", orgEnv) }()

	for _, path := range []string{"/etc/hosts", "/etc/hostname"} {
		path := path
		t.Run(path, func(t *testing.T) {
			os.Setenv("SIDECAR_CONFIG", path)
			_, err := readConfigFile()
			assert.NoError(t, err, "expected error not present")
		})
	}
}

func TestGetConfigFail1(t *testing.T) {
	readConfigFileF = func() ([]byte, error) { return []byte{}, errors.New("test") }
	logFatalfF = func(string, ...interface{}) { panic("test") }

	assert.Panics(t, func() { GetConfig() }, "expected mocked panic not present")
}

func TestGetConfigFail2(t *testing.T) {
	buildConfigF = func([]byte) (Config, error) { return Config{}, errors.New("test") }
	logFatalfF = func(string, ...interface{}) { panic("test") }

	assert.Panics(t, func() { GetConfig() }, "expected mocked panic not present")
}

func TestGetConfigMockedEmpty(t *testing.T) {
	buildConfigF = func([]byte) (Config, error) { return Config{}, nil }
	readConfigFileF = func() ([]byte, error) { return []byte{}, nil }
	config := GetConfig()

	assert.Equal(t, config, Config{}, "expected empty config for mocked input")
}
