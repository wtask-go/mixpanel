package assets

import (
	_ "embed"
)

// Embedded asset content
var (
	// go:embed event.schema.json
	EventSchemaJSON []byte
)
