package form

import (
	"fmt"
	"net/url"
)

// OptionalValue is intended to pass not mandatory request values.
type OptionalValue func(*url.Values)

// NewValues builds url.Values with required `data` item.
func NewValues(data []byte, optional ...OptionalValue) (*url.Values, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}

	values := &url.Values{}
	values.Set("data", string(data))

	for _, value := range optional {
		value(values)
	}

	return values, nil
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
// Option works like OneOf between redirect, image, callback.
func WithRedirectResponse(redirect string) OptionalValue {
	return func(values *url.Values) {
		if redirect != "" {
			values.Set("redirect", redirect)
			values.Del("image")
			values.Del("callback")
		} else {
			values.Del("redirect")
		}
	}
}

// WithImageResponse builds option to pass image flag.
// As result Mixpanel will serve 1x1 transparent image as a response to the request.
// Expected content type is `image/png`.
// Option works like OneOf between redirect, image, callback.
func WithImageResponse(image bool) OptionalValue {
	return func(values *url.Values) {
		if image {
			values.Set("image", "1")
			values.Del("redirect")
			values.Del("callback")
		} else {
			values.Del("image")
		}
	}
}

// WithJavascriptCallback builds option to pass javascript callback value.
// As result Mixpanel will return a content-type: `text/javascript`
// with a body that calls a function by value provided.
// Option works like OneOf between redirect, image, callback.
func WithJavascriptCallback(callback string) OptionalValue {
	return func(values *url.Values) {
		if callback != "" {
			values.Set("callback", callback)
			values.Del("redirect")
			values.Del("image")
		} else {
			values.Del("callback")
		}
	}
}
