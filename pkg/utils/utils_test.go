package utils

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestOpenSSH(t *testing.T) {
	err := OpenSSH(DefaultTerminal, "rancher")
	assert.NilError(t, err)
}
