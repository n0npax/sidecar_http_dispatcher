package utils_test

import (
	"testing"

	"github.com/n0npax/sidecar_http_dispatcher/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvPATH(t *testing.T) {
	defaultVal := "FOO"
	val := utils.GetEnv("PATH", defaultVal)
	assert.NotEqual(t, val, defaultVal, "not expected fallback for PATH env has happened")
}

func TestGetEnvFOO(t *testing.T) {
	defaultVal := "FOO"
	val := utils.GetEnv("FOO", defaultVal)
	assert.Equal(t, val, defaultVal, "ENV value FOO expected to fallback, but didn't")
}
