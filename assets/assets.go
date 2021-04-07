// Package assets is intended to manage static content used by module.
package assets

import (
	"embed"
)

// FS represents filesystem with embedded files
//go:embed openapi
var FS embed.FS
