# RemitTrack Implementation Plan

The definitive build order. Follow it top to bottom. Each numbered item is one atomic commit.

## How to use this
- Work one step at a time, in order. A step is not done until its test is green and it is committed.
- Declare a learning mode at the start of each step (default A). Core logic steps should be Mode A.
- TDD always: write the failing test named in the step first, watch it fail, then implement.
- One step = one commit, using the message shown. Reference the step in the commit body (e.g. "step 4.2").
- Do not start a step whose dependencies (earlier steps) are not green.

## Testing ladder (what "tested" means at each level)
- Unit: pure logic in a package, no I/O. Table-driven in Go.
- Integration (seam): the boundary between components, using real Postgres and a mock SQS via
  tests/regression/docker-compose.yml. Added the moment a seam exists.
- End-to-end: a full vertical slice (scraper -> SQS -> consumer -> DB, or HTTP request -> DB -> response).
- Regression: replays of historical failures; every case in failure_registry.json stays green.

## When surprises happen (read before deviating)
The plan is a living document, not a contract. If reality diverges (an assumption breaks, a better
approach appears, a library misbehaves):
1. STOP. Do not silently deviate. Silent drift makes this file a lie.
2. Discuss and decide with Claude in plan mode.
3. Update the affected steps here, and append a dated entry to the Decisions & changes log below.
4. If it is architecture-level, write an ADR in docs/decisions/ and update docs/sadd.md.
5. Commit the doc update as its own `docs:` commit BEFORE resuming code.

---

## Phase 0 - Bootstrap
- [x] 0.1 chore(build): confirm go.work ties shared, backend, lambdas; go mod init each. Test: `go build ./...` in each module.
- [ ] 0.2 chore(lint): add .golangci.yml (govet, staticcheck, errcheck, revive, gofumpt) using the v2 schema. Test: `golangci-lint config verify` passes. (A full `golangci-lint run` returns exit 5 on a module with no .go files; it becomes green once the first package lands at 1.1.)
- [ ] 0.3 chore(build): Makefile (test, lint, build, up/down targets) + dev docker-compose.yml (Postgres for Phase 2 integration tests). The test/lint/build targets skip modules with no .go files so they are clean on the empty workspace. Test: `make test` runs.
- [ ] 0.4 ci: backend-ci.yml runs go test + golangci-lint on PR (via `make lint`, which skips empty modules). Test: workflow parses (act or push).

## Phase 1 - Shared domain (foundation, no dependencies)
- [ ] 1.1 feat(shared): corridor CurrencyPair type with Parse/String ("USD:GHS"). Test: parse valid/invalid, round-trip.
- [ ] 1.2 feat(shared): RateRecord struct + Validate (positive rate, known pair, non-zero timestamp). Test: valid and each invalid case.
- [ ] 1.3 feat(shared): RateRecord JSON marshal/unmarshal with stable field names. Test: round-trip + field-name assertions.

## Phase 2 - Persistence (the source of truth)
- [ ] 2.1 feat(backend): internal/models structs (Quote view, AlertConfig, User) + validation. Test: zero-value and validation (unit).
- [ ] 2.2 feat(database): migration 0001 historical_rates + current_rates + UNIQUE(provider,currency_pair,scrape_timestamp). Test: up/down applies (integration, real Postgres).
- [ ] 2.3 feat(database): connection pool + config from env. Test: connect and ping (integration).
- [ ] 2.4 feat(database): idempotent UPSERT into historical_rates (ON CONFLICT DO NOTHING). Test: duplicate ignored, distinct timestamps kept (integration).
- [ ] 2.5 feat(database): current_rates upsert on insert + read-latest query. Test: one row per pair, historical not scanned (integration).

## Phase 3 - Cache
- [ ] 3.1 feat(cache): RateCache interface + sync.Map implementation. Test: get/set/miss + concurrent access (unit).

## Phase 4 - Ingestion (SQS consumer, in-process)
- [ ] 4.1 feat(ingestion): batch builder accumulates messages into multi-row batches by size and time. Test: batching rules (unit).
- [ ] 4.2 feat(ingestion): consumer long-polls a queue, decodes RateRecord, writes via database. Test: mock SQS -> rows in Postgres (integration).
- [ ] 4.3 feat(ingestion): recover() boundary + poison-message handling (log/DLQ, loop survives). Test: bad message does not kill the loop (integration).

## Phase 5 - Alerting (in-memory, coupled to ingestion)
- [ ] 5.1 feat(alerting): corridor index map build/refresh from AlertConfig. Test: build and lookup (unit).
- [ ] 5.2 feat(alerting): evaluation loop fires on threshold crossing only. Test: crossing fires, non-crossing does not (unit).
- [ ] 5.3 feat(alerting): notification dispatch interface (SES behind iface, fake in test). Test: fired alert dispatches once (unit).
- [ ] 5.4 feat(backend): wire ingestion batch completion to alerting evaluation in-process. Test: batch write triggers eval (integration).

## Phase 6 - Server (read path first, then edges)
- [ ] 6.1 feat(server): current-rates-by-corridor handler backed by cache/DB. Test: 200 + payload (httptest integration).
- [ ] 6.2 feat(server): logging + panic-recovery middleware. Test: handler panic returns 500, server survives (unit).
- [ ] 6.3 feat(server): token-bucket ratelimit middleware. Test: over-limit returns 429 (unit).
- [ ] 6.4 feat(server): B2B token auth middleware for /api/v1/b2b/*. Test: valid vs invalid key (integration).
- [ ] 6.5 feat(server): Stripe webhook idempotency via evt_id ledger. Test: duplicate event ignored (integration).
- [ ] 6.6 feat(server): masked /out/* redirect against allowlist. Test: approved 302, unknown 403 (unit).

## Phase 7 - Binary wiring
- [ ] 7.1 feat(backend): cmd/api-server wires config, db, cache, server, and starts the consumer goroutine with recover(). Test: boot smoke, /health 200 (integration).
- [ ] 7.2 feat(infra): systemd unit + env profile load (REMITTRACK_INFRA_PROFILE). Test: config parse (unit).

## Phase 8 - Scrapers (first full external slice)
- [ ] 8.1 feat(shared): scrapeharness uTLS transport + proxy rotation + backoff. Test: builds client, honors HTTP_PROXY, retries (unit, fake transport).
- [ ] 8.2 feat(lambdas): wise parser fixture -> RateRecord. Test: saved payload -> expected record (unit).
- [ ] 8.3 feat(lambdas): wise handler pushes RateRecord to SQS. Test: fake SQS receives normalized message (unit).
- [ ] 8.4 test(integration): wise -> SQS -> consumer -> Postgres end to end. Test: full path writes expected row.
- [ ] 8.5 feat(lambdas): sendwave parser + handler (repeat 8.2-8.3). Test: fixture and SQS push.

## Phase 9 - Infra + deploy (lean env)
- [ ] 9.1 feat(infra): modules/messaging SQS + DLQ, prevent_destroy. Test: terraform validate + reviewed plan.
- [ ] 9.2 feat(infra): modules/data RDS Multi-AZ, prevent_destroy. Test: validate + reviewed plan.
- [ ] 9.3 feat(infra): envs/lean composes network + data + messaging + compute. Test: validate + reviewed plan.
- [ ] 9.4 feat(infra): deploy/ nginx template + systemd + swap.sh (rewrite-then-reload). Test: shellcheck + swap.sh dry run.
- [ ] 9.5 ci: backend-cd.yml builds, ships, SSH swap, manual prod approval. Test: workflow parses.
- [ ] 9.6 chore(infra): first real deploy of api-server to lean. Test: live /health 200.

## Phase 10 - Frontend (SEO-critical)
- [ ] 10.1 chore(frontend): Next.js App Router init + tsconfig + eslint (Google) + prettier. Test: npm run build.
- [ ] 10.2 feat(frontend): lib/schema JSON-LD builders (ExchangeRateSpecification, FAQPage), pure. Test: builder output shape (unit).
- [ ] 10.3 feat(frontend): lib/api fetch current rates. Test: mock fetch (unit).
- [ ] 10.4 feat(frontend): app/[locale]/exchange-rates/[corridor]/page.tsx ISR page + JSON-LD. Test: render + JSON-LD snapshot.
- [ ] 10.5 feat(frontend): hreflang + canonical tags in layout. Test: header tags present (unit).
- [ ] 10.6 feat(frontend): revalidation webhook endpoint + ingestion fires it. Test: webhook triggers revalidate (integration).
- [ ] 10.7 ci: frontend-cd.yml build + push to S3 + CloudFront invalidate. Test: workflow parses.

## Phase 11 - Agents (they observe a system that must already exist)
- [ ] 11.1 feat(services): services/shared LLM client, S3/Chroma helper, PR opener (all mockable). Test: mocked calls (unit).
- [ ] 11.2 feat(services): agent_healer/linter.go AST allowlist gate. Test: allowed passes; forbidden import and //go:linkname rejected (table-driven unit).
- [ ] 11.3 feat(services): agent_healer/healer.py RAG -> candidate -> linter -> canary -> PR. Test: mocked pipeline never auto-merges (unit).
- [ ] 11.4 feat(services): compliance_sentinel diffs robots/ToS, PRs provider_map toggle. Test: change detected -> PR (mocked unit).
- [ ] 11.5 feat(services): onboarding_agent sniff -> scraper gen -> PR. Test: mocked capture -> PR (unit).
- [ ] 11.6 ci: agent-heal.yml on repository_dispatch. Test: workflow parses.

## Phase 12 - Regression + hardening
- [ ] 12.1 chore(tests): regression/docker-compose.yml mock proxy network. Test: compose up/down.
- [ ] 12.2 test(tests): failure_registry.json schema + first case. Test: schema validates.
- [ ] 12.3 test(integration): regression replay harness; each registry case replays green.
- [ ] 12.4 feat(services): canary statistical gate (3-sigma + provider corridor + cold-start skip). Test: pass/fail bands, new provider held (unit).
- [ ] 12.5 ci: wire regression into the PR gate. Test: failing case blocks merge.

---

## Decisions & changes log
Append newest first. Format: `YYYY-MM-DD [step] what changed and why (link ADR if any)`.

- 2026-07-11 [0.3] Two additions. (1) `go test ./...` also exits non-zero on a module with no
  .go files (like golangci-lint), so the skip-empty guard now covers test/lint/build, not just
  lint. (2) The plan named an `up (docker)` target but no step created the compose file it needs
  and Phase 2 assumes a local Postgres; 0.3 now also adds a minimal dev docker-compose.yml
  (Postgres only). No ADR — dev tooling, not architecture.
- 2026-07-10 [0.2-0.4] golangci-lint v2 returns exit 5 (not clean) on a module with zero .go
  files, so the original 0.2 test ("run exits clean on empty pkgs") was unachievable. Decision
  (Option A): 0.2's acceptance is now `golangci-lint config verify`; the Makefile lint target and
  CI skip modules that have no .go files, so linting is clean on the empty workspace and self-heals
  as packages land. No ADR — tooling behavior, not architecture.