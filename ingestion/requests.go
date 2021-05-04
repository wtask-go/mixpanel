package ingestion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/wtask-go/mixpanel/ingestion/event"
	"github.com/wtask-go/mixpanel/internal/form"
)

func (c *client) makeTrackRequest(data *event.Data) (*http.Request, error) {
	body, err := makeEventForm(data, form.WithVerboseResponse(true))
	if err != nil {
		return nil, err
	}

	return makeFormURLEncodedPost(c.endpoint.live.String(), body)
}

func (c *client) makeTrackDeduplicateRequest(data *event.Data) (*http.Request, error) {
	body, err := makeEventForm(data, form.WithVerboseResponse(true))
	if err != nil {
		return nil, err
	}

	return makeFormURLEncodedPost(c.endpoint.deduplicate.String(), body)
}

func (c *client) makeBatchRequest(data []*event.Data) (*http.Request, error) {
	body, err := makeEventBatchForm(data, form.WithVerboseResponse(true))
	if err != nil {
		return nil, err
	}

	return makeFormURLEncodedPost(c.endpoint.batch.String(), body)
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

func makeEventBatchForm(batch []*event.Data, options ...form.OptionalValue) (*url.Values, error) {
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
