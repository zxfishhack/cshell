package main

import (
	"os/exec"
)

func main() {
	cmd := exec.Command("osascript", "-e", `tell application "Terminal" to do script "touch ~/vvv"`)
	cmd.Start()
	if cmd.Process != nil {
		cmd.Process.Release()
	}
}
