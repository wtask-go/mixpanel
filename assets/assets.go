// Package assets is intended to manage static content used by module.
package assets

import (
	"embed"
)

var (
	// FS represents filesystem with embedded files
	//go:embed openapi
	FS embed.FS
)
