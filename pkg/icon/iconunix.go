//go:build linux || darwin
// +build linux darwin

package icon

import _ "embed"

//go:embed icon.png
var Icon []byte
