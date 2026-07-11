---
name: feature
description: Use when starting any change to a system that is already running (new feature, behavior change, refactor, schema change, or dependency change). Walks through the architecture-change decision, safe data migration, and safe rollout. Invoke before writing code for any post-launch change.
---

# Adding a feature to a live system

Teach me through each gate. Do not skip to code. Default learning mode is A.

## 1. Frame it
- State the user-facing change in one sentence, and the smallest slice that delivers value.
- Find the seam it touches (which internal/ package, which route, which table). Read that area's CLAUDE.md.

## 2. Architecture-change gate (default: do NOT change architecture)
Changing structure is expensive and risky. Only change it when the CURRENT design cannot meet a REAL
requirement, not for taste. Ask me these and make me answer:
- Can this be done inside the existing package seams? If yes, do that.
- Is the pain a real limit (latency, correctness, scale we actually have) or hypothetical?
- What is the smallest change that fits the current architecture?
If a real architecture change is warranted: STOP and write an ADR in docs/decisions/NNNN-title.md
(Context, Options considered, Decision, Consequences). Get my approval before any code.
Then update docs/sadd.md if the design itself changed.

## 3. Data-migration gate (expand / contract; never a breaking change in one step)
If the schema changes, never rename or drop in the same deploy that ships the new code. Use expand/contract:
1. EXPAND: add new columns/tables as nullable/additive. Migration is forward-only; reversible by a new migration.
2. BACKFILL: populate new shape in a separate step; keep it idempotent and restartable.
3. DUAL-WRITE / DUAL-READ: new code writes both shapes and reads the new one behind a flag.
4. SWITCH: flip reads to the new shape once backfill is verified.
5. CONTRACT: in a LATER deploy, remove the old columns/dual-write.
Every migration gets a test that applies up and down against a real Postgres.

## 4. Rollout gate (keep the old path working the whole time)
- Keep the existing behavior intact until the new path is proven. Guard the new path with a feature flag.
- API changes are additive: add fields, do not remove or rename. If truly breaking, version the endpoint.
- Ship with the blue/green swap (rewrite upstream, health-check, reload). Verify on green before cutover.
- Have a rollback in mind before deploying: flag off, or swap back to blue.

## 5. Close the loop
- TDD throughout: failing test first at every step.
- Update docs/implementation-plan.md changes log, the ADR, and sadd.md/repo.md if they changed.
- Small atomic commits (skill: commit). One expand, one backfill, one switch, one contract: separate commits.