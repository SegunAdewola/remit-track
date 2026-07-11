# backend/ (Go modulith)

Single binary: cmd/api-server. The SQS consumer runs as a goroutine, not a second process
(docs/sadd.md Section B). recover() at the ingestion goroutine boundary so a scraper panic cannot kill the API.

- internal/ packages are the service seams: cache, database, models, alerting, ingestion, server.
  Dependencies point inward. models/ imports no third-party packages.
- Domain types (RateRecord) come from the shared module. Do not redefine them here.
- Handlers are thin: no business logic and no direct DB calls in handlers. Logic lives in internal packages.
- Tests: table-driven, colocated _test.go. Ask before adding testify or any test dependency.
- Run go test ./... and golangci-lint run before every commit.