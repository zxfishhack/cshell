package sshc

import (
	"github.com/kevinburke/ssh_config"
	"io/ioutil"
	"os"
	"sync"
)

var sshConfigMutex sync.RWMutex

func LoadSSHConfig(sshFile string) (cfg *ssh_config.Config, err error) {
	f, err := os.Open(sshFile)
	if err != nil {
		return
	}
	return ssh_config.Decode(f)
}

func SaveSSHConfig(sshFile string, cfg *ssh_config.Config) (err error) {
	return ioutil.WriteFile(sshFile, []byte(cfg.String()), 0600)
}
