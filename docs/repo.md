# RemitTrack Repository Layout

This tree is the physical expression of the SADD. Module boundaries, the lean/enterprise toggle, and the human-gated deploy path are all visible in the structure itself.

```
remit-track/
├── shared/                              # Shared Go module (own go.mod): single source of truth
│   ├── go.mod
│   ├── raterecord/                      # Canonical RateRecord domain type
│   ├── corridor/                        # Currency-pair and corridor value types
│   └── scrapeharness/                   # uTLS transport, proxy rotation, retry/backoff
├── backend/                             # Go modulith core (own go.mod)
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── api-server/                  # REST API, B2B layer (/api/v1/b2b/*), Stripe webhook
│   │       └── main.go                  # Single modulith binary; SQS consumer runs as a goroutine
│   └── internal/
│       ├── cache/                       # RateCache iface: sync.Map (lean) / ElastiCache (enterprise)
│       │   └── rate_cache.go
│       ├── database/                    # PG pool, idempotent UPSERT, current_rates hot table
│       │   ├── postgres.go
│       │   └── migrations/              # Versioned SQL migrations
│       ├── models/                      # Domain entities (no third-party deps)
│       │   ├── alert.go
│       │   ├── marketing.go             # Masked outbound routing allowlist (/out/*)
│       │   ├── quote.go                 # Persistence view over shared/raterecord
│       │   └── user.go                  # Token auth, Stripe evt_id hashes, tier metadata
│       ├── alerting/                    # In-memory corridor index map and evaluation loop
│       │   └── alert_pool.go
│       ├── ingestion/                   # SQS consumer: batch builder and dedupe
│       │   └── consumer.go
│       └── server/
│           ├── handlers/                # REST, B2B, and chat handlers
│           ├── middleware/              # Token auth, Stripe idempotency, logging
│           └── ratelimit/               # Token-bucket limiter
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
│   │   │   ├── edge/                    # CloudFront, S3, Lambda@Edge (ISR)
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
│   └── provider_map.json                # Scraper registry (scrapers and compliance agent)
├── tests/
│   └── regression/                      # Chaos sandbox for the canary/regression gate
│       ├── docker-compose.yml           # Mock proxy network
│       └── failure_registry.json        # Historical failure cases
├── .github/
│   └── workflows/
│       ├── backend-ci.yml               # go test ./..., vet, build
│       ├── backend-cd.yml               # Cross-compile, upload, SSH swap.sh (prod approval)
│       ├── frontend-ci.yml              # Lint, typecheck, next build
│       ├── frontend-cd.yml              # Build and push static to S3, CloudFront invalidate
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

5. **One modulith binary, by design.** The API, chat, alert pool, and SQS ingestion consumer all run inside a single binary under one systemd unit (`remittrack.service`). This keeps the SADD's in-process hot path intact: the consumer writes a batch and the in-memory alert pool evaluates it immediately, with no DB round trip. The ingestion goroutine must `recover()` at its boundary so a scraper-side panic cannot take down the API. The `internal/ingestion` and `internal/alerting` packages are the seams to extract into a separate `cmd/` entrypoint later, if and when ingestion needs to scale independently.

