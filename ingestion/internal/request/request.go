package request

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// NewPostFormURLEncoded builds http request to post url-encoded form.
// Adds `Content-Type` and ``Content-Length  headers only.
func NewPostFormURLEncoded(ctx context.Context, url string, values *url.Values) (*http.Request, error) {
	body := strings.NewReader(values.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", body.Size()))

	return req, nil
}
