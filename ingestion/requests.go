package ingestion

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/wtask-go/mixpanel/ingestion/event"
	"github.com/wtask-go/mixpanel/internal/form"
)

func (c *client) makeTrackRequest(data *event.Data) (*http.Request, error) {
	body, err := form.NewValues(
		data,
		form.WithVerboseResponse(true),
	)
	if err != nil {
		return nil, err
	}

	req, err := makeFormURLEncodedPost(c.endpoint.live.String(), body)
	if err != nil {
		return nil, err
	}

	c.addDefaultHeaders(req)

	return req, nil
}

func (c *client) makeTrackDeduplicateRequest(data *event.Data) (*http.Request, error) {
	body, err := form.NewValues(
		data,
		form.WithVerboseResponse(true),
	)
	if err != nil {
		return nil, err
	}

	req, err := makeFormURLEncodedPost(c.endpoint.deduplicate.String(), body)
	if err != nil {
		return nil, err
	}

	c.addDefaultHeaders(req)

	return req, nil
}

func (c *client) makeBatchRequest(data []event.Data) (*http.Request, error) {
	body, err := form.NewBatchValues(
		data,
		form.WithVerboseResponse(true),
	)
	if err != nil {
		return nil, err
	}

	req, err := makeFormURLEncodedPost(c.endpoint.batch.String(), body)
	if err != nil {
		return nil, err
	}

	c.addDefaultHeaders(req)

	return req, nil
}

// addDefaultHeaders adds default request headers.
func (c *client) addDefaultHeaders(req *http.Request) {
	req.Header.Add("Accept", "plain/text")
	req.Header.Add("Accept", "application/json")

	if c.agent != "" {
		req.Header.Set("User-Agent", c.agent)
	}
}

// makeFormURLEncodedPost builds http request to post url-encoded form.
// Adds `Content-Type` and `Content-Length  headers only.
func makeFormURLEncodedPost(url string, values *url.Values) (*http.Request, error) {
	body := strings.NewReader(values.Encode())

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", body.Size()))

	return req, nil
}
