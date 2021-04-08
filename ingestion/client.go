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
	"github.com/wtask-go/mixpanel/ingestion/internal/request"
)

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
	httpc *http.Client
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
		Path: "/track#live-event",
	})
	cli.endpoint.deduplicate = server.ResolveReference(&url.URL{
		Path: "/track#live-event-deduplicate",
	})
	cli.endpoint.batch = server.ResolveReference(&url.URL{
		Path: "/track#past-events-batch",
	})
	cli.httpc = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	cli.agent = fmt.Sprintf(
		"ingestion.Client/v*.*.* (%s; %s; %s;)", runtime.GOOS, runtime.GOARCH, runtime.Version(),
	)

	for _, option := range options {
		if err := option(cli); err != nil {
			return nil, fmt.Errorf("client option: %w", err)
		}
	}

	return cli, nil
}

// WithHTTPClient allows to change default internal HTTP client with specified one.
func WithHTTPClient(httpc *http.Client) ClientOption {
	return func(c *client) error {
		if httpc == nil {
			return fmt.Errorf("%w: http client is nil", errs.ErrInvalidArgument)
		}

		c.httpc = httpc

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
	body, err := request.NewValues(
		data,
		request.WithVerboseResponse(true),
	)
	if err != nil {
		return err
	}

	req, err := request.NewPostFormURLEncoded(ctx, c.endpoint.live.String(), body)
	if err != nil {
		return err
	}

	c.addDefaultHeaders(req)

	resp, err := c.httpc.Do(req)
	if err != nil {
		return err
	}

	return c.parseTrackResponse(resp)
}

func (c *client) TrackDeduplicate(ctx context.Context, data *event.Data) error {
	body, err := request.NewValues(
		data,
		request.WithVerboseResponse(true),
	)
	if err != nil {
		return err
	}

	req, err := request.NewPostFormURLEncoded(ctx, c.endpoint.deduplicate.String(), body)
	if err != nil {
		return err
	}

	c.addDefaultHeaders(req)

	resp, err := c.httpc.Do(req)
	if err != nil {
		return err
	}

	return c.parseTrackResponse(resp)
}

func (c *client) TrackBatch(ctx context.Context, data []event.Data) error {
	body, err := request.NewBatchValues(
		data,
		request.WithVerboseResponse(true),
	)
	if err != nil {
		return err
	}

	req, err := request.NewPostFormURLEncoded(ctx, c.endpoint.deduplicate.String(), body)
	if err != nil {
		return err
	}

	c.addDefaultHeaders(req)

	resp, err := c.httpc.Do(req)
	if err != nil {
		return err
	}

	return c.parseTrackResponse(resp)
}

func (c *client) addDefaultHeaders(req *http.Request) {
	req.Header.Add("Accept", "plain/text")
	req.Header.Add("Accept", "application/json")

	if c.agent != "" {
		req.Header.Set("User-Agent", c.agent)
	}
}

func (c *client) parseTrackResponse(resp *http.Response) error {
	return nil
}
