// Package assets is intended to manage static content used by module.
package assets

import (
	_ "embed" // required to be anonymous until embed.FS will be used
)

// Embedded asset content
var (
	//go:embed event.schema.json
	EventSchemaJSON []byte
)
