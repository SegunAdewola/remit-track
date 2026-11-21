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
- [x] 0.2 chore(lint): add .golangci.yml (govet, staticcheck, errcheck, revive, gofumpt) using the v2 schema. Test: `golangci-lint config verify` passes. (A full `golangci-lint run` returns exit 5 on a module with no .go files; it becomes green once the first package lands at 1.1.)
- [x] 0.3 chore(build): Makefile (test, lint, build, up/down targets) + dev docker-compose.yml (Postgres for Phase 2 integration tests). The test/lint/build targets skip modules with no .go files so they are clean on the empty workspace. Test: `make test` runs.
- [x] 0.4 ci: backend-ci.yml runs `make lint` + `make test` on PR and push to main; Go version read from go.work, golangci-lint pinned to the local version. Test: `actionlint` passes.

## Phase 1 - Shared domain (foundation, no dependencies)
- [x] 1.1 feat(shared): corridor CurrencyPair type with Parse/String ("USD:GHS"). Test: parse valid/invalid, round-trip.
- [x] 1.2 feat(shared): money.Money value type (integer minor units + ISO currency) with Validate (currency is 3 uppercase letters; sign is the caller's concern). Test: valid/invalid currency (table-driven unit). (Display formatting deferred: needs per-currency exponent metadata - JPY=0 vs USD=2.)
- [x] 1.3 feat(shared): RateRecord quote struct (Provider, Pair, SendAmount/TotalPaid/ReceiveAmount/Fees as Money, OfferedRate float64, PayIn/PayOut as PayMethod, Class as RateClass, Timestamp) + Validate, plus RateClass.Valid over the closed set {standard, promo_first_transfer}. Validate is pure: it checks each field's own rule and currency/side matching only (Send/TotalPaid in Base, Receive in Quote, Fees in either side of the pair), never cross-field arithmetic, and checks the pair is well-formed (valid codes, Base != Quote) but NOT that it is a supported corridor (that allowlist check lives at ingestion, 4.2). Test: valid and each invalid case - empty provider, malformed pair, non-positive amounts, negative fees, currency/side mismatch, non-finite or non-positive rate, empty method, bad class, zero timestamp (table-driven unit).
- [x] 1.4 feat(shared): money.MinorUnits(code) (int, bool), an owned ISO 4217 minor-unit table for the currencies RemitTrack serves (CAD, USD, GBP, EUR, NGN, GHS, KES = 2; XOF, XAF = 0). The single source of truth for both parsers (8.2) and EffectiveRate; adding a corridor adds a row. Test: each listed currency pinned to its value; unknown and malformed codes return false (table-driven unit).
- [ ] 1.5 refactor(shared): store Money as fixed-scale micros (ADR 0001). Field `Amount` becomes `Micros int64` (millionths of a whole unit for every currency), plus const `MicrosPerUnit = 1_000_000`. MinorUnits now drives validation and display only. RateRecord code and tests converted to micros with the same cases and expected errors. Test: money tests in micros + MicrosPerUnit pinned; all shared tests green.
- [ ] 1.6 feat(shared): RateRecord.EffectiveRate() (float64, error) = float64(ReceiveAmount.Micros) / float64(TotalPaid.Micros). No per-pair scaling (micros are the same scale for every currency). Calls Validate first. Test: known CAD->NGN ratio, CAD->XOF regression (430, not 4.3), invalid record returns an error; float comparisons use a tolerance (unit).
- [ ] 1.7 feat(shared): RateRecord JSON marshal/unmarshal with stable field names. Test: round-trip + field-name assertions.

## Phase 2 - Persistence (the source of truth)
- [ ] 2.1 feat(backend): internal/models structs (Quote view = a ranked comparison result carrying EffectiveRate + promo badge, AlertConfig, User) + validation. Test: zero-value and validation (unit).
- [ ] 2.2 feat(database): migration 0001 historical_rates + current_rates with full quote columns (send_amount, total_paid, receive_amount, fees, offered_rate, pay_in, pay_out, rate_class; money as bigint micros, e.g. send_amount_micros, per ADR 0001) + UNIQUE(provider, currency_pair, send_amount, pay_in, pay_out, rate_class, scrape_timestamp). Test: up/down applies (integration, real Postgres).
- [ ] 2.3 feat(database): connection pool + config from env. Test: connect and ping (integration).
- [ ] 2.4 feat(database): idempotent UPSERT into historical_rates (ON CONFLICT on the full quote key DO NOTHING). Test: duplicate ignored, distinct timestamps/amounts kept (integration).
- [ ] 2.5 feat(database): current_rates upsert on insert + read-latest-by-quote-key query ranked by effective rate. Test: one row per quote key, historical not scanned (integration).

## Phase 3 - Cache
- [ ] 3.1 feat(cache): RateCache interface + sync.Map implementation. Test: get/set/miss + concurrent access (unit).

## Phase 4 - Ingestion (SQS consumer, in-process)
- [ ] 4.1 feat(ingestion): batch builder accumulates messages into multi-row batches by size and time. Test: batching rules (unit).
- [ ] 4.2 feat(ingestion): consumer long-polls a queue, decodes RateRecord, validates it, rejects corridors not in the injected supported-corridors allowlist, writes via database. Test: mock SQS -> rows in Postgres; unsupported corridor rejected (integration).
- [ ] 4.3 feat(ingestion): recover() boundary + poison-message handling (log/DLQ, loop survives). Test: bad message does not kill the loop (integration).

## Phase 5 - Alerting (in-memory evaluation, durable rules + delivery)
- [ ] 5.1 feat(alerting): corridor index map build/refresh from AlertConfig; DB is source of truth, index rebuilt on boot. Test: build, lookup, and rebuild-from-DB (unit).
- [ ] 5.2 feat(alerting): evaluation loop fires on threshold crossing only, with persisted per-alert trigger state (armed/triggered) for edge-triggered dedup. Test: crossing fires once, stays-crossed does not re-fire, re-arms after reset (unit).
- [ ] 5.3 feat(alerting): transactional outbox - firing writes a notifications row (status=pending) in the SAME tx as the trigger-state flip. Test: fired alert commits a pending outbox row atomically (integration).
- [ ] 5.4 feat(alerting): notification dispatcher drains pending outbox via SES (behind iface, fake in test), marks sent, dedups by notification_id. Test: pending row delivered once; crash-before-send leaves pending; retry does not double-send (unit).
- [ ] 5.5 feat(backend): wire ingestion batch completion to alerting evaluation in-process. Test: batch write triggers eval (integration).

## Phase 6 - Server (read path first, then edges)
- [ ] 6.1 feat(server): current-rates-by-corridor handler backed by cache/DB, ranked by effective rate (received/paid), standard rate class only with promos flagged. Test: 200 + correctly ranked payload (httptest integration).
- [ ] 6.2 feat(server): logging + panic-recovery middleware. Test: handler panic returns 500, server survives (unit).
- [ ] 6.3 feat(server): token-bucket ratelimit middleware. Test: over-limit returns 429 (unit).
- [ ] 6.4 feat(server): B2B token auth middleware for /api/v1/b2b/*. Test: valid vs invalid key (integration).
- [ ] 6.5 feat(server): Stripe webhook idempotency via evt_id ledger, plus users.tier derivation from subscription state. Test: duplicate event ignored; tier updates on event (integration).
- [ ] 6.6 feat(server): masked /out/* redirect against allowlist. Test: approved 302, unknown 403 (unit).

## Phase 7 - Binary wiring
- [ ] 7.1 feat(backend): cmd/api-server wires config, db, cache, server, and starts the consumer goroutine with recover(). Test: boot smoke, /health 200 (integration).
- [ ] 7.2 feat(infra): systemd unit + env profile load (REMITTRACK_INFRA_PROFILE). Test: config parse (unit).

## Phase 8 - Scrapers (first full external slice)
- [ ] 8.1 feat(shared): scrapeharness uTLS transport + proxy rotation + backoff. Test: builds client, honors HTTP_PROXY, retries (unit, fake transport).
- [ ] 8.2 feat(lambdas): wise quote parser: fixture -> RateRecord quote (amounts, fees, methods, rate class - the observed output). Test: saved payload -> expected quote (unit).
- [ ] 8.3 feat(lambdas): wise handler pushes RateRecord to SQS. Test: fake SQS receives normalized message (unit).
- [ ] 8.4 test(integration): wise -> SQS -> consumer -> Postgres end to end. Test: full path writes expected row.
- [ ] 8.5 feat(lambdas): sendwave parser + handler (repeat 8.2-8.3). Test: fixture and SQS push.
- [ ] 8.6 feat(lambdas): provider_map.json ingestion-strategy field + API-first client (official API preferred, scrape fallback on absence or failure). Test: api path chosen when present, falls back on failure (unit).

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
- [ ] 10.4 feat(frontend): app/[locale]/exchange-rates/[corridor]/page.tsx static page + JSON-LD, with a real (build-time) rate baked into the HTML as the crawler-visible baseline and no-JS fallback. Test: render + JSON-LD snapshot + baked rate present in server HTML.
- [ ] 10.5 feat(frontend): hreflang + canonical tags in layout. Test: header tags present (unit).
- [ ] 10.6 feat(frontend): revalidation webhook endpoint + ingestion fires it; scoped rebuild writes S3 objects with a short s-maxage TTL and does NOT invalidate CloudFront per tick. Test: webhook triggers revalidate + written objects carry Cache-Control s-maxage (integration).
- [ ] 10.7 ci: frontend-cd.yml build + push to S3 + CloudFront invalidate (deploy path only, for template changes; the per-tick ingestion path relies on s-maxage TTL, see 10.6). Test: workflow parses.
- [ ] 10.8 feat(frontend): live rate hydration on corridor pages - small client fetch of the cached rates API that freshens the baked-in baseline, degrading gracefully to the baked value when JS is off. Interactive surfaces (alerts/dashboard/chat) are client-rendered against the same API, no static rebuild. Test: hydration updates the number; no-JS shows baked rate (unit + render).

## Phase 11 - Agent guardrails, corpus, and gates (build before any code generation)
- [ ] 11.1 feat(observability): structured scraper error/empty-payload emission + CloudWatch metric filter + repository_dispatch webhook. Test: an error log produces a dispatch payload (unit + integration).
- [ ] 11.2 chore(tests): regression/docker-compose.yml mock proxy network that serves historical payloads. Test: compose up/down.
- [ ] 11.3 test(tests): failure_registry.json schema + first case. Test: schema validates.
- [ ] 11.4 feat(services): agent_healer/linter.go AST allowlist gate. Test: allowed passes; forbidden import and //go:linkname rejected (table-driven unit).
- [ ] 11.5 feat(services): canary statistical gate (3-sigma + provider corridor + cold-start skip). Test: pass/fail bands, new provider held (unit).
- [ ] 11.6 test(integration): regression replay harness; each registry case replays green.

## Phase 12 - Agents (rely on the Phase 11 guardrails and corpus)
- [ ] 12.1 feat(services): services/shared LLM client, S3/Chroma helper, PR opener (all mockable). Test: mocked calls (unit).
- [ ] 12.2 feat(services): agent_healer/healer.py RAG -> candidate -> linter -> canary -> PR, max 3 repair loops then escalate. Test: mocked pipeline never auto-merges; stops after 3 (unit).
- [ ] 12.3 feat(services): compliance_sentinel diffs robots/ToS, PRs provider_map toggle. Test: change detected -> PR (mocked unit).
- [ ] 12.4 feat(services): onboarding_agent sniff -> scraper gen -> PR (last: highest blast radius). Test: mocked capture -> PR (unit).
- [ ] 12.5 ci: agent-heal.yml on repository_dispatch. Test: workflow parses.
- [ ] 12.6 ci: wire regression into the PR gate. Test: failing case blocks merge.

## Phase 13 - Deferred / roadmap (pending explicit scope decisions)
- [ ] 13.1 feat(backend): operator dashboard rollup (materialized view or summary table) + read-only B2B-gated endpoints. DECISION NEEDED: v1 or deferred.
- [ ] 13.2 feat(infra): enterprise compute profile (CloudFront -> ALB -> ECS Fargate) behind REMITTRACK_INFRA_PROFILE; resolve alert-pool ownership (ElastiCache vs single owner) before enabling. Deferred until scale requires it.
- [ ] 13.3 feat(server): tier reconciliation job (periodic Stripe re-read to correct missed-webhook drift). Deferred; safe to add post-launch.
- [ ] 13.4 feat(server): on-demand live quoting for exact amounts - fan-out to provider quote endpoints, deadline-bounded, single-flight coalesced, short-TTL cached, partial results tolerated; API-quote providers only. Deferred until the benchmark cache proves insufficient.
- [ ] 13.5 feat(backend): landing default resolver - send country from CloudFront-Viewer-Country, corridor from a popularity prior + remembered selector, amount from user-last -> corridor median -> global fallback (static shows population default, hydration personalizes). Deferred; needs usage data.
- [ ] 13.6 feat(backend): demand-shaped benchmark amounts - track per-corridor requested-amount distribution and tune the precomputed benchmark set (higher / more granular) to maximize cache hit rate. Deferred; needs traffic.

---

## Open follow-ups (not yet decided; resolve at the trigger shown)
Items raised in discussion that need a decision or a plan edit. When one is resolved, delete it here
and record the outcome in the Decisions & changes log below.

- [ ] (before 2.1, plan mode + likely ADR) Modulith boundaries: replace the planned kind-grouped
  `internal/models` + shared `internal/database` with domain-owned modules (e.g. rates, alerting,
  users/billing), each owning its types and tables (Parnas: hide each module's secret). Update
  repo.md, sadd.md and Phases 2/5/6 accordingly.
- [ ] (before 7.1) Write performance hypotheses as numbers (e.g. cached read p99 at N rps, consumer
  records/s, alert eval latency at 10k rules) in sadd.md.
- [ ] (after 7.1) Add a step for a k6 baseline load test with thresholds as pass/fail (dev tool via
  Homebrew, like actionlint). Prometheus metrics via client_golang later: new dependency, discuss.
- [ ] (Phase 10) Adapt the V0-built frontend for v1: split `"use client"` components into a server
  shell + client islands so corridor pages render the rate server-side; decide on shadcn/Tailwind
  dependencies. Confirm where the V0 export lives.
- [ ] (before 4.2 and Phase 9) Local cloud emulation: evaluate LocalStack (and any AWS-native
  emulator) to run SQS/Lambda/S3/SSM locally and to `terraform apply` against before deploying.
  Could replace the hand-rolled mock SQS in tests/regression. Check current licensing/free-tier
  coverage (RDS is not in LocalStack's free tier; keep real Postgres via docker compose).

### API contract and HTTP stack (decide in plan mode before Phase 6)
- [ ] Add a step before 6.1: write the public API contract as OpenAPI (contract-first, decided
  2026-11-14) in a versioned file, e.g. `backend/api/openapi.yaml`. Handlers in Phase 6 implement it;
  the frontend (10.3) and B2B docs read it.
- [ ] Decide whether to generate Go types/server stubs from the spec (oapi-codegen, std-http target)
  or hand-write handlers and add a conformance test against the spec. Either way a new dependency
  or tool: discuss tradeoff first.
- [ ] HTTP stack: leaning stdlib `net/http` with Go 1.22+ method/path patterns
  (`"GET /api/v1/rates/{corridor}"`, `r.PathValue`) and `func(http.Handler) http.Handler`
  middleware; no framework. chi is the fallback if routing outgrows ServeMux. Avoid Gin/Echo
  (framework-specific handler types) and Fiber (fasthttp, not net/http).
- [ ] Keep the public API shape separate from RateRecord's internal SQS JSON (1.6): the API is a
  contract with customers, the SQS message is an internal detail that may change.
- [ ] Decide the API's error response shape (e.g. RFC 9457 problem+json) and versioning policy
  (`/api/v1`, what counts as breaking) as part of the contract.
- [ ] Postman/Newman is optional: import the OpenAPI spec for manual exploration and smoke tests
  against a deployed server (from 7.1/9.6). Correctness stays in Go httptest in CI.

### Feature flags (decide before 6.1; wire at 7.1)
- [ ] Every user-facing feature must be switchable off without a redeploy of code: chat, alerts,
  B2B API, live quoting, operator dashboard. Leaning: plain config flags (env / SSM Parameter Store)
  read at boot that gate route registration and background goroutines; no third-party flag
  service. Disabled = route not registered + goroutine not started, not just hidden in the UI.
- [ ] LLM chat is a learning exercise, not a kept v1 feature: build it behind a flag defaulting to
  off in production. It has no plan step yet; add one (likely after Phase 7) when we get there.

### Cost & simplification
Decide together in plan mode after Phase 1, before 2.2 builds the schema. On decision, update
sadd.md, design.md, repo.md and the affected steps, and write an ADR. Estimated 2026-11-13: current design ~$52-55/mo infra + $25-300 residential proxy + ~$5-30 LLM.
Candidate lean target ~$27/mo infra, ~$40-80 all-in. Figures are hand-estimated list prices;
verify in the AWS Pricing Calculator before deciding.
- [ ] RDS Multi-AZ -> single-AZ with automated backups + PITR until revenue (~$14/mo saved; failover
  becomes a ~30 min restore, no data loss beyond PITR). Flip back is a Terraform flag.
- [ ] Edge on Cloudflare free (DNS, CDN, static frontend via Pages, WAF/DDoS in front of the API);
  AWS keeps compute + data. Max two vendors; stateful pieces stay next to their compute.
  Replaces CloudFront + WAF + Route 53 (~$10/mo saved).
- [ ] Scrape cadence 5 min -> ~15 min, and API-first providers bypass the residential proxy.
  The proxy is the dominant cost (~$230-290/mo at 5 min, 450 quotes/tick; ~$10-40 after both).
- [ ] Drop the per-tick ingestion -> webhook -> scoped rebuild pipeline (10.6): since hydration
  (2026-09-18) only crawlers need the static rate, so rebuild on a schedule (every few hours) +
  on deploy. Simplifies Phase 10.
- [ ] SQS: one batched message per provider per tick (not one per quote) to stay in the free tier.
  Fold into 8.3.
- [ ] (8.2) Parsers must reject amounts with more decimals than money.MinorUnits allows (never
  round), so a provider-side precision change fails loudly instead of being stored.
- [ ] (before 1.6) Decide whether RateRecord.Validate should also reject amounts with more precision
  than their currency allows (Micros not a multiple of 10^(6 - MinorUnits)), or leave that to
  parsers only. Doing it in Validate needs every currency to be in the table.
- [ ] (Phase 11, with 11.1) Scheduled ISO 4217 check: a weekly GitHub Actions job fetches the official
  list published by SIX Group (ISO's maintenance agency), compares minor units and withdrawals for
  supported codes against money.MinorUnits, and opens an issue (never an auto-PR) on any difference.
  Runbook on an issue: new code (redenominations usually get one, e.g. GHC->GHS, ZWL->ZWG) -> add a
  row, accept both codes through the overlap, then retire the old one; same code with new minor
  units -> update the table, test and parser fixtures before the effective date. Stored amounts are
  micros (ADR 0001), so no rescaling migration is needed.
- [ ] historical_rates growth (~0.6-1 GB/mo at the above volumes): decide a retention/rollup policy
  (e.g. raw 30-90 days, hourly rollups kept) before 2.2.
- [ ] Compute stays EC2 (or Lightsail; verify pricing). Rejected for now: ECS Fargate (needs an ALB,
  ~$28-32/mo; revisit at the enterprise profile, where it also replaces nginx blue/green) and
  Lambda-only (breaks the single-process alert pool/cache/consumer; RDS Proxy ~$22/mo minimum).
- [ ] Keep Lambda + SQS for scrapers (not in-process goroutines): the scrapers are the code agents
  rewrite, so Lambda isolates generated code from the API process and its secrets.
- [ ] Premature for v1, consider deferring: Multi-AZ, 5-min scraping, per-tick rebuilds, RAG over an
  empty failure corpus (start the healer with the registry + fixtures, add RAG when there is history).

### Security (each tagged with the step that must include it)
- [ ] (Phase 12 / 12.5) Agents read attacker-controlled text (ToS pages, sniffed XHR). Least-privilege
  GITHUB_TOKEN in agent workflows, no secrets exposed to agent-PR-triggered workflows (avoid
  `pull_request_target`), AWS via GitHub OIDC with narrow roles, never long-lived keys.
- [ ] (10.2 / FAQ mining) Chat text flows into FAQPage JSON-LD inside a <script> tag: escape `<`
  (and `</script>`) in JSON-LD serialization and moderate before publishing (stored XSS risk).
- [ ] (6.3 + chat step) Denial-of-wallet: per-IP/per-user limits and max_tokens on chat; AWS Budgets
  alerts (free) plus hard spend caps on the LLM and proxy accounts.
- [ ] (9.4 / 9.5) Deploy via SSM Run Command instead of SSH from Actions: no inbound port 22.
- [ ] (7.1 / Phase 9) Secrets (Stripe, LLM, DB) in SSM Parameter Store SecureString; never .env on
  the host or in the repo.
- [ ] (0.4 follow-on) Supply chain: pin deps, add govulncheck + npm audit + Dependabot to CI.
- [ ] (8.1 / product) Legal/partner risk of TLS-fingerprint forging + residential proxies to bypass
  bot protection on providers who are also prospective affiliate partners. Prefer official APIs and
  affiliate rate feeds; justify evasion per provider.

## Decisions & changes log
Append newest first. Format: `YYYY-MM-DD [step] what changed and why (link ADR if any)`.

- 2026-11-16 [1.5-1.7, 2.2, ADR 0001] Money is stored as fixed-scale micros (int64 millionths of a
  whole unit, every currency) instead of per-currency minor units. Stored amounts no longer depend on
  money.MinorUnits, so a future minor-unit change cannot misread them, and EffectiveRate needs no
  per-pair scaling. MinorUnits stays for validation and display. Inserted 1.5 (the refactor);
  EffectiveRate is now 1.6 and JSON 1.7; 2.2 columns are bigint micros. Dropped the "store each
  row's scale" follow-up. See docs/decisions/0001-money-as-fixed-scale-micros.md.
- 2026-11-14 [Phase 6] The public API will be designed contract-first in OpenAPI, before handlers,
  because the frontend and B2B customers depend on it and field names become hard to change once
  consumed. A step before 6.1 will be added in plan mode; codegen vs hand-written handlers and the
  HTTP stack (leaning stdlib net/http) are open follow-ups. No ADR yet.
- 2026-11-14 [1.4-1.6] Inserted step 1.4 (money.MinorUnits) before EffectiveRate; EffectiveRate is now
  1.5 and JSON 1.6. Dividing raw minor units gives a wrong rate whenever the two currencies have
  different minor units (CAD->XOF reads 4.3 instead of 430). Chose an owned ISO 4217 table over
  golang.org/x/text/currency: stored amounts are interpreted through this table, so it must change
  only by our own reviewed commit, never by an upstream data update (x/text uses CLDR display data,
  which differs from ISO 4217 for some currencies). This matches how payment processors such as
  Stripe and Adyen keep their own minor-unit lists. EffectiveRate returns (float64, error) and
  validates first, avoiding in-band errors (Inf/NaN) per Google Go style. Closes the x/text follow-up.
  No ADR.
- 2026-11-13 [1.3, 4.2] Settled the RateRecord quote interface before writing its test. (1) "Unknown
  pair" split in two: RateRecord.Validate stays pure and checks only that the pair is well-formed
  (valid codes, Base != Quote); the supported-corridors allowlist is injected config checked at
  ingestion (4.2), keeping shared/ free of config and global state. (2) Fees may be in either
  currency of the pair (none seen in the receive currency yet, but cheap flexibility). (3) RateClass
  and PayMethod are named string types; RateClass is a closed set with Valid(), PayMethod only needs
  to be non-empty for now. (4) Validate does no cross-field arithmetic (e.g. TotalPaid ==
  SendAmount + Fees) because fee models differ (added on top vs deducted) and the record holds
  observed values. OfferedRate stays float64: displayed, never used in money arithmetic. No ADR.

- 2026-11-13 [1.2 follow-up] corridor now validates codes via money.IsValidCurrency and its private
  copy is deleted, so there is one definition of a well-formed currency code. Dependency direction is
  corridor -> money because a pair is built from currencies, not the reverse. Corridor tests stayed
  green unedited. No ADR.

- 2026-09-19 [Phase 1/2/6/8/13, docs/architecture] Reshaped the rate model into a fee-inclusive quote so comparisons are accurate. RateRecord grows from a single `Rate` into a quote (Money-typed SendAmount/TotalPaid/ReceiveAmount/Fees, OfferedRate, PayIn/PayOut, rate Class) captured as the provider's *observed* output — never a reconstructed fee formula, because a miscomputed price is indistinguishable from a correct one (the worst failure for a comparison product). Ranking is by effective rate = received/paid, standard class only, promos shown as a badge. Money is integer minor units (no float). Schema: historical_rates/current_rates carry the full quote columns and the idempotency key expands to (provider, currency_pair, send_amount, pay_in, pay_out, rate_class, scrape_timestamp). Read path serves benchmark-cached quotes by default, with on-demand deadline-bounded, single-flighted, short-TTL-cached live quoting for exact amounts (API-quote providers only); the landing default resolves send country from CloudFront-Viewer-Country, corridor from a popularity prior + remembered selector, and amount from user-last → corridor median → global fallback. Plan: Phase 1 reordered/expanded (1.2 Money, 1.3 enriched RateRecord + Validate, 1.4 EffectiveRate, 1.5 JSON); Phase 2 schema/keys updated (2.1/2.2/2.4/2.5); 6.1 ranks by effective rate; 8.2 emits quotes; roadmap 13.4-13.6 added (live quoting, default resolver, demand-shaped benchmarks). Updated sadd.md Sections B & C (+ trade-offs) and design.md. No ADR yet.
- 2026-09-18 [Phase 10, docs/architecture] Added an explicit per-route rendering strategy to resolve the "why not a plain SPA + cached API" question. Rendering mode is chosen per route by "does a crawler need to see this": public corridor pages are statically rendered with a real rate baked into the HTML (Googlebot only runs JS on a deferred pass, and PerplexityBot/GPTBot run none, so a client-only page is blank to the AI crawlers we court); logged-in/interactive surfaces (alerts, dashboard, chat, B2B console) are client-rendered against the cached API with no static rebuild. Adopted "static baseline + live client hydration" for corridor pages: HTML ships a build-time rate (crawler-visible, no-JS fallback), a small client fetch freshens it live. This decouples crawler freshness (hourly is ample) from human freshness (sub-minute via API), letting the static rebuild cadence relax below the 5-minute ingestion tick and shrinking the edge pipeline. Updated sadd.md Section A (new blueprint item 6 + trade-off bullet), and Phase 10 (10.4 bakes the baseline rate, new 10.8 for hydration + interactive surfaces). No ADR yet.
- 2026-09-18 [Phase 5, docs/architecture] Closed a durability gap in the alert pool. Previously the design described only an in-memory `map[string][]AlertConfig` evaluated by a worker pool with fire-and-forget SES dispatch, which risked losing alert rules on restart and dropping (or duplicating) notifications on a crash between "fired" and "sent". New design separates disposable compute from durable state: (A) alert rules + per-alert trigger state are persisted in PostgreSQL as source of truth, and the in-memory index is a cache rebuilt on boot; (B) delivery uses a transactional outbox — the trigger-state flip and a `notifications` (status=pending) row commit in one tx, and a separate dispatcher drains pending rows via SES, marking sent and deduping by notification_id, so at-least-once delivery replaces silent loss. Firing is edge-triggered (armed→triggered→armed) to avoid re-firing a still-crossed threshold every batch. This durability axis is orthogonal to the existing in-memory-vs-ElastiCache scaling trade-off. Updated sadd.md Section B (blueprint item 2 + new trade-off bullet), design.md alert-pool feature bullet, and Phase 5 steps (now 5.1-5.5, adding the outbox write and dispatcher). No ADR yet.
- 2026-09-18 [10.6/10.7, docs/architecture] Reversed the per-tick CloudFront invalidation decision from 2026-07-11. The ISR ingestion path (every 5 min, per corridor) now relies on a short `Cache-Control: s-maxage` TTL instead of calling the Invalidation API. Rationale: invalidating ~288×/day per corridor costs $0.005/path beyond the 1,000/month free tier — roughly $860/mo at 20 corridors, $2,000+/mo at 50 — for a freshness guarantee a TTL aligned to the cadence already provides. Tradeoff: a page may be up to ~1 cadence stale at a given edge; immaterial for crawler/SEO freshness and covered by the visible "last updated" timestamp. If instant freshness is ever needed the escalation is versioned object keys + pointer, not per-tick invalidation. Invalidation is retained for the frontend *deploy* path only (frontend-cd.yml, 10.7), where the HTML template changes rarely and unpredictably. Updated sadd.md Section A (blueprint step 3 + diagram + new trade-off bullet), design.md ISR bullet, repo.md edge/ and frontend-cd.yml comments. ADR still to be written (folds into the pending ISR ADR from 2026-07-11).
- 2026-07-11 [docs/architecture] SADD and design.md reconciliation. Resolved three contradictions: (1) ISR regeneration is a scoped rebuild-to-S3 then CloudFront invalidation, not Lambda@Edge (Lambda@Edge deferred until per-user rendering is needed); (2) TLS terminates at CloudFront, NGINX serves HTTP :80 with the origin locked to CloudFront; (3) B2B path standardized to /api/v1/b2b/*. Added to the SADD: API-first fallback ingestion (provider_map strategy field, now plan step 8.6), operator dashboard (rollup-sourced, v1 TBD), enterprise ALB/Fargate compute profile + alert-pool coupling note, subscription tiers (users.tier, folded into 6.5) + reconciliation, JA4/HTTP2 emulation detail, 3-repair-loop cap, LLM chat deep-dive, and the durable-queue-over-in-process-channel trade-off. Deepened Section E (self-improving loop, six-gate trust ladder) and Section 3 (failure registry as living corpus, chaos build-order). Reordered plan so guardrails + corpus + gates (Phase 11) precede agents (Phase 12); added Phase 13 roadmap for deferred items. ADRs still to be written for the three contradiction resolutions and the enterprise-compute/alert-pool coupling. Architecture-level.
- 2026-07-11 [0.4] Validate "workflow parses" with actionlint (+ shellcheck) locally instead of
  act/push: act needs Docker and push is remote, while actionlint checks the Actions schema and
  shell-checks run: scripts. This is now the standard check for every workflow step (9.5, 10.7,
  11.6). Dev tool, installed via Homebrew. No ADR.
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