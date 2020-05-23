package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
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
	randomJsonConfigData = `
{"aaa": "aaa"}
`
)

func TestConfigNotConflictedData(t *testing.T) {
	for i, data := range []string{notMatchingYamlConfigData, randomJsonConfigData} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			config, err := GetConfig([]byte(data))
			assert.Error(t, err, "expected error not present")
			assert.Equal(t, config, Config{}, "expected empty config")
			//tassert.Contains(t, buf.String(), "unsupported protocol scheme")
		})
	}
}
func TestConfigConflictedData(t *testing.T) {
	for i, data := range []string{conflictingYamlConfigData} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Panics(t, func() { GetConfig([]byte(data)) }, "Expected panic did not happened")
		})
	}
}

func TestConfigokData(t *testing.T) {
	for i, data := range []string{okYamlConfigData0, okYamlConfigData1} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			config, err := GetConfig([]byte(data))
			assert.NoError(t, err, "expected error not present")
			assert.NotEqual(t, config, Config{}, "expected empty config")
			assert.Contains(t, config.Destination, "test")
		})
	}
}
