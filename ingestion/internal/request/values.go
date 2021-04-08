package request

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/wtask-go/mixpanel/errs"
	"github.com/wtask-go/mixpanel/ingestion/event"
)

// MaxBatchLength is Mixpanel limitation to track events batch.
const MaxBatchLength = 50

// OptionalValue is intended to pass not mandatory values.
type OptionalValue func(*url.Values)

func makeValues(data string, optional ...OptionalValue) (*url.Values, error) {
	if data == "" {
		return nil, fmt.Errorf("%w: data is empty", errs.ErrInvalidArgument)
	}

	values := &url.Values{}
	values.Set("data", data)

	for _, value := range optional {
		value(values)
	}

	return values, nil
}

// NewValues builds request values required to track single event.
func NewValues(data *event.Data, optional ...OptionalValue) (*url.Values, error) {
	if data == nil {
		return nil, fmt.Errorf("%w: event data is nil", errs.ErrInvalidArgument)
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("event data marshaling: %w", err)
	}

	return makeValues(string(encoded), optional...)
}

// NewBatchValues builds request values required to track batch events.
func NewBatchValues(data []event.Data, optional ...OptionalValue) (*url.Values, error) {
	switch l := len(data); {
	case l == 0:
		return nil, fmt.Errorf("%w: events batch empty", errs.ErrInvalidArgument)
	case l > MaxBatchLength:
		return nil, fmt.Errorf("%w: events batch exceeds %d items ", errs.ErrInvalidArgument, MaxBatchLength)
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("events batch marshaling: %w", err)
	}

	optional = append(
		optional,
		// are not supported for batch
		WithIPAsDistinctID(false),
		WithRedirectResponse(""),
		WithImageResponse(false),
		WithJavascriptCallback(""),
	)

	return makeValues(string(encoded), optional...)
}

// WithVerboseResponse builds option to pass verbose flag.
// As result Mixpanel should response with `application/json` content-type and JSON.
func WithVerboseResponse(verbose bool) OptionalValue {
	return func(values *url.Values) {
		if verbose {
			values.Set("verbose", "1")
		} else {
			values.Del("verbose")
		}
	}
}

// WithIPAsDistinctID builds option to pass ip flag.
// As result Mixpanel will use IP-address of incoming request as distinct_id if none is provided for event data.
func WithIPAsDistinctID(ip bool) OptionalValue {
	return func(values *url.Values) {
		if ip {
			values.Set("ip", "1")
		} else {
			values.Del("ip")
		}
	}
}

// WithRedirectResponse builds option to pass redirection URL.
// As result Mixpanel will serve redirect as a response to the request.
func WithRedirectResponse(redirect string) OptionalValue {
	return func(values *url.Values) {
		if redirect == "" {
			values.Del("redirect")
		} else {
			values.Set("redirect", redirect)
		}
	}
}

// WithImageResponse builds option to pass image flag.
// As result Mixpanel will serve 1x1 transparent image as a response to the request.
// Expected content type is `image/png`.
func WithImageResponse(image bool) OptionalValue {
	return func(values *url.Values) {
		if image {
			values.Set("image", "1")
		} else {
			values.Del("image")
		}
	}
}

// WithJavascriptCallback builds option to pass javascript callback value.
// As result Mixpanel will return a content-type: `text/javascript`
// with a body that calls a function by value provided.
func WithJavascriptCallback(callback string) OptionalValue {
	return func(values *url.Values) {
		if callback != "" {
			values.Set("callback", callback)
		} else {
			values.Del("callback")
		}
	}
}
