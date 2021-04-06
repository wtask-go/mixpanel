package assets_test

import (
	"testing"

	"github.com/wtask-go/mixpanel/assets"
)

func TestFS_embedded(t *testing.T) {
	embeds := map[string][]string{
		// Check embedded dirs.
		// If `key` describes entries all of them will be checked they were successfully embedded.
		// If not, dir will be checked is not empty.
		"openapi": {
			"event.schema.json",
			"ingestion.yml",
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
