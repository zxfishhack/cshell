//+build darwin
//go:build darwin
// +build darwin

package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const terminalScript = `
tell application "Terminal"
    set newWindow to (do script "ssh %s")
    tell the application named "Terminal"
        activate
    end tell
end tell
`

const iTerm2Script = `
tell application "iTerm2"
    set newWindow to (create window with default profile command "ssh %s")
    tell the application named "iTerm2"
        activate
    end tell
end tell
`

type TerminalType string

const (
	DefaultTerminal TerminalType = "Terminal"
	ITerm2                       = "iTerm2"
)

func OpenSSH(termType TerminalType, name string) (err error) {
	fn := filepath.Join(os.TempDir(), "run")
	script := ""
	switch termType {
	case ITerm2:
		script = fmt.Sprintf(iTerm2Script, name)
	case DefaultTerminal:
		script = fmt.Sprintf(terminalScript, name)
	default:
		script = fmt.Sprintf(terminalScript, name)
	}
	_ = ioutil.WriteFile(fn, []byte(script), 0755)
	cmd := exec.Command("osascript", fn)
	err = cmd.Start()
	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}
	return
}
