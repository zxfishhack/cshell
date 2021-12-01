package store

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestGetSSHHostList(t *testing.T) {
	assert.NilError(t, Init())
	t.Log(GetSSHHostList("个人", false))
}
