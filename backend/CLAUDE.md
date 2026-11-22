# backend/ (Go modulith)

Single binary: cmd/api-server. The SQS consumer runs as a goroutine, not a second process
(docs/sadd.md Section B). recover() at the ingestion goroutine boundary so a scraper panic cannot kill the API.

Modules (docs/decisions/0002-domain-owned-modules.md): internal/rates, ingestion, alerting, accounts,
outbound, httpapi, plus technical internal/platform/* packages.
- Each module owns its tables in its own Postgres schema. Never read or join another module's tables.
- Modules never import each other. A module declares the small interfaces it needs; cmd/api-server is
  the only place that wires modules together. Only shared/ types (RateRecord, Money, CurrencyPair)
  cross module boundaries.
- Each module starts as one package. Only when it splits into sub-packages do they go under
  <module>/internal/. A depguard rule in .golangci.yml rejects cross-module imports.
- platform/ packages hold technical code only (no domain logic) and may be imported by any module.
- Domain types (RateRecord) come from the shared module. Do not redefine them here.
- Handlers are thin: no business logic and no direct DB calls in handlers. Logic lives in modules.
- Tests: table-driven, colocated _test.go; test a module with fakes for the interfaces it declares.
  Ask before adding testify or any test dependency.
- Run go test ./... and golangci-lint run before every commit.
