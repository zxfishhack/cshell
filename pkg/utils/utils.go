//go:build darwin

package utils

import (
	"fmt"
	"os/exec"
)

func OpenSSH(name string) (err error) {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "ssh %s"`, name))
	err = cmd.Start()
	if cmd.Process != nil {
		cmd.Process.Release()
	}
	return
}
