package assets_test

import (
	"testing"

	"github.com/wtask-go/mixpanel/internal/assets"
)

func TestFS_embedded(t *testing.T) {
	embeds := map[string][]string{
		// Check embedded dirs.
		// If `key` describes entries all of them will be checked they were successfully embedded.
		// If not, dir will be checked is not empty.
		"openapi": {
			"event.schema.json",
			"ingestion.openapi.yml",
			"engage.schema.json",
		},
	}
	for dir, files := range embeds {
		entries, err := assets.FS.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}

		switch {
		case len(files) > 0:
			if len(files) != len(entries) {
				t.Logf(
					"perhaps lame case, %q expected having %d entries, but actual %d",
					dir,
					len(files),
					len(entries),
				)
			}

			for _, filename := range files {
				_, err := assets.FS.ReadFile(dir + "/" + filename)
				if err != nil {
					t.Fatal(err)
				}
			}
		case len(entries) == 0:
			t.Fatalf("%q have no entries", dir)
		}
	}
}

func TestMustCompileSchema(t *testing.T) {
	schemas := []string{
		"openapi/event.schema.json",
		"openapi/engage.schema.json",
	}

	for _, asset := range schemas {
		schema := assets.MustCompileSchema(asset)
		if schema == nil {
			t.Fatal("required schema is nil", asset)
		}
	}
}

func TestMustCompileIngestionSpecification(t *testing.T) {
	spec := assets.MustCompileIngestionSpecification()
	if spec == nil {
		t.Fatal("required Ingestion OpenAPI specification is nil")
	}
}
