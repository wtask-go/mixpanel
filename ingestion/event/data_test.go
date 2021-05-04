package event_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/wtask-go/mixpanel/ingestion/event"
	"github.com/wtask-go/mixpanel/internal/assets"
)

func TestData_json_encoding(t *testing.T) {
	now := time.Now().Truncate(1 * time.Second).UTC()
	src := event.Data{
		Event: "test",
		Properties: event.Properties{
			InsertID:   "test-insert-id",
			DistinctID: "test-distinct-id",
			IP:         "ip-address",
			Time:       now,
			Token:      "test-token",
			CustomProperties: event.CustomProperties{
				"string":  "test-custom-string",
				"logical": true,
				"number":  3.14,
			},
		},
	}

	srcJSON, err := json.Marshal(&src)
	if err != nil {
		t.Fatal(err)
	}

	data := event.Data{}
	if err = json.Unmarshal(srcJSON, &data); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(data, src) {
		t.Fatalf("source: %+v, actual: %+v", src, data)
	}
}

func TestData_json_schema(t *testing.T) {
	cases := []struct {
		data  interface{}
		valid bool
	}{
		{&event.Data{}, true}, // required fields are present (have zero values)
		{&event.Data{
			Properties: event.Properties{
				CustomProperties: event.CustomProperties{
					"user": "John Smith",
				},
			},
		}, true}, // required Event and Properties.Token fields have zero values
		{&event.Data{
			Properties: event.Properties{
				IP:    "127:0:1:1",
				Token: "token-value",
			},
		}, false}, // invalid IP format
		{&event.Data{
			Properties: event.Properties{
				InsertID:   "uuid",
				Time:       time.Now(),
				DistinctID: "john@smith.tech",
				IP:         "127.0.1.1",
				Token:      "token-value",
				CustomProperties: event.CustomProperties{
					"username": "John Smith",
					"country":  "UK",
					"city":     "London",
					"age":      49,
				},
			},
		}, true},
	}

	schema := assets.MustCompileSchema("openapi/event.schema.json")

	for i, c := range cases {
		data, err := json.Marshal(c.data)
		if err != nil {
			t.Fatalf("#%d: %s", i, err)
		}

		err = assets.ValidateJSON(schema, data)

		switch {
		case err != nil && c.valid:
			t.Fatalf("[#%d] unexpected error: %s, data: %+v, json: %s", i, err, c.data, data)
		case err == nil && !c.valid:
			t.Fatalf("[#%d] unexpectedly valid json %s", i, data)
		}
	}
}
