---
name: commit
description: Use when creating a git commit. Produces small atomic commits with Conventional Commits messages, one logical change each. Invoke after a test goes green or a refactor completes.
---

# Commit conventions

One commit per logical change. Do not batch unrelated changes.

Format: type(scope): summary
- types: feat, fix, refactor, test, docs, chore, build, ci
- scope: the area, e.g. backend, ingestion, frontend, infra, shared
- summary: imperative, lower-case, no trailing period, under ~72 chars
- body (optional): the WHY, wrapped at ~72 cols. Not what changed (the diff shows that).

Before committing:
1. Run the relevant tests and linter. Do not commit red or unformatted code.
2. git status and git diff. Show me what will be staged and confirm it is one logical change.
3. If the diff mixes concerns, split into separate commits.

Reference the plan step in the body when relevant, e.g. "Implements step 3.2".
Never commit secrets, .env files, or large binaries.