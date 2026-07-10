## RemitTrack: Technical System Architecture Design Document (SADD)
Target Phase: Production Blueprint (Hardened, Cost-Optimized)

Author: Oluwasegun Adewola

Last Updated: July 2026

------------------------------
## 1. Global System Topology
The global system topology isolates real-time data ingestion, public presentation caching, core business state management, and the agentic self-healing verification network into clean, decoupled tiers.
```
                                    [ 1. EDGE INGRESS LAYER ]
                                               │
               ┌───────────────────────────────┴───────────────────────────────┐
               ▼ (Public Web / Bot Traffic)                                    ▼ (B2B Enterprise Requests)
    +──────────────────────────────+                                +──────────────────────────────+
    │     AWS CloudFront (CDN)     │                                │     AWS CloudFront (CDN)     │
    │ - Applies WAF Rate-Limits    │                                │ - Restricts to /v1/b2b/*     │
    │ - Serves Static Edge ISR Pages│                               │ - Enforces Token Auth Header │
    +──────────────────────────────+                                +──────────────────────────────+
               │                                                               │
               └───────────────────────────────┬───────────────────────────────┘
                                               │ (Clear, Filtered HTTP Routes)
                                               ▼
                                   [ 2. COMPUTE HOST NODE ]
+─────────────────────────────────────────────────────────────────────────────────────────────────────────────+
│ AWS EC2 INSTANCE (t4g.micro • 1 GB RAM • Linux Ubuntu • Managed via systemd)                                │
│                                                                                                             │
│   +─────────────────────────────────────────────────────────────────────────────────────────────────────+   │
│   │                                      NGINX REVERSE PROXY                                            │   │
│   │   - Port 80/443 External Ingress Interception & TLS Termination                                     │   │
│   │   - Active Blue/Green Upstream Ports [8080 / 8081] Routing Matrix with Zero-Downtime Hot Reloads    │   │
│   +─────────────────────────────────────────────────────────────────────────────────────────────────────+   │
│                                              │                                                              │
│                                              ▼ (Local Socket Proxy Pass)                                    │
│   +─────────────────────────────────────────────────────────────────────────────────────────────────────+   │
│   │                                  GO APPLICATION MODULITH CORE                                       │   │
│   │                                                                                                     │   │
│   │   ├── [API Package]        : Consumer REST Endpoints & B2B Token-Gate Middleware                    │   │
│   │   ├── [Chat Package]       : Conversational LLM Context Router (Streaming Chunk Workers)            │   │
│   │   ├── [Alert Pool Worker]  : In-Memory Map Corridor Index Matcher (`map[string][]AlertConfig`)      │   │
│   │   └── [Ingestion Consumer] : Persistent SQS Long-Poll Loop (Idempotent Multi-Row DB Builder)        │   │
│   +─────────────────────────────────────────────────────────────────────────────────────────────────────+   │
│                                              │                                                              │
│                                              ▼ (Sub-Millisecond Thread Reading)                             │
│                                ==============================                                               │
│                                [   RATECACHE STRUCT INTERFACE ]                                             │
│                                ==============================                                               │
│                                      │                │                                                     │
│                                      ├─► [Lean]       ▼ [Enterprise Scale-Up]                               │
│                                      │   sync.Map     Amazon ElastiCache Redis Cluster                      │
+──────────────────────────────────────┼──────────────────────────────────────────────────────────────────────+
                                       │
                                       ▼ (Persistent Database Connections)
                                    [ 3. CORE COUPLING STORAGE LAYER ]
  +───────────────────────────────────────────────────────────────────────────────────────────────────────────+
  │ AWS RDS POSTGRESQL (db.t4g.micro • Multi-AZ Active-Passive Replication • PITR Logging • 20GB gp3 SSD)     │
  +───────────────────────────────────────────────────────────────────────────────────────────────────────────+
                                       │
                                       ▼ (Tails stdout Structured Streams)
                                    [ 4. ASYNCHRONOUS TELEMETRY LAYER ]
  +────────────────────────────────────────────────────────────────────────────────────────────────────────────+
  │ AWS CLOUDWATCH ENGINE                                                                                      │
  │  ├── Log Group Ingestion : Streams application JSON logs                                                   │
  │  └── Metric Filters      : Matches scraper error / empty-payload log patterns, firing Repository Dispatch webhooks │
  +────────────────────────────────────────────────────────────────────────────────────────────────────────────+
                                       │
                                       ▼ (Encrypted Webhook JSON Payload)
                                    [ 5. INTEGRATED PIPELINE HARNESS ]
  +───────────────────────────────────────────────────────────────────────────────────────────────────────────+
  │ GITHUB ACTIONS VM RUNNER (Ubuntu-Latest • Ephemeral Sandbox)                                              │
  │  ├── Agentic Core        : Multi-Turn LangGraph Orchestrator (RAG Vector S3 Unpacker)                     │
  │  ├── Security Guardian   : Abstract Syntax Tree (AST) Compiler Import Whitelist Sanitizer                 │
  │  ├── Simulation Sandbox  : Docker Compose Proxy Network (JA4 Handshake Testing Layouts)                   │
  │  └── Canary System       : Live Residential Proxy Smoke Tester & 3-Sigma Mathematical Verification        │
  +───────────────────────────────────────────────────────────────────────────────────────────────────────────+
                                       │
                                       ▼ (On Full Manual Verification Merge Approval)
                                    [ 6. SERVERLESS EXTRACTION CLUSTER ]
  +───────────────────────────────────────────────────────────────────────────────────────────────────────────+
  │ AWS LAMBDA INGESTION FUNCTIONS                                                                            │
  │  ├── Schedule Profile    : EventBridge Cron Invocations (Every 5 Minutes)                                 │
  │  ├── Network Mask        : uTLS Browser Signature Forging across Rotating Residential Proxies             │
  │  └── Buffer Channel      : Pushes Standardized JSON Objects directly to AWS SQS                           │
  +───────────────────────────────────────────────────────────────────────────────────────────────────────────+
```

------------------------------
## 2. Component Deep Dives & Architectural Zoom-Ins
### Section A: Frontend, Ingress, & Programmatic Edge Caching
### Architectural Design & Implementation Blueprint
The public delivery architecture completely avoids hosting a live, resource-intensive Node.js JavaScript rendering process on the production machine. Instead, it shifts the computing weight outward to the Content Delivery Network (CDN) edge.
```
[ AI Crawler / PerplexityBot ] ──► [ AWS CloudFront Edge Node ] ──► (Cache Hit: 0ms Host Compute)
                                              │
                                              ▼ (On Ingestion Completion / Cache Expired)
                                   [ Executing Edge ISR Trigger ]
                                              │
                                              ▼ (Purges & Rebuilds Specific Path JSON/HTML)
                                   [ Next.js Build Workspace (CI/CD) ] ──► Pushes Flat Files to S3
```

   1. **Static Build Ingestion Platform**: The frontend application package (`/apps/ssr-frontend`) is developed using Next.js. During continuous delivery deployment steps, it compiles directly into flat, static HTML and asset arrays.

   2. **CDN Distribution Layer**: These pre-compiled pages are hosted inside an **AWS S3 Bucket** and distributed using **AWS CloudFront**. When an AI search engine, scraping bot (e.g., PerplexityBot), or standard user requests content, CloudFront resolves the connection directly from edge cache memory, ensuring **0 MB memory allocation** on your core host machine.

   3. **Edge Incremental Static Regeneration (ISR) Webhook Integration**: When the ingestion pipeline finishes compiling fresh pricing values (every 5 minutes), the background processor fires an authenticated webhook execution request to the CloudFront Invalidation API. This purges and regenerates the specific currency corridor page (e.g., `/exchange-rates/gbp-to-ngn`) across edge caches. This pattern ensures that web crawlers that bypass local client-side Javascript execution always receive pre-rendered, real-time metrics.

   4. **Structured Semantic Injections**: The static file generation layer maps data directly into rich JSON-LD data graphs. It injects `ExchangeRateSpecification` schemas for table rows and parses user intent queries from conversational logs into valid `FAQPage` structures, capturing organic long-tail search visibility.

   5. **Internationalization Execution Architecture (i18n)**: Next.js handles regional localized paths natively through directory layout structures: `/app/[locale]/exchange-rates/[corridor]/page.tsx`. It pulls strings from disk dictionaries (e.g., `/locales/fr.json`), and the build process automatically appends explicit canonical tags and hreflang cross-references into the HTML header:
        ```html
        <link rel="alternate" hreflang="en-gb" href="https://remittrack.com" />
        <link rel="alternate" hreflang="fr-fr" href="https://remittrack.com" />
        ```
   
   
### Product Engineering & Operational Configurations

* **WAF Layer Guard Rules**: Set up AWS WAF rules on CloudFront with custom rate limits. Block or puzzle-challenge IP addresses hitting `/api/v1/auth/*` or dynamic conversational endpoints more than 100 times within 5 minutes. Maintain a dedicated header override block to guarantee entry for verified search crawlers like Googlebot and PerplexityBot.

### Architectural Trade-offs

* **Static Edge Generation over Live SSR Runtime**: Choosing edge-triggered static generation limits your ability to customize page components dynamically based on real-time user header variables (like local currency localization guessing via incoming IP data). However, it drops host operational server expenses down to zero and protects the compute engine from crash states during heavy web scraping sweeps.

------------------------------
## Section B: Core Backend Modulith, In-Memory Alert Engine, & B2B Layer
### Architectural Design & Implementation Blueprint
The core application (`/apps/remittrack-api`) is constructed inside a single Go compiled modulith binary package running under **Linux Systemd Process Supervision** on an **AWS EC2 `t4g.micro`** virtual node.
```
                      ┌───────────────────────────────────────────────┐
                      │    Go Application Modulith Process (Port 8080)│
                      │                                               │
  [ SQS Message Queue]├──► [ Ingestion Thread ]                       │
                      │          │                                    │
                      │          ▼ (Dispatches Async Internal Event)  │
                      │    [ Alert Pool Evaluation Loop ]             │
                      │          │                                    │
                      │          ▼ (Instant Look-up)                  │
                      │    Memory Corridor Index Map                  │
                      │    `map[string][]AlertConfig`                 │
                      └───────────────────────────────────────────────┘
```

   1. **NGINX Upstream Rolling Swap Routing**: NGINX handles ingress data flows on the virtual host machine, shielding the Go execution process. It implements a dual-port upstream configuration matrix to support zero-downtime hot reloads:
        ```nginx
        # Deploy tooling rewrites this file to point at the newly-started colour,
        # then runs `sudo systemctl reload nginx`. Only one colour is ever active.
        upstream go_modulith_active {
            server 127.0.0.1:8080;   # active colour (blue); swapped to :8081 (green) on deploy
        }
        ```
   
        When the deployment pipeline transfers a new binary artifact onto the host, it starts the new process on the idle colour (say `8081`), health-checks it, rewrites the `upstream` block to point at `8081`, and triggers a safe configuration refresh: `sudo systemctl reload nginx`. NGINX spawns fresh worker threads that route new connections to `8081` while the old workers let in-flight transactions on `8080` drain gracefully, after which the retired process is stopped. This keeps the swap free of dropped requests.
   
   2. **In-Memory Worker Alert Pool Loop**: To maintain performance, the platform separates rate alert processing from the database layer. The Modulith instantiates a background pool of worker goroutines linked to an optimized memory map index cache: `map[string][]AlertConfig`. The key maps to specific currency channels (e.g., `"USD:GHS"`). When a fresh data batch completes ingestion, the worker thread accesses the target array index instantly, evaluating price-target alerts in memory. If a threshold condition evaluates to true, the entry is dispatched to an asynchronous thread that triggers **AWS SES (Simple Email Service)** to notify the user.

   3. **Stripe Webhook Idempotency Architecture**: The Stripe billing endpoint (`/api/v1/webhooks/stripe`) enforces crypto-signature verification using Stripe's Webhook Secret Key. To block duplicate transaction processing streams from network retry events, it passes payloads through an **Idempotency Execution Ledger** table inside PostgreSQL. The controller hashes the transaction token (`evt_id`); if a collision matches an existing hash, the pipeline discards the request before running business logic.

   4. **B2B Token Verification Middleware Engine**: Corporate consumers pull real-time rates through an independent token-gated endpoint (`/api/v1/b2b/rates/latest`). The request passes through custom Go auth middleware that validates the provided `x-api-key` against hashed records inside PostgreSQL. Requests are rate-limited via a token-bucket algorithm hosted inside the instance memory space.

### Product Engineering & Operational Configurations

* **Systemd Configuration File (`/etc/systemd/system/remittrack.service`):**

    ```ini
    [Unit]
    Description=RemitTrack Go Modulith Core Engine
    After=network.target

    [Service]
    Type=simple
    User=ubuntu
    WorkingDirectory=/home/ubuntu/remittrack
    ExecStart=/home/ubuntu/remittrack/main
    Restart=always
    RestartSec=1s
    Environment=APP_ENV=production REMITTRACK_INFRA_PROFILE=LEAN
    StandardOutput=append:/var/log/remittrack/app.log
    StandardError=append:/var/log/remittrack/err.log

    [Install]
    WantedBy=multi-user.target
    ```

### Architectural Trade-offs

* **In-Memory Map Cache over Distributed Redis Infrastructure**: Running alert caches natively inside a Go `sync.Map` or local pointer allocation layer delivers incredibly low latency at zero cost. However, it binds your state memory to a single host computer footprint. If you expand the deployment cluster to support horizontal multi-node scaling across separate machines, this in-memory cache architecture must be refactored to use an external, shared distributed cache layer like **Amazon ElastiCache Redis**.

------------------------------
## Section C: Hardened Serverless Ingestion & Network Evasion Data Pipeline
### Architectural Design & Implementation Blueprint
The scraping and pipeline framework decouples data extraction completely from the host machine to isolate process threads and prevent connection exhaustion.
```
[ AWS EventBridge Cron ] ──► [ AWS Lambda Ingestion Runner ]
                                        │
                                        ▼ (uTLS Fingerprint Handshake)
                            [ Rotating Residential Proxy Node ]
                                        │
                                        ▼ (Fetches payload anonymously)
                            [ Target Money Transfer Provider ]
                                        │
                                        ▼ (Maps to Normalized RateRecord)
                            [ Standard AWS SQS Message Queue ]
                                        │
                                        ▼ (Long-Polls concurrent batches)
                            [ Core Modulith Ingestion Processor ]
                                        │
                                        ▼ (Idempotent Multi-Row UPSERT Block)
                            [ AWS RDS PostgreSQL Storage Engine ]
```

   1. **Stateless Serverless Execution Extractor**: Scraper code elements are extracted out of the core app process space into standalone **AWS Lambda functions** written in Go. An AWS EventBridge cron schedule executes the Lambda cluster every 5 minutes.

   2. **Network Evasion & Handshake Forging**: To bypass firewall systems (such as Cloudflare or Akamai), the scraping Lambdas implement the uTLS library (`github.com/refraction-networking/utls`) instead of Go's default `crypto/tls` module, allowing them to forge client TLS extension handshakes to look exactly like standard consumer browsers. The Lambda configuration routes all outbound connections through an external, encrypted **Rotating Residential Proxy Network Service** by setting the HTTP_PROXY shell variable, hiding data center IP signatures.

   3. **The Standard SQS Buffering Infrastructure**: The Lambda normalizes scraped rates into a unified JSON model at the edge, pushes the payload into a standard **AWS SQS Ingestion Queue**, and instantly terminates. This decouples the scrapers from the database layer. Inside the core Go modulith on the EC2 host, a persistent background consumer thread long-polls the SQS queue, pooling rows into unified batches.

   4. **Idempotent Storage Ingestion Layer**: The Modulith writes data batches to the **AWS RDS PostgreSQL** (`db.t4g.micro`) engine using an idempotent insert-or-ignore block tied to a unique constraint array:
        ```postgresql
        -- scrape_timestamp is stored at full resolution, so distinct fetches are never
        -- collapsed; DO NOTHING only discards true duplicate deliveries (SQS at-least-once).
        INSERT INTO historical_rates (provider, currency_pair, mid_rate, scrape_timestamp)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (provider, currency_pair, scrape_timestamp) DO NOTHING;
        ```
        A separate `current_rates` table holding one row per `(provider, currency_pair)` is updated on each successful insert and backs the comparison read path, so the growing `historical_rates` table is never scanned on hot user queries.
   
   5. **Multi-AZ Availability Durability Profile**: The RDS instance is configured within an active-passive **AWS Multi-AZ Deployment Map**. Transaction logs stream continuously to an isolated secondary availability zone, providing automated failover routing, automated snapshot states, and Point-In-Time Recovery (PITR) protection.

### Product Engineering & Operational Configurations

* **Outbound Redirection Protection**: To secure affiliate links from open-redirect phishing exploitation, the internal redirection utility `/out/` checks incoming link addresses against a strict domain whitelist before returning a `302 Found` response:
    ```go
    var ApprovedProviders = map[string]string{
        "lemfi":    "https://example.com",
        "sendwave": "https://example.com",
    }
    ```


### Architectural Trade-offs

* **Standard SQS Ingestion over SQS FIFO Queues**: Switching from a FIFO queue to a standard SQS queue eliminates strict message delivery order sorting at the queue level. However, it removes throughput caps and simplifies scaling, while relying on idempotent SQL database constraints to safely handle out-of-order or duplicate message delivery streams.

------------------------------
## Section D: DevOps, CI/CD, & Automated Infrastructure
### Architectural Design & Implementation Blueprint
The platform uses a complete GitOps model to manage deployment automation and cloud resource provisioning through a single GitHub monorepo configuration.
```
[ Code Commit Push / Merge ] ──► [ GitHub Actions Global Pipeline Runner ]
                                              │
             ┌────────────────────────────────┴────────────────────────────────┐
             ▼ (Validation Track)                                              ▼ (Infrastructure Automation Track)
[ Run Golang linter / Unit Tests ]                                [ Initialize Terraform Workspace ]
             │                                                                 │
             ▼                                                                 ▼
[ Launch Docker Compose Chaos Sandbox ]                           [ Run terraform plan / apply ]
             │                                                                 │
             ▼ (Executes Outlier Canary Check)                                 ▼
[ Cross-Compile Lambda Binaries & Zip packages ]                 [ Provision Multi-AZ RDS & SQS Pools ]
             │                                                                 │
             └────────────────────────────────┬────────────────────────────────┘
                                              ▼
                        [ Execute Production Rolling Swap Live Update ]
```

   1. **The Terraform Infrastructure Automation Engine**: The cloud footprint is fully managed as code within the `/terraform/` repository directory. When modifications pass verification tests, continuous integration systems call the deployment toolchain:

        ```bash
        cd terraform
        terraform init
        terraform plan -out=tfplan          # posted to the PR for human review
        # apply runs only after manual approval on the protected prod environment
        terraform apply tfplan
        ```

        This automatically manages access rights, configurations, and scaling parameters for components like your SQS message queues, Lambda runtimes, and multi-AZ RDS storage blocks. Stateful resources, the multi-AZ RDS instance above all, carry a `lifecycle { prevent_destroy = true }` guard so a bad plan can never replace or delete them, and the production `apply` is gated behind a manual approval rather than auto-approved.

   2. **Unified Continuous Integration Verification Engine**: When changes merge into the main application branch, a GitHub Actions runner boots up an ephemeral Ubuntu-latest workspace environment. The workflow executes all standard internal package validations: `go test -v ./...`.

   3. **Production Rolling Server Deployments**: If the global validation matrix returns full success checks, the runner cross-compiles the Go API Modulith core executable, packages the Lambda files into compressed zip folders, and ships the artifacts to AWS using the AWS CLI. It then securely connects to the EC2 host via SSH to reload NGINX and execute the zero-downtime rolling update process.

### Product Engineering & Operational Configurations

* **Terraform SQS Configuration Block (`terraform/sqs.tf)`**:
```hcl
resource "aws_sqs_queue" "remittrack_ingestion_queue" {
  name                      = "remittrack-rate-ingestion"
  delay_seconds             = 0
  max_message_size          = 262144
  message_retention_seconds = 86400
  receive_wait_time_seconds = 20 # Enables long-polling to minimize billing overhead

  lifecycle {
    prevent_destroy = true # applied to all stateful resources, RDS above all
  }
}
```

### Architectural Trade-offs

* **GitHub Actions Runners over Dedicated Cloud Testing Environments**: Offloading chaos tests and compilation suites entirely to GitHub Actions introduces workflow queue latency compared to running a dedicated cloud container node 24/7. However, it eliminates recurring testing infrastructure expenses, leveraging GitHub's free computation allotments to keep operational overhead at zero.

------------------------------
## Section E: AI Engineering & Multi-Agent Oversight
### Architectural Design & Implementation Blueprint
The AI Engineering layer operates within an isolated sandbox environment inside the GitHub Actions pipeline (`/services/`), running a multi-agent system designed to manage code healing, legal compliance tracking, and automated platform onboarding safely.

```
                       [ ACTION WEBHOOK TRIGGER ]
                                   │
                                   ▼
                [ Ephemeral GitHub Actions Runner VM ]
                                   │
      ┌────────────────────────────┼────────────────────────────┐
      ▼                            ▼                            ▼
[ 1. SELF-HEALING AGENT ]    [ 2. COMPLIANCE AGENT ]      [ 3. ONBOARDING AGENT ]
  ├── Pulls broken code        ├── Crawls robots.txt        ├── Uses playwright-stealth
  ├── Queries ChromaDB RAG     ├── LLM Policy Check         ├── Sniffs background XHR
  └── Runs AST Linter Gate     └── Toggles provider_map.json└── Generates Go Lambda Code
```

#### 1. The Agentic Self-Healing Engine (Python + LangGraph)
When an active data extraction pipeline triggers an execution crash or empty JSON data array, an AWS alarm hits the GitHub workflow interface via an authenticated repository dispatch webhook.

* **The Context Execution Step**: The runner downloads the compressed vector memory workspace archive (`chroma_db.tar.gz`) from an isolated AWS S3 Bucket into local storage (`/tmp/chroma_db`) and initializes an in-memory client session using `chromadb.PersistentClient`. The agent runs semantic search queries across historical logs to retrieve the top 3 code patches that successfully resolved similar error profiles.

* **The AST Whitelist Linter Security Gate**: Code candidate modifications generated by the LLM (e.g., Claude 3.5 Sonnet) are parsed by a dedicated Go AST validator (Go's `go/parser` package, invoked as a subprocess by the Python agent) before compilation to block remote code execution (RCE) vulnerabilities. The linter enforces a strict import allowlist: **anything not on the allowlist of pure data-mapping packages (`strings`, `strconv`, `encoding/json`, `errors`) is instantly discarded, and `//go:linkname` compiler directives are rejected outright**. The engine's scope is strictly restricted to updating JSON field-mapping dictionary keys. If a legacy HTML DOM page structure fails, the agent stops and alerts a human engineer.
* **The Statistical Canary Cross-Validation Gate**: If the code compiles safely, it runs inside the localized Docker Compose simulation proxy network. If simulation checks pass, the runner performs a **Live Internet Smoke Test** across a residential proxy. The extracted rate payload must pass a statistical evaluation before deployment approval:

$$\text{Pass Constraint} = \mu - 3\sigma \le \text{Candidate Rate} \le \mu + 3\sigma$$

* The rate must also align with an adaptive, provider-calibrated variance corridor compared to the other active platforms to catch data correctness or inversion errors. For a brand-new provider or corridor with insufficient history to compute a stable mean and standard deviation, the Z-score check is skipped and the candidate is held for mandatory human review rather than auto-passed. If all validation criteria pass, the system logs the entry to the registry and **opens a Git Pull Request for manual human engineering review**.

#### 2. The Compliance Sentinel Agent
This agent runs daily via a GitHub Actions cron schedule. It downloads and parses the `robots.txt` and Terms of Service (ToS) files for all 15+ monitored platforms, saving historical snapshots in an S3 bucket as an audit trail. The agent calculates a SHA256 checksum against yesterday's copy. If a change is detected, an LLM evaluates the updated text. If it uncovers new scraping restrictions, the agent automatically opens a Git Pull Request to update `provider_map.json` and deactivate that scraper function in production.

#### 3. The Automated Onboarding Agent
Triggered via an input command (`make onboard URL=...`), this agent first checks the url parameters against the Compliance Sentinel rules engine. If compliant, it launches a headless browser session inside GitHub Actions using `playwright-stealth` to bypass automated traffic filters. Playwright binds an event listener to the browser's virtual network interface and records background traffic for 45 seconds. It filters responses using regex expressions like `.*(rate|exchange|pricing).*`. If it captures an internal JSON response containing currency data, it maps the structure and prompts the LLM to write a brand new Go-based Lambda scraping function, opening a pull request for manual merge.

### Product Engineering & Operational Configurations
* **Go AST Linter Component** `(/services/agent_healer/linter.go)`. The Python LangGraph agent shells out to this Go validator, because Go source can only be parsed by Go's own `go/parser`, not Python's `ast`:
```go
// Validates an LLM-generated Go patch using Go's own AST parser.
// Enforces a strict import ALLOWLIST: anything not explicitly permitted is rejected.
package main

import (
    "errors"
    "fmt"
    "go/parser"
    "go/token"
    "strconv"
    "strings"
)

// Only pure data-mapping libraries are permitted. Everything else is rejected.
var allowedImports = map[string]bool{
    "strings":       true,
    "strconv":       true,
    "encoding/json": true,
    "errors":        true,
}

func VerifyASTSafety(goSource string) error {
    // Reject compiler directives that reach runtime symbols without an import.
    if strings.Contains(goSource, "//go:linkname") {
        return errors.New("security alert: //go:linkname directive is not permitted")
    }
    fset := token.NewFileSet()
    file, err := parser.ParseFile(fset, "patch.go", goSource, parser.ImportsOnly)
    if err != nil {
        return fmt.Errorf("patch rejected: unparseable Go source: %w", err)
    }
    for _, imp := range file.Imports {
        path, _ := strconv.Unquote(imp.Path.Value)
        if !allowedImports[path] {
            return fmt.Errorf("security alert: forbidden import %q (not on allowlist)", path)
        }
    }
    return nil
}
```

### Architectural Trade-offs

* **Strict AST Code Generation Restrictions over Broad Agent Autonomy**: Restricting your automated code repair engine to basic JSON mapping updates prevents it from fixing structural HTML changes on legacy platforms. However, as one defense-in-depth layer ahead of the mandatory human merge, it materially reduces the attack surface against compiler bypass vectors like `//go:linkname` and resource-exhaustion loops.

------------------------------
## 3. Runtime Verification & Failure Registries
To track pipeline health and maintain high data accuracy across all currency corridors, the system records and verifies every automation event using a strict logging and testing matrix.

### The System Log Ledger Schema (`failure_registry.json`)
Every automated execution step, simulation check pass, and canary validation event is logged directly inside an **AWS S3 JSON Analytics Registry** file for regression verification tracking:

```json
{
  "case_id": "ERR-LEMFI-JSON-073",
  "timestamp": 1783635600,
  "provider": "LemFi",
  "error_profile": "KeyError: 'mid_market_rate_value' payload layout missing",
  "agent_execution_metrics": {
    "turns_attempted": 2,
    "tokens_consumed": 4350,
    "ast_linter_checks": "PASSED"
  },
  "validation_canary_telemetry": {
    "historical_z_score": 0.42,
    "cross_provider_variance_delta": "0.012",
    "canary_exit_status": "VERIFIED_SUCCESSFUL"
  },
  "patch_output_artifact_s3": "s3://remittrack-artifacts/patches/lemfi_fix_073.go"
}
```

### Continuous Integration Regression Testing Sequence
When a new code change or onboarding integration request drops into the monorepo branch, the global evaluation test suite runs automated validation checks to prevent regressions before deployment:
```
                  [ Pull Request Event / Auto-Heal Test Init ]
                                       │
                                       ▼
                     [ Boot Docker Compose Proxy Sandbox ]
                                       │
                                       ▼
                     [ Read failure_registry.json Targets ]
                                       │
                                       ▼
                     [ Execute Automated Regression Runs ]
       ├── Mock Proxy simulates each historical failure case sequentially
       └── Harness asserts new scraper code passes without breaking older rules
                                       │
                    ┌──────────────────┴──────────────────┐
                    ▼ (Any test fails)                    ▼ (All tests pass)
            [ REJECT BUILD ]                      [ APPROVE & PUSH TO MAIN ]
```

By verifying your entire historical failure log array against every code candidate inside an isolated, containerized local loop, you protect production pipelines from breaking changes. This ensures high data reliability, solidifies security perimeters, and provides an exceptionally robust system architecture.

