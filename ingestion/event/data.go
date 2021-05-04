// Package event is intended to describe Mixpanel Ingestion API event.
package event

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	// Data describes action that take place within your product.
	// An event contains properties that describe the action.
	Data struct {
		// The name of the action to track
		Event string `json:"event"`

		// A dictionary of properties to hold metadata about your event
		Properties Properties `json:"properties"`
	}

	// Properties is a dictionary to hold metadata about your event.
	Properties struct {
		// A random 36 character string of hyphenated alphanumeric characters that is unique to an event.
		// Hyphen (-) is optional.
		// $insert_id can contain less than 36 characters,
		// but any string longer than 36 characters will be truncated.
		// If an $insert_id contains non-alphanumeric or non-hyphen characters
		// then Mixpanel replaces it with a random alphanumeric value.
		InsertID string `json:"$insert_id,omitempty"`

		// The value of distinct_id will be treated as a string,
		// and used to uniquely identify a user associated with your event.
		// If you provide a distinct_id property with your events,
		// you can track a given user through funnels and distinguish unique users for retention analyzes.
		// You should always send the same distinct_id when an event is triggered by the same user.
		DistinctID string `json:"distinct_id,omitempty"`

		// An IP address string (e.g. "127.0.0.1") associated with the event.
		// This is used for adding geolocation data to events,
		// and should only be required if you are making requests from your backend.
		// If ip is absent (and ip=1 is not provided as a URL parameter),
		// Mixpanel will ignore the IP address of the request.
		IP string `json:"ip,omitempty"`

		// The time that this event occurred.
		// If present, the value should be a unix timestamp (seconds since midnight, January 1st, 1970 - UTC).
		// If this property is not included in your request, Mixpanel will use the time the event
		// arrives at the server. If you're using our mobile SDKs, it will be set automatically for you.
		Time time.Time `json:"time,omitempty"`

		// The Mixpanel token associated with your project.
		// You can find your Mixpanel token in the project settings dialog in the Mixpanel app.
		// Events without a valid token will be ignored.
		Token string `json:"token"`

		// Additional properties to send to Mixpanale as part of event.
		CustomProperties `json:"-"`
	}

	// CustomProperties is extra event metadata.
	CustomProperties map[string]interface{}
)

// MarshalJSON implements json.Marshaler interface.
func (p *Properties) MarshalJSON() ([]byte, error) {
	if p == nil {
		return nil, nil
	}

	obj := map[string]json.RawMessage{}
	marshal := func(key string, value interface{}, include bool) (err error) {
		if !include {
			return nil
		}

		obj[key], err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("marshal %s: %w", key, err)
		}

		return nil
	}

	if err := marshal("$insert_id", p.InsertID, p.InsertID != ""); err != nil {
		return nil, err
	}

	if err := marshal("distinct_id", p.DistinctID, p.DistinctID != ""); err != nil {
		return nil, err
	}

	if err := marshal("ip", p.IP, p.IP != ""); err != nil {
		return nil, err
	}

	if err := marshal("time", p.Time.UTC().Unix(), !p.Time.IsZero()); err != nil {
		return nil, err
	}

	if err := marshal("token", p.Token, true); err != nil {
		return nil, err
	}

	for name, value := range p.CustomProperties {
		if err := marshal(name, value, true); err != nil {
			return nil, err
		}
	}

	return json.Marshal(obj)
}

// UnmarshalJSON implements json.Unmarshaler interface.
// nolint:funlen // it is required parse raw data
func (p *Properties) UnmarshalJSON(raw []byte) error {
	if p == nil {
		return fmt.Errorf("unmarshal event.Properties to nil")
	}

	data := map[string]json.RawMessage{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}

	unmarshal := func(key string, target interface{}) error {
		if v, ok := data[key]; ok {
			if err := json.Unmarshal(v, target); err != nil {
				return err
			}

			delete(data, key)
		}

		return nil
	}

	if err := unmarshal("$insert_id", &p.InsertID); err != nil {
		return err
	}

	if err := unmarshal("distinct_id", &p.DistinctID); err != nil {
		return err
	}

	if err := unmarshal("ip", &p.IP); err != nil {
		return err
	}

	var unix int64
	if err := unmarshal("time", &unix); err != nil {
		return err
	}

	if unix != 0 {
		p.Time = time.Unix(unix, 0).UTC()
	}

	if err := unmarshal("token", &p.Token); err != nil {
		return err
	}

	if len(data) > 0 {
		if p.CustomProperties == nil {
			p.CustomProperties = CustomProperties{}
		}

		for k, v := range data {
			var c interface{}
			if err := json.Unmarshal(v, &c); err != nil {
				return err
			}

			p.CustomProperties[k] = c
		}
	}

	return nil
}
