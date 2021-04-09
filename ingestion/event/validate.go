package event

import (
	"bytes"

	"github.com/santhosh-tekuri/jsonschema/v3"
	"github.com/wtask-go/mixpanel/internal/assets"
)

// schema is compiled JSON schema to validate event data.
var schema = compileSchema("openapi/event.schema.json")

// compileSchema compiles schema from specified asset source.
func compileSchema(asset string) *jsonschema.Schema {
	f, err := assets.FS.Open(asset)
	if err != nil {
		panic(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource(asset, f); err != nil {
		panic(err)
	}

	return c.MustCompile(asset)
}

// ValidateJSON checks event data against Mixpanel event JSON schema.
func ValidateJSON(json []byte) error {
	return schema.Validate(bytes.NewReader(json))
}
