# RemitTrack Repository Layout

This tree is the physical expression of the SADD. Module boundaries, the lean/enterprise toggle, and the human-gated deploy path are all visible in the structure itself.

```
remit-track/
├── shared/                              # Shared Go module (own go.mod): single source of truth
│   ├── go.mod
│   ├── raterecord/                      # Canonical RateRecord quote (amounts, fees, methods, class); EffectiveRate ranking key
│   ├── corridor/                        # Currency-pair and corridor value types
│   ├── money/                           # Money value type (int64 micros + ISO currency)
│   └── scrapeharness/                   # uTLS transport, proxy rotation, retry/backoff
├── backend/                             # Go modulith core (own go.mod)
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── api-server/
│   │       └── main.go                  # Single modulith binary; the only code that wires modules together
│   └── internal/                        # Domain-owned modules (ADR 0002); modules never import each other.
│       │                                # Each starts as one package; sub-packages go under <module>/internal/.
│       ├── rates/                       # Owns schema `rates`: Save, Current (ranked), Quote view, cache
│       ├── ingestion/                   # SQS consumer, batching, DLQ, corridor allowlist; declares QuoteWriter + BatchObserver
│       ├── alerting/                    # Owns schema `alerting`: rules, trigger state, outbox; OnBatch
│       ├── accounts/                    # Owns schema `accounts`: users, API keys, Stripe ledger, tiers
│       ├── outbound/                    # Masked /out/* redirect allowlist; no tables
│       ├── httpapi/                     # HTTP transport only: routes, middleware, rate limiting
│       └── platform/
│           └── postgres/                # Technical only: pool, config from env, migration runner
├── lambdas/                             # Stateless serverless ingestion (own go.mod)
│   ├── go.mod                           # Depends on ../shared
│   └── scrapers/                        # Per-provider deployment archives
│       ├── wise/
│       │   └── main.go
│       └── sendwave/
│           └── main.go
├── services/                            # Python multi-agent suite (GitHub Actions sandbox)
│   ├── agent_healer/                    # LangGraph self-healing engine
│   │   ├── healer.py
│   │   ├── linter.go                    # Go AST allowlist gate (invoked as subprocess)
│   │   └── requirements.txt
│   ├── compliance_sentinel/             # robots.txt / ToS diffing, provider_map toggling
│   │   ├── sentinel.py
│   │   └── requirements.txt
│   ├── onboarding_agent/                # playwright-stealth XHR sniffer, scraper generator
│   │   ├── onboard.py
│   │   └── requirements.txt
│   └── shared/                          # LLM client, S3/Chroma helpers, PR opener
│       └── __init__.py
├── frontend/                            # Next.js (App Router) + TypeScript
│   ├── app/
│   │   └── [locale]/
│   │       └── exchange-rates/
│   │           └── [corridor]/
│   │               ├── page.tsx         # ISR page, revalidated by ingestion webhook
│   │               └── opengraph-image.tsx
│   ├── components/                      # ComparisonTable, AlertManager, ChatWidget
│   ├── lib/
│   │   ├── api/                         # Data fetching
│   │   └── schema/                      # ExchangeRateSpecification and FAQPage JSON-LD builders
│   ├── locales/                         # en.json, fr.json, ...
│   ├── public/
│   ├── next.config.ts
│   ├── package.json
│   └── tsconfig.json
├── infra/                               # Infrastructure as Code
│   ├── terraform/
│   │   ├── modules/
│   │   │   ├── network/                 # VPC, security groups
│   │   │   ├── data/                    # RDS Multi-AZ, ElastiCache (enterprise)
│   │   │   ├── messaging/               # SQS and dead-letter queue
│   │   │   ├── ingestion/               # Lambda functions, EventBridge cron
│   │   │   ├── edge/                    # CloudFront, S3 origin, scoped rebuild + s-maxage TTL (ISR; invalidation is deploy-only)
│   │   │   └── compute/                 # EC2 (lean) / ALB + Fargate (enterprise)
│   │   ├── envs/
│   │   │   ├── lean/                    # Single EC2 + single RDS, prevent_destroy guards
│   │   │   │   ├── main.tf
│   │   │   │   ├── variables.tf
│   │   │   │   └── outputs.tf
│   │   │   └── enterprise/
│   │   │       ├── main.tf
│   │   │       ├── variables.tf
│   │   │       └── outputs.tf
│   │   └── README.md
│   └── deploy/
│       ├── nginx/
│       │   └── go_modulith.conf.tmpl    # Blue/green upstream template
│       ├── systemd/
│       │   └── remittrack.service      # Single unit for the modulith binary
│       └── swap.sh                      # Rewrite upstream, health-check, reload nginx
├── config/
│   └── provider_map.json                # Provider registry: per-provider ingestion strategy (api/scrape), quote endpoint, and pricing config (pay-in/out methods, rate classes, per-corridor overrides); read by scrapers, compliance, onboarding
├── tests/
│   └── regression/                      # Chaos sandbox for the canary/regression gate
│       ├── docker-compose.yml           # Mock proxy network
│       └── failure_registry.json        # Historical failure cases
├── .github/
│   └── workflows/
│       ├── backend-ci.yml               # go test ./..., vet, build
│       ├── backend-cd.yml               # Cross-compile, upload, SSH swap.sh (prod approval)
│       ├── frontend-ci.yml              # Lint, typecheck, next build
│       ├── frontend-cd.yml              # Build and push static to S3, CloudFront invalidate (deploy only, not per-tick)
│       ├── infra-plan.yml               # terraform plan on PR, posts plan for review
│       ├── infra-apply.yml              # terraform apply on manual approval
│       └── agent-heal.yml               # Episodic self-healing runner
├── docs/
│   ├── design.md
│   ├── sadd.md
│   └── repo.md
├── go.work                              # Ties shared/, backend/, lambdas/ for local dev
└── README.md
```

## Notes

1. **Three Go modules, one domain model.** `shared/` is its own module holding the canonical `RateRecord` and the uTLS scrape harness. `backend/` and `lambdas/` both depend on it, so the domain type and network-evasion code are never duplicated across scrapers. `go.work` ties the modules together for local development.

2. **Frontend is Next.js App Router, not a SPA.** The SADD's SEO strategy (Edge ISR, per-corridor static generation, hreflang, `ExchangeRateSpecification` and `FAQPage` JSON-LD) requires server and build-time rendering. The `app/[locale]/exchange-rates/[corridor]/` route and `lib/schema/` builders are the load-bearing pieces.

3. **Infra separates modules from environments.** `terraform/modules/` are reusable; `terraform/envs/lean` and `terraform/envs/enterprise` compose them. Lean vs enterprise is an explicit environment choice, not a fork, and `prevent_destroy` guards live on the stateful modules. `deploy/` is the home the SADD's blue/green swap previously lacked: the nginx upstream template, the systemd unit, and `swap.sh`.

4. **CI and CD are distinct; infra is human-gated.** `*-ci` runs tests and builds on every PR. `*-cd` ships. `infra-plan.yml` posts a `terraform plan` to the PR and `infra-apply.yml` runs only after manual approval, matching the SADD's review gate.

5. **One modulith binary, by design.** The API, chat, alert pool, and SQS ingestion consumer all run inside a single binary under one systemd unit (`remittrack.service`). This keeps the SADD's in-process hot path intact: the consumer writes a batch and the in-memory alert pool evaluates it immediately, with no DB round trip. The ingestion goroutine must `recover()` at its boundary so a scraper-side panic cannot take down the API. Inside the binary the code is split into domain-owned modules (ADR 0002): each owns its tables in its own Postgres schema, starts as one package (sub-packages go under its own `internal/` only once it splits), and declares small interfaces for what it needs from others, which `cmd/api-server` wires together. Because modules never import each other, any of them can later move to its own `cmd/` entrypoint or service by changing the wiring and moving its schema.

6. **Roadmap items not yet in the tree.** Three SADD features are intentionally deferred and will add structure when scheduled (see implementation-plan Phase 13): the operator market intelligence dashboard (a rollup-backed, B2B-gated read surface; v1-vs-deferred still to be decided), the enterprise compute profile (ALB + ECS Fargate, already stubbed at `infra/terraform/modules/compute`), and the tier reconciliation job. API-first ingestion is already reflected: each `provider_map.json` entry declares an `api` or `scrape` strategy, and the scraper Lambdas try the official API before falling back.

