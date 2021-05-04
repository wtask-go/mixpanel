package ingestion_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
)

// urlencodedBodyDecoder is required for `kin-openapi` to decode custom urlencoded form.
func urlencodedBodyDecoder(
	body io.Reader,
	_ http.Header,
	schema *openapi3.SchemaRef,
	_ openapi3filter.EncodingFn,
) (interface{}, error) {
	// Validate schema of request body.
	// By the OpenAPI 3 specification request body's schema must have type "object".
	// Properties of the schema describes individual parts of request body.
	if schema.Value.Type != "object" {
		return nil, errors.New("unsupported schema of request body")
	}

	// Parse form.
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}

	// check schema for required fields
	required := map[string]bool{}
	for _, name := range schema.Value.Required {
		required[name] = true
	}

	// Make an object from form body.
	obj := make(map[string]interface{})

	for name, prop := range schema.Value.Properties {
		_, ok := values[name]
		if !ok {
			if required[name] {
				return nil, fmt.Errorf("required field %s is missing", name)
			}

			continue
		}

		switch t := prop.Value.Type; t {
		case "object":
			raw := map[string]interface{}{}
			if err := json.Unmarshal([]byte(values.Get(name)), &raw); err != nil {
				return err, fmt.Errorf("invalid object %s: %w", name, err)
			}

			obj[name] = raw
		case "array":
			raw := []interface{}{}
			if err := json.Unmarshal([]byte(values.Get(name)), &raw); err != nil {
				return err, fmt.Errorf("invalid array %s: %w", name, err)
			}

			obj[name] = raw
		case "integer":
			v, err := strconv.ParseFloat(values.Get(name), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %s (%s): %w", name, values.Get(name), err)
			}

			obj[name] = v
		case "string":
			obj[name] = values.Get(name)
		default:
			return nil, fmt.Errorf("unrecognized field type %q", t)
		}
	}

	return obj, nil
}

type HTTPDoerMock func(*http.Request) (*http.Response, error)

func (mock HTTPDoerMock) Do(req *http.Request) (*http.Response, error) {
	return mock(req)
}

func OpenAPIRequestValidator(t *testing.T, router routers.Router) HTTPDoerMock {
	t.Helper()

	return func(req *http.Request) (*http.Response, error) {
		// router is helper to traverse compiled OpenAPI document
		// to find schema corresponded to request
		route, params, err := router.FindRoute(req)
		if err != nil {
			t.Fatalf("OpenAPI router failed: %s", err)
		}

		err = openapi3filter.ValidateRequest(
			context.Background(),
			&openapi3filter.RequestValidationInput{
				Request:    req,
				PathParams: params,
				Route:      route,
			},
		)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(req.URL.String(), "valid")

		return ResponseOK("1", req), nil
	}
}

func ResponseOK(body string, req *http.Request) *http.Response {
	return &http.Response{
		Request:    req,
		Header:     http.Header{},
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(body)),
	}
}
