package ingestion_test

import (
	"context"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/legacy"
	"github.com/wtask-go/mixpanel/ingestion"
	"github.com/wtask-go/mixpanel/ingestion/event"
)

func init() {
	// this formats are not automatically registered
	openapi3.DefineIPv4Format()
	openapi3.DefineIPv6Format()
}

func Test_client_openapi_requests(t *testing.T) {
	openapi3filter.RegisterBodyDecoder("application/x-www-form-urlencoded", urlencodedBodyDecoder)

	router, err := legacy.NewRouter(spec)
	if err != nil {
		t.Fatal(err)
	}

	cli, err := ingestion.NewClient(
		"https://api-eu.mixpanel.com",
		ingestion.WithHTTPDoer(OpenAPIRequestValidator(t, router)),
	)
	if err != nil {
		t.Fatal(err)
	}

	_ = cli.Track(context.Background(), &event.Data{
		Event: "live-1",
		Properties: event.Properties{
			InsertID:   "uuid",
			Token:      "token",
			DistinctID: "account@server.com",
			IP:         "127.0.0.1",
			Time:       time.Now(),
			CustomProperties: event.CustomProperties{
				"role": "manager",
			},
		},
	})

	_ = cli.TrackDeduplicate(context.Background(), &event.Data{
		Event: "deduplicate-1",
		Properties: event.Properties{
			InsertID:   "uuid",
			Token:      "token",
			DistinctID: "account@server.com",
			IP:         "127.0.0.1",
			Time:       time.Now(),
			CustomProperties: event.CustomProperties{
				"role": "manager",
			},
		},
	})
}
