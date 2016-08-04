package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseEnvAllowsPoundSign(t *testing.T) {
	assert := assert.New(t)
	val, err := ParseEnv("FOO=bar#baz")
	assert.NoError(err)
	assert.Equal(map[string]string{"FOO": "bar#baz"}, val)
}
