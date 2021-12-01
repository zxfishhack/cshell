package store

import (
	"github.com/kevinburke/ssh_config"
	"github.com/zxfishhack/cshell/pkg/sshc"
	"os"
	"path/filepath"
)

var sshConfig *ssh_config.Config

func Init() (err error) {
	sshConfig, err = sshc.LoadSSHConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	return
}

func SaveSSHConfig() (err error) {
	return sshc.SaveSSHConfig(filepath.Join(os.Getenv("HOME"), ".ssh", "config"), sshConfig)
}

func GetSSHHostList() (hosts []string) {
	if sshConfig != nil {
		for _, v := range sshConfig.Hosts {
			for _, p := range v.Patterns {
				if p.String() == "*" {
					continue
				}
				hosts = append(hosts, p.String())
			}
		}
	}
	return
}
