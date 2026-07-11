# tests/ (cross-cutting integration + regression)

Unit tests live next to the code they test. This dir is for integration and the regression/canary harness only.

- regression/ holds docker-compose.yml (mock proxy network) and failure_registry.json (historical cases).
- Integration tests spin up real Postgres and a mock SQS via docker-compose; assert wiring, not units.
- Every historical failure in failure_registry.json must have a replay test that stays green.