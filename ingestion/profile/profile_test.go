package profile_test

import (
	"encoding/json"
	"testing"

	"github.com/wtask-go/mixpanel/ingestion/profile"
	"github.com/wtask-go/mixpanel/internal/assets"
)

func Test_json_schema(t *testing.T) {
	cases := []struct {
		data interface{}
		ok   bool
	}{
		{
			&profile.Set{}, false,
		},
		{
			&profile.Set{Set: map[string]interface{}{
				"key": "value",
			}}, true,
		},
		{
			&profile.SetOnce{}, false,
		},
		{
			&profile.SetOnce{SetOnce: map[string]interface{}{
				"key": "value",
			}}, true,
		},
		{
			&profile.NumberAdd{}, false,
		},
		{
			&profile.NumberAdd{Add: map[string]interface{}{
				"numeric_key": 1,
			}}, true,
		},
		{
			&profile.ListAppend{}, false,
		},
		{
			&profile.ListAppend{Append: map[string]interface{}{
				"list_key": "value",
			}}, true,
		},
		{
			&profile.ListRemove{}, false,
		},
		{
			&profile.ListRemove{Remove: map[string]interface{}{
				"list_key": "value",
			}}, true,
		},
		{
			&profile.Unset{}, false,
		},
		{
			&profile.Unset{Unset: []string{
				"numeric_key",
			}}, true,
		},
	}

	schema := assets.MustCompileSchema("openapi/engage.schema.json")

	for i, c := range cases {
		data, err := json.Marshal(c.data)
		if err != nil {
			t.Fatalf("[#%d] %s", i, err)
		}

		err = assets.ValidateJSON(schema, data)

		switch {
		case err != nil && c.ok:
			t.Fatalf("[#%d] unexpected error: %s, data: %+v, json: %s", i, err, c.data, data)
		case err == nil && !c.ok:
			t.Fatalf("[#%d] unexpectedly valid json %s", i, data)
		}
	}
}
