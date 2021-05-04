// Package assets is intended to manage static content used by module.
package assets

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/santhosh-tekuri/jsonschema/v3"
)

// FS represents filesystem with embedded files
//go:embed openapi
var FS embed.FS

// MustCompileSchema builds json schema from specified asset.
func MustCompileSchema(asset string) *jsonschema.Schema {
	f, err := FS.Open(asset)
	if err != nil {
		panic(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource(asset, f); err != nil {
		panic(err)
	}

	return c.MustCompile(asset)
}

// MustCompileIngestionSpecification builds OpenAPI specification for Mixpanel Ingestion API.
func MustCompileIngestionSpecification() *openapi3.T {
	spec, err := FS.ReadFile("openapi/ingestion.openapi.yml")
	if err != nil {
		panic(err)
	}

	loader := openapi3.NewLoader()
	loader.Context = context.Background()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		switch url.String() {
		case "./event.schema.json":
			return FS.ReadFile("openapi/event.schema.json")
		case "./engage.schema.json":
			return FS.ReadFile("openapi/engage.schema.json")
		}

		return nil, fmt.Errorf("not found: %s", url.String())
	}

	oapi, err := loader.LoadFromData(spec)
	if err != nil {
		panic(err)
	}

	if err = oapi.Validate(loader.Context); err != nil {
		panic(err)
	}

	return oapi
}

// ValidateJSON is a helper func to validate JSON data with specified schema.
func ValidateJSON(schema *jsonschema.Schema, data []byte) error {
	return schema.Validate(bytes.NewReader(data))
}
