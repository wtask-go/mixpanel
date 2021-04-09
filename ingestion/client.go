// Package ingestion contains client to interact with Mixpanel Ingestion API.
package ingestion

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/wtask-go/mixpanel/errs"
	"github.com/wtask-go/mixpanel/ingestion/event"
)

// HTTPDoer represent HTTP client interface only required for this package.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client represents top-level interface of Mixpanel Ingestion API
type Client interface {
	Track(context.Context, *event.Data) error
	TrackDeduplicate(context.Context, *event.Data) error
	TrackBatch(context.Context, []event.Data) error
}

type client struct {
	endpoint struct {
		live        *url.URL
		deduplicate *url.URL
		batch       *url.URL
	}
	httpc HTTPDoer
	agent string
}

// ClientOption provides customization for Ingestion API client.
type ClientOption func(*client) error

// NewClient builds default implementation of Ingestion API client.
func NewClient(serverURL string, options ...ClientOption) (Client, error) {
	serverURL = strings.TrimRight(serverURL, "/")
	if serverURL == "" {
		return nil, fmt.Errorf("%w: insufficient or empty server URL", errs.ErrInvalidArgument)
	}

	server, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("parse server URL: %w", err)
	}

	cli := &client{}
	cli.endpoint.live = server.ResolveReference(&url.URL{
		Path:     "/track",
		Fragment: "live-event",
	})
	cli.endpoint.deduplicate = server.ResolveReference(&url.URL{
		Path:     "/track",
		Fragment: "live-event-deduplicate",
	})
	cli.endpoint.batch = server.ResolveReference(&url.URL{
		Path:     "/track",
		Fragment: "past-events-batch",
	})
	cli.httpc = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	cli.agent = fmt.Sprintf(
		"ingestion.Client/v* (%s; %s; %s;)", runtime.GOOS, runtime.GOARCH, runtime.Version(),
	)

	for _, option := range options {
		if err := option(cli); err != nil {
			return nil, fmt.Errorf("client option: %w", err)
		}
	}

	return cli, nil
}

// WithHTTPDoer allows to change default internal HTTP client with specified one.
func WithHTTPDoer(doer HTTPDoer) ClientOption {
	return func(c *client) error {
		if doer == nil {
			return fmt.Errorf("%w: http doer is nil", errs.ErrInvalidArgument)
		}

		c.httpc = doer

		return nil
	}
}

// WithUserAgent allows to set desired User-Agent header value.
func WithUserAgent(agent string) ClientOption {
	return func(c *client) error {
		c.agent = agent

		return nil
	}
}

func (c *client) Track(ctx context.Context, data *event.Data) error {
	req, err := c.makeTrackRequest(data)
	if err != nil {
		return err
	}

	resp, err := c.httpc.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	return c.parseResponse(resp)
}

func (c *client) TrackDeduplicate(ctx context.Context, data *event.Data) error {
	req, err := c.makeTrackDeduplicateRequest(data)
	if err != nil {
		return err
	}

	resp, err := c.httpc.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	return c.parseResponse(resp)
}

func (c *client) TrackBatch(ctx context.Context, data []event.Data) error {
	req, err := c.makeBatchRequest(data)
	if err != nil {
		return err
	}

	resp, err := c.httpc.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	return c.parseResponse(resp)
}
