# RemitTrack developer tasks.
#
# Every target operates per-module because the repo root is not itself a Go
# module — the workspace is defined in go.work. The test/lint/build targets
# skip modules with no .go files: `go test ./...` and golangci-lint both exit
# non-zero on an empty module, so skipping keeps the empty workspace green and
# self-heals as packages land. Run `make <target>`.
MODULES := shared backend lambdas
FUZZTIME ?= 30s

.PHONY: test lint build fuzz up down

## test: run all tests in every non-empty module
test:
	@for m in $(MODULES); do \
		if find $$m -name '*.go' | grep -q .; then \
			echo "==> test $$m"; \
			(cd $$m && go test ./...) || exit 1; \
		else \
			echo "==> skip $$m (no .go files)"; \
		fi; \
	done

## lint: golangci-lint every non-empty module
lint:
	@for m in $(MODULES); do \
		if find $$m -name '*.go' | grep -q .; then \
			echo "==> lint $$m"; \
			(cd $$m && golangci-lint run ./...) || exit 1; \
		else \
			echo "==> skip $$m (no .go files)"; \
		fi; \
	done

## build: compile every non-empty module
build:
	@for m in $(MODULES); do \
		if find $$m -name '*.go' | grep -q .; then \
			echo "==> build $$m"; \
			(cd $$m && go build ./...) || exit 1; \
		else \
			echo "==> skip $$m (no .go files)"; \
		fi; \
	done

## fuzz: run every fuzz target in shared for FUZZTIME each (default 30s).
## Go fuzzes one target per run, so targets are listed and run in turn.
## Seed inputs already run in `make test`; this explores beyond them.
fuzz:
	@cd shared && for pkg in $$(go list ./...); do \
		for f in $$(go test -list '^Fuzz' $$pkg | grep '^Fuzz'); do \
			echo "==> fuzz $$pkg $$f"; \
			go test $$pkg -run '^$$' -fuzz "^$$f\$$" -fuzztime $(FUZZTIME) || exit 1; \
		done; \
	done

## up: start local dependencies (Postgres) in the background
up:
	docker compose up -d

## down: stop and remove local dependencies
down:
	docker compose down
