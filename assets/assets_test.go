package assets_test

import (
	"testing"

	"github.com/wtask-go/mixpanel/assets"
)

func TestEventSchemaJSON_embedded(t *testing.T) {
	if len(assets.EventSchemaJSON) == 0 {
		t.Fail()
	}
}
