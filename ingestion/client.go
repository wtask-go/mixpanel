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

	"github.com/wtask-go/mixpanel/ingestion/event"
	"github.com/wtask-go/mixpanel/ingestion/profile"
)

// HTTPDoer represent HTTP client interface only required for this package.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client represents top-level interface of Mixpanel Ingestion API
type Client interface {
	Track(context.Context, *event.Data) error
	TrackDeduplicate(context.Context, *event.Data) error
	TrackBatch(context.Context, []*event.Data) error
	Engage(context.Context, profile.Mutator) error
	EngageBatch(context.Context, []profile.Mutator) error
}

// TrackBatchLimit is Mixpanel limitation for events batch.
const TrackBatchLimit = 50

type client struct {
	endpoint struct {
		track struct {
			live        *url.URL
			deduplicate *url.URL
			batch       *url.URL
		}
		engage struct {
			set     *url.URL
			setOnce *url.URL
			add     *url.URL
			append  *url.URL
			remove  *url.URL
			unset   *url.URL
			batch   *url.URL
		}
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
		return nil, fmt.Errorf("bad or empty server URL")
	}

	server, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("parse server URL: %w", err)
	}

	serverRef := func(path, fragment string) *url.URL {
		return server.ResolveReference(&url.URL{Path: path, Fragment: fragment})
	}

	cli := &client{}
	cli.endpoint.track.live = serverRef("/track", "live-event")
	cli.endpoint.track.deduplicate = serverRef("/track", "live-event-deduplicate")
	cli.endpoint.track.batch = serverRef("/track", "past-events-batch")

	cli.endpoint.engage.set = serverRef("/engage", "profile-set")
	cli.endpoint.engage.setOnce = serverRef("/engage", "profile-set-once")
	cli.endpoint.engage.add = serverRef("/engage", "profile-numerical-add")
	cli.endpoint.engage.append = serverRef("/engage", "profile-list-append")
	cli.endpoint.engage.remove = serverRef("/engage", "profile-list-remove")
	cli.endpoint.engage.unset = serverRef("/engage", "profile-unset")
	cli.endpoint.engage.batch = serverRef("/engage", "profile-batch-update")

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
			return fmt.Errorf("HTTPDoer is nil")
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

func (c *client) send(ctx context.Context, req *http.Request) error {
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "plain/text")
	req.Header.Add("Accept", "application/json")

	if c.agent != "" {
		req.Header.Set("User-Agent", c.agent)
	}

	resp, err := c.httpc.Do(req)
	if err != nil {
		return err
	}

	return c.parseResponse(resp)
}

func (c *client) Track(ctx context.Context, data *event.Data) error {
	req, err := c.makeTrackRequest(data)
	if err != nil {
		return err
	}

	return c.send(ctx, req)
}

func (c *client) TrackDeduplicate(ctx context.Context, data *event.Data) error {
	req, err := c.makeTrackDeduplicateRequest(data)
	if err != nil {
		return err
	}

	return c.send(ctx, req)
}

func (c *client) TrackBatch(ctx context.Context, data []*event.Data) error {
	req, err := c.makeTrackBatchRequest(data)
	if err != nil {
		return err
	}

	return c.send(ctx, req)
}

func (c *client) Engage(ctx context.Context, action profile.Mutator) error {
	req, err := c.makeEngageRequest(action)
	if err != nil {
		return err
	}

	return c.send(ctx, req)
}

func (c *client) EngageBatch(ctx context.Context, batch []profile.Mutator) error {
	req, err := c.makeEngageBatchRequest(batch)
	if err != nil {
		return nil
	}

	return c.send(ctx, req)
}
