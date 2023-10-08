//go:build windows
// +build windows

package icon

import _ "embed"

//go:embed icon.ico
var Icon []byte
