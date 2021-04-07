package event

import (
	"bytes"

	"github.com/santhosh-tekuri/jsonschema/v3"
	"github.com/wtask-go/mixpanel/assets"
)

// schema is compiled JSON schema to validate event data.
var schema = func() *jsonschema.Schema {
	name := "openapi/event.schema.json"
	f, err := assets.FS.Open(name)
	if err != nil {
		panic(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource(name, f); err != nil {
		panic(err)
	}

	return c.MustCompile(name)
}()

// ValidateJSON checks event data against Mixpanel JSON schema.
func ValidateJSON(json []byte) error {
	return schema.Validate(bytes.NewReader(json))
}
