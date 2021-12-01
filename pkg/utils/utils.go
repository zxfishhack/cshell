//go:build darwin

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

func OpenSSH(name string) (err error) {
	fn := filepath.Join(os.TempDir(), "run")
	script := fmt.Sprintf(iTerm2Script, name)
	_ = ioutil.WriteFile(fn, []byte(script), 0755)
	cmd := exec.Command("osascript", fn)
	err = cmd.Start()
	if cmd.Process != nil {
		_ = cmd.Process.Release()
	}
	return
}
