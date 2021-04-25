package event_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/wtask-go/mixpanel/ingestion/event"
)

func TestData_json(t *testing.T) {
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
