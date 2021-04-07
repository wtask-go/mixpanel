package event_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/wtask-go/mixpanel/ingestion/event"
)

func TestValidateJSON(t *testing.T) {
	cases := []struct {
		data  interface{}
		valid bool
	}{
		{(*event.Properties)(nil), false},
		{"event", false},
		{1, false},
		{true, false},
		{false, false},
		{map[string]interface{}{"event": "", "properties": nil}, false},
		{&struct{}{}, false},
		{struct{}{}, false},
		{
			struct {
				Event      string      `json:"event"`
				Properties interface{} `json:"properties"`
			}{"event", nil},
			false,
		},
		{
			struct {
				Event      string      `json:"event"`
				Properties interface{} `json:"properties"`
			}{"event", struct{}{}},
			false,
		},
		{
			struct {
				Event      string      `json:"event"`
				Properties interface{} `json:"properties"`
			}{"event", struct {
				Token string `json:"token"`
			}{""}},
			true, // all required fields are present
		},
		{
			map[string]interface{}{
				"event": "",
				"properties": map[string]interface{}{
					"token": "",
				}},
			true, // all require fields are present
		},
		{&event.Data{}, true}, // required are present (but are empty)
		{&event.Data{
			Properties: event.Properties{
				CustomProperties: event.CustomProperties{
					"user": "John Smith",
				},
			},
		}, true}, // required Event and Token are empty
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
	for i, c := range cases {
		data, err := json.Marshal(c.data)
		if err != nil {
			t.Fatalf("#%d: %s", i, err)
		}

		err = event.ValidateJSON(data)

		switch {
		case err != nil && c.valid:
			t.Fatalf("#%d: unexpected error: %s, data: %+v, json: %s", i, err, c.data, data)
		case err == nil && !c.valid:
			t.Fatalf("#%d: unexpectedly valid json %s", i, data)
		}
	}
}
