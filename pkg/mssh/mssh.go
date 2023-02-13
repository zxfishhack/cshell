package mssh

import (
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/zxfishhack/cshell/pkg/store"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func Init() (err error) {
	err = store.Init()
	if err != nil {
		return
	}

	return
}

var outputMtx sync.Mutex

func Execute(tagList, hostList []string, password string, cmd []string) (err error) {
	var filterHosts []string
	if len(tagList) == 0 {
		filterHosts = store.GetSSHHostList("", false)
	} else {
		for _, tag := range tagList {
			hosts := store.GetSSHHostList(tag, false)
			filterHosts = append(filterHosts, hosts...)
		}
	}
	if len(hostList) > 0 {
		var tmp []string
		for _, host := range hostList {
			if host == "" {
				continue
			}
			if exp, e := regexp.Compile(host); e == nil {
				for _, h := range filterHosts {
					if exp.MatchString(h) {
						tmp = append(tmp, h)
					}
				}
			} else {
				for _, h := range filterHosts {
					if h == host {
						tmp = append(tmp, h)
					}
				}
			}
		}
		filterHosts = nil
		uniq := make(map[string]bool)
		for _, h := range tmp {
			if _, ok := uniq[h]; !ok {
				filterHosts = append(filterHosts, h)
				uniq[h] = true
			}
		}
	}
	signers := make(map[string]ssh.Signer)
	log.Printf("ready to execute command [%s] on remote [%s]",
		strings.Join(cmd, " "), strings.Join(filterHosts, ","))
	var wg sync.WaitGroup
	for _, h := range filterHosts {
		res := store.GetHostKV(h)
		hostname := h
		port := 22
		username := ""
		homeDir := ""
		if u, e := user.Current(); e == nil {
			username = u.Username
			homeDir = u.HomeDir
		}
		idFile := ""
		for _, kv := range res {
			if kv.Key == "HostName" {
				hostname = kv.Value
			} else if kv.Key == "Port" {
				port, _ = strconv.Atoi(kv.Value)
			} else if kv.Key == "User" {
				username = kv.Value
			} else if kv.Key == "IdentityFile" {
				idFile = kv.Value
			}
		}
		var signer ssh.Signer
		if idFile != "" {
			var fn string
			if strings.HasPrefix(idFile, "~") {
				fn = strings.Replace(idFile, "~", homeDir, 1)
			} else if fn, err = filepath.Abs(idFile); err != nil {
				fn = idFile
			}
			if b, e := os.ReadFile(fn); e == nil {
				if signer, e = ssh.ParsePrivateKey(b); e == nil {
					signers[fn] = signer
				} else if _, ok := e.(*ssh.PassphraseMissingError); ok {
					if signer, e = ssh.ParsePrivateKeyWithPassphrase(b, []byte(password)); e == nil {
						signers[fn] = signer
					} else if errors.Is(e, x509.IncorrectPasswordError) {
						fmt.Printf("please enter passphrase for identity file %s [%s]:", idFile, fn)
						outputMtx.Lock()
						if passwd, e := terminal.ReadPassword(int(os.Stdin.Fd())); e == nil {
							if signer, e = ssh.ParsePrivateKeyWithPassphrase(b, passwd); e == nil {
								signers[fn] = signer
							} else {
								log.Printf("parse private key [%s] failed %s", idFile, e.Error())
							}
						}
						fmt.Println()
						outputMtx.Unlock()
					}
				} else {
					log.Printf("parse private key [%s] failed %s", idFile, e.Error())
				}
			} else {
				log.Printf("open file %s failed %s.", fn, e.Error())
			}
		}
		am := make([]ssh.AuthMethod, 0)
		if signer != nil {
			am = append(am, ssh.PublicKeys(signer))
		}
		if password != "" {
			am = append(am, ssh.Password(password))
		}
		am = append(am, ssh.PasswordCallback(func() (secret string, err error) {
			outputMtx.Lock()
			defer outputMtx.Unlock()
			fmt.Printf("please enter passsword for host %s [%s]:", h, hostname)
			passwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			return string(passwd), err
		}))
		conf := &ssh.ClientConfig{
			User: username,
			Auth: am,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
		wg.Add(1)
		go func(conf *ssh.ClientConfig, addr, host string) {
			defer wg.Done()
			if e := sshCmd(conf, addr, host, cmd); e != nil {
				log.Printf("ERROR: execute command on %s [%s] failed %s.", host, addr, e.Error())
			}
		}(conf, fmt.Sprintf("%s:%d", hostname, port), h)
	}
	err = nil
	wg.Wait()
	return
}

func sshCmd(conf *ssh.ClientConfig, addr string, host string, cmd []string) (err error) {
	var outputBytes []byte
	c, err := ssh.Dial("tcp", addr, conf)
	if err != nil {
		log.Printf("ERROR: dial %s [%s] failed %s.", host, addr, err.Error())
		return
	}
	defer c.Close()
	defer func() {
		outputMtx.Lock()
		fmt.Println()
		log.Printf("execute command [%s] on %s [%s] done.", strings.Join(cmd, " "), host, addr)
		_, _ = os.Stdout.Write(outputBytes)
		fmt.Println()
		outputMtx.Unlock()
	}()
	if len(cmd) == 0 {
		return
	}
	var s *ssh.Session
	if s, err = c.NewSession(); err != nil {
		return
	}
	defer s.Close()
	cmdLine := cmd[0]
	if len(cmd) > 1 {
		cmdLine += " \"" + strings.Join(cmd[1:], `" "`) + "\""
	}
	outputBytes, err = s.CombinedOutput(cmdLine)
	return
}
