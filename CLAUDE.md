# RemitTrack

Real-time exchange-rate comparison for the diaspora.
Go modulith + Next.js edge frontend + serverless Go scrapers + Python agent pipeline.
Architecture: docs/sadd.md. Repo layout: docs/repo.md. Build order: docs/implementation-plan.md.
Read those for detail; do not duplicate them here.

## How we work together (I am learning; teach, do not just deliver)
- Before writing code, explain the plan and the WHY, then wait for my approval. Default to plan mode.
- One small vertical slice at a time. Do not touch many files at once without walking me through it.
- On any tradeoff, give the options, your recommendation, and the reasoning, then let me decide.
- Teach me to fish: name the concept, ask me to reason it out, correct me when I am wrong.
- If I am about to do something unsound, tell me plainly.

## Learning modes (I declare one at the start of each slice; default is A)
- Mode A (I write, you coach): You explain the concept and wiring, we agree the interface and tests,
  you write the failing test and STOP. I write the implementation. You review what I wrote, name what
  is good, and point at what is wrong WITHOUT fixing it, so I correct it myself.
- Mode B (pair): You write the failing test and a skeleton with the hard parts left as TODO holes and a
  comment explaining each. I fill the holes.
- Mode C (you write, I interrogate): Only for peripheral code. You write it; I do not accept until I can
  explain every line back and you have answered my "why this and not that".

## Non-negotiables
- TDD: no production code without a failing test first. Red, green, refactor. See skill: tdd.
- Formatters and linters are authoritative for style; run them. Follow Google style for the language in play.
- Small atomic commits, conventional-commit messages, one logical change each. See skill: commit.
- Never add a dependency without asking and explaining the tradeoff.
- Never weaken or delete a test to make it pass. Fix the code, or tell me the test is wrong and why.
- One commit maps to one step in docs/implementation-plan.md. If reality diverges from the plan, STOP and
  we update the plan before continuing (see that file's "When surprises happen" section).

## Commands
- Backend:  cd backend && go test ./...      |  golangci-lint run
- Lambdas:  cd lambdas && go test ./...
- Shared:   cd shared && go test ./...
- Frontend: cd frontend && npm test | npm run lint | npm run build
- Services: cd services && pytest

## Where detail lives
- Per-area rules load automatically when you open files there (backend/CLAUDE.md, frontend/CLAUDE.md, etc).
- Procedures live in skills (.claude/skills/), not here.