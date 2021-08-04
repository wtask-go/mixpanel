# Mixpanel

Mixpanel API client for Go

## Requirements

* Go 1.16 and above. This module uses `embed` package introduced for Go 1.16

## Dependencies

* [github.com/santhosh-tekuri/jsonschema](https://github.com/santhosh-tekuri/jsonschema) -- for Ingestion API Event validation
* [github.com/getkin/kin-openapi](https://github.com/getkin/kin-openapi) -- for requests validation against OpenAPI v3 specification

## Ingestion API

Check official docs for [Ingestion API](https://developer.mixpanel.com/reference/ingestion-api) details.

Our client offers top-level interface to interact with Mixpanel endpoints.
We use semi-official json schema of Event object in tests to validate prepared event data. Check [this page in docs](https://developer.mixpanel.com/docs/data-model#anatomy-of-an-event) for the [schema link](https://gist.github.com/jbwyme/f01f0a6f6f8b8db2472cb8771f7a505c).

Also we made own [OpenAPI schema](./internal/assets/assets/ingestion.openapi.yml) to describe external Mixpanel Ingestion API. The module uses mentioned schema to validate prepared HTTP requests in tests only.

### Events

* Track Event: `ingestion.Client.Track()`
* Track Event with Deduplication: `ingestion.Client.TrackDeduplicate()`
* Track Multiple Events: `ingestion.Client.TrackBatch()`

### User Profiles

* Set Property, Set Property Once, Increment Numerical Property, Append to List Property, Remove from List Property, Delete Property: `ingestion.Client.Engage()`
* Update Multiple Profiles: `ingestion.Client.EngageBatch()`
