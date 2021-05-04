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
	"github.com/wtask-go/mixpanel/ingestion/profile"
	"github.com/wtask-go/mixpanel/internal/assets"
)

func init() {
	// this formats are not automatically registered
	openapi3.DefineIPv4Format()
	openapi3.DefineIPv6Format()
}

func Test_Client_openapi_requests(t *testing.T) {
	openapi3filter.RegisterBodyDecoder("application/x-www-form-urlencoded", urlencodedBodyDecoder)

	router, err := legacy.NewRouter(
		assets.MustCompileIngestionSpecification(),
	)
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

	_ = cli.TrackBatch(context.Background(), []*event.Data{
		{
			Event: "outdated-1",
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
		},
	})

	_ = cli.Engage(context.Background(), &profile.Set{
		Token:      "token",
		DistinctID: "user-id",
		Set: map[string]interface{}{
			"counter": 0,
		},
	})

	_ = cli.Engage(context.Background(), &profile.SetOnce{
		Token:      "token",
		DistinctID: "user-id",
		SetOnce: map[string]interface{}{
			"verified": true,
		},
	})

	_ = cli.Engage(context.Background(), &profile.NumberAdd{
		Token:      "token",
		DistinctID: "user-id",
		Add: map[string]interface{}{
			"counter": 1,
		},
	})

	_ = cli.Engage(context.Background(), &profile.ListAppend{
		Token:      "token",
		DistinctID: "user-id",
		Append: map[string]interface{}{
			"roles": "user",
		},
	})

	_ = cli.Engage(context.Background(), &profile.ListRemove{
		Token:      "token",
		DistinctID: "user-id",
		Remove: map[string]interface{}{
			"roles": "manager",
		},
	})

	_ = cli.Engage(context.Background(), &profile.Unset{
		Token:      "token",
		DistinctID: "user-id",
		Unset: []string{
			"counter",
		},
	})

	_ = cli.EngageBatch(context.Background(), []profile.Mutator{
		&profile.Set{
			Token:      "token",
			DistinctID: "user-id",
			Set: map[string]interface{}{
				"counter": 0,
			},
		},
		&profile.SetOnce{
			Token:      "token",
			DistinctID: "user-id",
			SetOnce: map[string]interface{}{
				"verified": true,
			},
		},
		&profile.NumberAdd{
			Token:      "token",
			DistinctID: "user-id",
			Add: map[string]interface{}{
				"counter": 1,
			},
		},
		&profile.ListAppend{
			Token:      "token",
			DistinctID: "user-id",
			Append: map[string]interface{}{
				"roles": "user",
			},
		},
		&profile.ListRemove{
			Token:      "token",
			DistinctID: "user-id",
			Remove: map[string]interface{}{
				"roles": "manager",
			},
		},
		&profile.Unset{
			Token:      "token",
			DistinctID: "user-id",
			Unset: []string{
				"counter",
			},
		},
	})
}
