# mixpanel

MixPanel API clients for Go

## Requirements

* Go 1.16 and above. Module uses `embed` package introduced for Go 1.16

## Dependencies

* github.com/santhosh-tekuri/jsonschema/v3 -- for Ingestion API Event validation
* github.com/getkin/kin-openapi -- for requests validation against OpenAPI v3 specification

## Ingestion API

Check official docs for [Ingestion API](https://developer.mixpanel.com/reference/ingestion-api) details.

Our client offers top-level interface to interact with Mixpanel endpoints.
We use semi-official json schema of Event object in tests to validate prepared event data. Check [this page in docs](https://developer.mixpanel.com/docs/data-model#anatomy-of-an-event) for the (schema link)[https://gist.github.com/jbwyme/f01f0a6f6f8b8db2472cb8771f7a505c].

Also we made own [OpenAPI schema](./internal/assets/ingestion.openapi.yml) to describe external Mixpanel Ingestion API. The module uses mentioned schema to validate prepared HTTP requests in tests only.
