# ADR 0002: Organize the backend as domain-owned modules

- Status: Accepted
- Date: 2026-11-21
- Plan steps: Phases 2 to 7

## Context

The backend is a single Go binary (a modulith). Its planned layout grouped code by kind: one
`internal/models` package holding every domain type (quotes, alert configs, users, the outbound
allowlist) and one `internal/database` package holding every table and query.

Grouping by kind gives every part of the backend access to every other part's data and SQL. A change
to how alerts are stored would touch the shared models and database packages, and nothing would stop
the HTTP layer or ingestion from reading alert tables directly. The binary would have folders but no
enforced boundaries, so splitting a part out later would mean untangling it first.

## Decision

Group the backend by domain. Each module owns one area and hides how it works (Parnas's information
hiding):

| Module | Owns | Hides |
|---|---|---|
| `rates` | `rates` schema: `historical_rates`, `current_rates` | storage, idempotent writes, ranking, caching |
| `ingestion` | no tables | SQS handling, batching, the dead-letter queue, the supported-corridor allowlist |
| `alerting` | `alerting` schema: alert configs, trigger state, notifications outbox | the in-memory index, edge triggering, email delivery |
| `accounts` | `accounts` schema: users, API keys, Stripe event ledger | key hashing, Stripe, subscription tiers |
| `outbound` | no tables | the affiliate redirect allowlist |
| `httpapi` | no tables | routing, middleware, rate limiting |

`platform/` packages (for example `platform/postgres`) hold technical code with no domain logic and
may be imported by any module.

Three rules keep the boundaries real:

1. **Each module's tables live in its own Postgres schema.** No module reads or joins another
   module's tables.
2. **Modules do not import each other.** A module declares small interfaces for what it needs. For
   example, `ingestion` declares a `QuoteWriter` and a `BatchObserver`. `cmd/api-server` is the only
   code that knows every module, and it connects them. Only `shared/` types such as `RateRecord`
   cross module boundaries.
3. **A module's implementation lives under `<module>/internal/`.** The Go compiler then stops other
   modules from importing it. A depguard rule in golangci-lint stops modules from importing each
   other's top-level packages.

## Consequences

- Each module can be tested on its own, with small fakes for the interfaces it declares, so most
  tests need no database.
- Storage changes stay inside one module and its schema.
- Extracting a module into its own service means changing the wiring in `cmd/api-server` and moving
  its schema, not untangling shared packages.
- Database permissions can later be granted per schema.
- There is more code up front: interface declarations in each consumer and wiring in `main`.
- Data that two modules need, such as fresh quotes for alert evaluation, must be passed explicitly
  (alerting receives each batch through `OnBatch`), never read from the other module's tables.

## Alternatives considered

- **Packages grouped by kind** (`internal/models`, `internal/database`). Rejected for the reasons
  in Context.
- **Modules importing each other directly.** Less code, but modules become coupled at compile time,
  tests need real dependencies such as a database, and splitting a module out is harder.
- **Table-name prefixes instead of schemas** (`rates_historical`). Simpler, but ownership is only a
  naming convention, and nothing prevents a cross-module join.
