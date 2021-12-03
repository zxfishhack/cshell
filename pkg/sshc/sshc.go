package sshc

import (
	"errors"
	"github.com/zxfishhack/ssh_config"
	"io/ioutil"
	"os"
	"sync"
)

var sshConfigMutex sync.RWMutex

func LoadSSHConfig(sshFile string) (cfg *ssh_config.Config, err error) {
	f, err := os.Open(sshFile)
	// 处理ssh config不存在的问题
	if errors.Is(err, os.ErrNotExist) {
		return &ssh_config.Config{
			Hosts: make([]*ssh_config.Host, 0),
		}, nil
	}
	if err != nil {
		return
	}
	return ssh_config.Decode(f)
}

func SaveSSHConfig(sshFile string, cfg *ssh_config.Config) (err error) {
	return ioutil.WriteFile(sshFile, []byte(cfg.String()), 0600)
}
