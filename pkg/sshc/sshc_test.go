package sshc

import (
	"gotest.tools/v3/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSshConfig(t *testing.T) {
	h, err := os.UserHomeDir()
	if err != nil {
		return
	}
	cfg, err := LoadSSHConfig(filepath.Join(h, ".ssh", "config"))

	assert.NilError(t, err)
	t.Log(cfg.Hosts)
}

func TestSaveSshConfig(t *testing.T) {
	h, err := os.UserHomeDir()
	if err != nil {
		return
	}
	cfg, err := LoadSSHConfig(filepath.Join(h, ".ssh", "config"))
	t.Logf("\n%s", cfg)
	for _, h := range cfg.Hosts {
		for _, p := range h.Patterns {
			t.Log(p)
			t.Log(h)
		}
	}
	assert.NilError(t, err)
	err = SaveSSHConfig("temp", cfg)
	assert.NilError(t, err)
}
