package ingestion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/wtask-go/mixpanel/ingestion/event"
	"github.com/wtask-go/mixpanel/ingestion/profile"
	"github.com/wtask-go/mixpanel/internal/form"
)

func (c *client) makeTrackRequest(data *event.Data) (*http.Request, error) {
	body, err := makeEventForm(data, form.WithVerboseResponse(true))
	if err != nil {
		return nil, err
	}

	return makeFormURLEncodedPost(c.endpoint.track.live.String(), body)
}

func (c *client) makeTrackDeduplicateRequest(data *event.Data) (*http.Request, error) {
	body, err := makeEventForm(data, form.WithVerboseResponse(true))
	if err != nil {
		return nil, err
	}

	return makeFormURLEncodedPost(c.endpoint.track.deduplicate.String(), body)
}

func (c *client) makeTrackBatchRequest(batch []*event.Data) (*http.Request, error) {
	switch l := len(batch); {
	case l == 0:
		return nil, fmt.Errorf("events batch is empty")
	case l > TrackBatchLimit:
		return nil, fmt.Errorf("events batch (%d) exceeds limit (%d)", l, TrackBatchLimit)
	}

	data, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}

	body, err := form.NewValues(data, form.WithVerboseResponse(true))
	if err != nil {
		return nil, err
	}

	return makeFormURLEncodedPost(c.endpoint.track.batch.String(), body)
}

func (c *client) makeEngageRequest(action profile.Mutator) (*http.Request, error) {
	var url string

	switch action.(type) {
	default:
		return nil, fmt.Errorf("unsupported engage action type %T", action)
	case nil:
		return nil, fmt.Errorf("engage action is nil")
	case *profile.Set:
		url = c.endpoint.engage.set.String()
	case *profile.SetOnce:
		url = c.endpoint.engage.setOnce.String()
	case *profile.NumberAdd:
		url = c.endpoint.engage.add.String()
	case *profile.ListAppend:
		url = c.endpoint.engage.append.String()
	case *profile.ListRemove:
		url = c.endpoint.engage.remove.String()
	case *profile.Unset:
		url = c.endpoint.engage.unset.String()
	}

	data, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}

	body, err := form.NewValues(data, form.WithVerboseResponse(true))
	if err != nil {
		return nil, err
	}

	return makeFormURLEncodedPost(url, body)
}

func (c *client) makeEngageBatchRequest(batch []profile.Mutator) (*http.Request, error) {
	if len(batch) == 0 {
		return nil, fmt.Errorf("empty profiles batch")
	}

	data, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}

	body, err := form.NewValues(data, form.WithVerboseResponse(true))
	if err != nil {
		return nil, err
	}

	return makeFormURLEncodedPost(c.endpoint.engage.batch.String(), body)
}

func makeEventForm(obj *event.Data, options ...form.OptionalValue) (*url.Values, error) {
	if obj == nil {
		return nil, fmt.Errorf("event object is nil")
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return form.NewValues(data, options...)
}

// makeFormURLEncodedPost builds http request to post url-encoded form.
// Adds `Content-Type` and `Content-Length  headers only.
func makeFormURLEncodedPost(url string, values *url.Values) (*http.Request, error) {
	body := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", body.Size()))

	return req, nil
}
