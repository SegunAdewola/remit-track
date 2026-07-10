## RemitTrack
Exchange Rate Comparison Platform • Jun. 2025 – Ongoing

1. Overview
2. System Architecture & Boundaries
3. Data Pipeline & Agentic Systems
4. Key Features
5. Engineering Trade-offs

------------------------------
## 1. Overview
RemitTrack is a real-time exchange rate comparison platform built for the African diaspora. Sending money home is something millions of people do every month, and the difference in rates across providers can be significant. RemitTrack pulls live rates from 15+ Money Transfer Providers (MTPs), normalises them into a consistent format, and surfaces them in one place so users can see where they get the best deal.
Beyond the consumer-facing side, the aggregated rate data serves as a market intelligence and arbitrage risk signal for MTP operators. The platform also has an LLM-powered chat interface that lets users ask questions about rates in plain English, built as an accessibility feature for users who find data tables harder to parse.

To address layout changes and anti-bot systems securely and legally, the infrastructure includes a decoupled multi-agent ecosystem. A Compliance Sentinel Agent monitors provider legal updates to enforce data compliance. An Automated Onboarding Agent uses stealth network sniffing to discover hidden JSON backend endpoints. When live extractions break, an Agentic Self-Healing Data Pipeline attempts up to three code repair loops on JSON field-mapping transformers. It validates code candidates through a GitHub-Actions-hosted Chaos Engineering Simulation Environment and a statistical cross-validation canary check before opening a Git Pull Request for manual human engineering review.

To turn this high-performance utility into a sustainable business engine, the architecture features tier-gated user personalization profiles, a Stripe-powered automated subscription verification layer, and a programmatic affiliate outbound routing module. Combined with a secure B2B Enterprise Data API layer and a dedicated Market Intelligence Dashboard for MTP operators, the platform supports multi-channel monetization. Localized multi-region hreflang structural configurations automatically translate dynamic currency matrices to capture global financial search volumes and conversational AI engine citations.

------------------------------
## 2. System Architecture & Boundaries
The core production platform is written in Go and structured as an optimized modulith hosted on AWS. It implements clean interface abstractions to shift smoothly from a resource-efficient local footprint to a high-availability enterprise configuration.

## Full Production, Monetization, B2B Data Feeds, & Multi-Agent Network Topology
```
                                      [ THE INGRESS EDGE ]
                       +───────────────────────────────────────────────+
                       │   Client Browser / AI Bot / Corporate API B2B  │
                       +───────────────────────────────────────────────+
                                               │
                                               ▼ (HTTPS / TLS via Geo-i18n & /b2b Routes)
                               +──────────────────────────────+
                               │     AWS CloudFront (CDN)     │
                               │ - Applies: AWS WAF Block     │ <── Rate Limits Bot Spam
                               │ - Serves: Pre-Rendered Pages │ <── Edge ISR Frontend Layer (0 MB Host RAM)
                               +──────────────────────────────+
                                               │
                                               ▼ (HTTP / Port 80)
+─────────────────────────────────────────────────────────────────────────────────────────────────────────────+
│ AWS EC2 INSTANCE (t4g.micro Modulith Host Node - 1 GB RAM Allocated)                                        │
│                                                                                                             │
│   +─────────────────────────────────────────────────────────────────────────────────────────────────────+   │
│   |                                    NGINX Reverse Proxy Layer                                        |   │
│   |        - Blue/Green Upstream Ports [8080 / 8081] Configuration Matrix for Zero-Downtime Swaps       |   │
│   +─────────────────────────────────────────────────────────────────────────────────────────────────────+   │
│                                              │                                                              │
│                                              ▼ (Localhost Ports 8080 / 8081 Rerouting)                      │
│   +─────────────────────────────────────────────────────────────────────────────────────────────────────+   │
│   |                                   GO APP MODULITH CORE (Systemd Key)                                |   │
│   |                                                                                                     |   │
│   |   - Go HTTP REST API & Core Routing Logic                                                    |   │
│   |   - In-Memory Worker Alert Pool Loop (Evaluates targets via specific corridor index maps)           |   │
│   |   - B2B API Token Validation Middleware & Stripe Webhook Idempotency Check Controllers              |   │
│   |   - Background SQS Long-Polling Consumer (Pools writes into single DB transactions)                 |   │
│   +─────────────────────────────────────────────────────────────────────────────────────────────────────+   │
│                                              │                                                              │
│                                              ▼                                                              │
│                                ==============================                                               │
│                                [   CACHE LAYER INTERFACE    ]                                               │
│                                ==============================                                               │
│                                      │                │                                                     │
│                                      ├─► [Lean]       ▼                                                     │
│                                      │   sync.Map ────┴─────────────────────────────────────────────────────┤
│                                      ▼                                                                      │
+──────────────────────────────────────┼──────────────────────────────────────────────────────────────────────+
                                       │
                                       ▼
  +───────────────────────────────────────────────────────────────────────────────────────────────────────────+
  │                      AWS RDS PostgreSQL (db.t4g.micro - Multi-AZ High-Availability Map)                   │
  +───────────────────────────────────────────────────────────────────────────────────────────────────────────+
                                                     │
                                                     ▼ (Zerolog Structured JSON Stream)
                                         [ CloudWatch Log Groups ]
                                                     │
                                                     ▼ (CloudWatch Metric Filter Trigger)
                                    ====================================
                                    [  GITHUB ACTIONS PIPELINE VM      ]
                                    ====================================
                                    │ - Isolated LangGraph Agent Pool  │
                                    │ - Go AST Lexical Whitelist Linter│
                                    │ - Statistical Outlier Canary     │
                                    +──────────────────────────────────+
                                                     │
                                      (Success: Code branch PR generated)
                                                     ▼
                                    +──────────────────────────────────+
                                    │  Serverless AWS Lambda Clusters  │
                                    │  - Buffered via Standard AWS SQS │
                                    │  - Normalized Residential Proxies│
                                    +──────────────────────────────────+

```

## Abstract Ingress & Routing Interface

* Lean Profile: Edge connection traffic hits AWS CloudFront coupled to AWS WAF. WAF enforces custom rate limits to mitigate scraping from competitors while allowing verified search crawlers through. Traffic routes over HTTP to an NGINX Reverse Proxy running on an AWS EC2 t4g.micro instance (1 GB RAM). NGINX utilizes an Upstream Rolling Matrix across twin ports (8080 Blue / 8081 Green). During core application updates, the pipeline reloads NGINX, switching active ports sequentially to enable updates with zero dropped user requests.
* Edge Incremental Static Regeneration (ISR): To keep crawler-visible rate metrics real-time without adding host server RAM overhead, the system offloads page assembly to an Edge ISR Framework. When the serverless ingestion pipeline pushes a fresh data batch (every 5 minutes), it fires an automated revalidation trigger that runs the render inside a Lambda@Edge function, not on the host EC2 box. This rebuilds and caches flat HTML/JSON corridor files on edge CDN storage nodes, keeping host application RAM use at 0 MB.
* Enterprise Scale-Up Profile: CloudFront re-points its target origin to an AWS Application Load Balancer (ALB). The ALB distributes incoming loads across an AWS ECS Fargate Cluster that auto-scales Go Modulith containers horizontally based on real-time CPU thresholds.

## Abstract Cache Interface

* Lean Profile: The Modulith manages temporary rate states through a RateCache interface backed by a thread-safe Go sync.Map. This stores real-time rates directly in the EC2 instance's RAM, ensuring sub-millisecond response times at zero infrastructure cost.
* Enterprise Scale-Up Profile: If the single-node environment is scaled out horizontally, the RateCache provider switches via an environment variable to an Amazon ElastiCache Redis cluster. This acts as a centralized state machine and avoids cache drift across distributed nodes.

------------------------------
## 3. Data Pipeline & Agentic Systems
Data collection is entirely decoupled from the primary application process, executing inside standalone, serverless AWS Lambda functions to prevent resource starvation or runtime leaks within the Modulith API.

```
[ EventBridge Cron ] ──► [ Lambda Scrapers ] ──► (Residential Proxy) ──► [ Outbound MTP Targets ]
                                 │
                                 ▼ (Standardized RateRecord Mapping)
                        [ AWS SQS Ingestion Queue ]
                                 │
                                 ▼ (Long-Poll Consumer Thread)
                     [ Go Modulith Single DB Client Connection ]
                                 │
                                 ▼ (Idempotent DB UPSERT Multi-Row Transactions)
                     [ AWS RDS PostgreSQL Multi-AZ Engine ]
```

* Extract: Ingestion tasks execute inside individual AWS Lambda functions managed by an AWS EventBridge scheduler. To bypass anti-bot defenses, the Go-based Lambda clients drop the standard network client stack, use the uTLS library (github.com/refraction-networking/utls) to forge the TLS signatures of major consumer browsers, and route all outbound calls through a Rotating Residential Proxy Network to hide data center IP signatures.
* Transform & Load: Scrapers normalize data into a unified RateRecord model at the edge, pushing payloads to an IngestionStream interface. The lean profile routes records to a standard AWS SQS Ingestion Queue. Inside the core Go modulith, a background consumer thread long-polls SQS, combining entries into efficient multi-row queries. Data deduplication and message sorting are handled at the database storage layer using an Idempotent UPSERT execution block bound to a unique constraint key array: UNIQUE(provider, currency_pair, scrape_timestamp).
* Storage Source of Truth: An AWS RDS PostgreSQL (db.t4g.micro) instance operates within a Multi-AZ Deployment Configuration, providing automated failover routing, continuous data replication, and point-in-time recovery (PITR) protection. This architecture choice intentionally doubles the core database billing footprint to safeguard paid customer subscription transactions and protect B2B data integrity.

## The Agent Ecosystem & Self-Healing Pipeline
The platform implements three specialized, event-driven agents inside the monorepo's /services/ directory:
## A. The Agentic Self-Healing Engine (Python + LangGraph)
When a scraper fails, it emits a structured JSON string via Zerolog. A CloudWatch Metric Filter alarm captures this event and sends a webhook to fire a GitHub Actions runner.

   1. RAG Context Search: The runner downloads a compressed vector archive (chroma_db.tar.gz) from AWS S3 into local storage (/tmp/chroma_db). The Python agent uses semantic search to locate the three most relevant historical code fixes for the encountered error string.
   2. AST Lexical Linter Security Gate: To protect against Remote Code Execution (RCE) via indirect prompt injection from compromised HTML pages, the generated Go code passes through a strict Abstract Syntax Tree (AST) validation parser before compilation. The linter enforces a strict import whitelist: the code is instantly rejected if it attempts to load unsafe packages like os, exec, net, unsafe, or reflect. It can only use pure data formatting libraries (like strings or encoding/json). The agent's scope is strictly limited to updating JSON field-mapping transformers for providers that expose JSON endpoints. If a legacy HTML DOM path breaks, the agent automatically escalates the issue to a human engineer.
   3. Statistical Cross-Validation Canary Gate: The runner compiles the agent's code candidate and evaluates it within the Chaos Proxy simulation network. If successful, it runs a live-internet smoke test through a residential proxy node. The extracted payload passes through a Statistical Cross-Validation Matrix, asserting that the calculated rate lands within a 3-sigma ($3\sigma$) historical Z-score window and falls within an adaptive, provider-calibrated variance corridor compared to the other active platforms. If validation checks pass, the engine updates the vector storage, appends the metadata to the registry, and opens a Git Pull Request for manual engineering sign-off and human deployment approval.

## B. The Compliance Sentinel Agent
This agent runs daily via a GitHub Actions cron schedule. It downloads and parses the robots.txt and Terms of Service (ToS) pages for all 15+ monitored platforms, storing historical snapshots in an S3-backed compliance directory as an audit trail. The agent calculates text checksums against previous versions. If a change is detected, an LLM evaluates the text. If the agent uncovers explicit scraping restrictions, it opens a Git Pull Request modifying provider_map.json to disable that specific scraper automatically.

## C. The Automated Onboarding Agent
Triggered manually or via a workflow argument (make onboard URL=...), this agent validates target links against the Compliance Sentinel rules engine to ensure data collection is legally permissible. It then spins up a headless browser session inside GitHub Actions using playwright-stealth to bypass automated traffic filters. Playwright binds an event listener to the browser's virtual network interface and records background traffic for 45 seconds. It filters responses using regex expressions like .*(rate|exchange|pricing).*. If it captures an internal JSON response containing currency data, it maps the structure and prompts the LLM to write a brand new Go-based Lambda scraping function, opening a pull request for manual merge.

------------------------------
## 4. Key Features

* Serverless Edge Ingestion: Offloads scraping workloads to standalone, auto-scaling AWS Lambda functions to eliminate process strain on the core application.
* JA4 Fingerprint Emulation: Forges transport-layer signatures and manages browser HTTP/2 frame properties via the uTLS library (github.com/refraction-networking/utls) to bypass advanced anti-bot systems.
* AST Security Linting: Validates agent-generated code structures using Abstract Syntax Trees as an initial defense-in-depth security layer ahead of human engineering review.
* API-First Fallback Ingestion: Prioritizes official partner API endpoints, falling back to web scraping only when compliant with a provider's legal policies.
* Edge Programmatic ISR Optimization: Injects real-time ExchangeRateSpecification JSON-LD markups into flat, automated edge-cached HTML files to optimize visibility for AI search engines without adding host server RAM overhead.
* Anonymised FAQ Mining Engine: Aggregates and filters conversational data logs from the LLM chat assistant, publishing structured FAQ schemas to capture search traffic.
* Regression-Checked Simulation Sandbox: Validates candidate scraper fixes against historical failure profiles inside an isolated Docker Compose network before deployment.
* In-Memory Worker Alert Pool Engine: Evaluates real-time rate ingestion matches using a background worker loop, distributing user notifications via a fast index map without causing database lockups.
* Masked Outbound Redirection: Routes user click traffic through an internal proxy endpoint (/out/*) matching a Strict Whitelist Array of approved provider domains to eliminate open-redirect phishing vectors.
* Stripe Webhook Idempotency Controls: Secures payment pipelines using cryptographic signature checks and an evt_id database ledger to block duplicate transaction streams.
* Corporate B2B API Access Layer: Exposes a high-throughput, token-authenticated JSON feed (/v1/b2b/*) that enterprise consumers can query to pull real-time exchange rates.
* Operator Market Intelligence Dashboard: Provides an isolated dashboard interface for MTP partners to track pricing trends, regional market shares, and potential arbitrage exposure.

------------------------------
## 5. Engineering Trade-offs

* Statically Compiled Lambdas over Dynamic Plugins: Moving data collection to AWS Lambda functions increases deployment footprints. However, it completely avoids the memory leak and compilation constraints of Go's native plugin system, allowing zero-downtime updates without restarting the core server process.
* File-Persisted S3 Vector Archives over Hosted Instances: Downloading and extracting a compressed vector database archive (chroma_db.tar.gz) inside an ephemeral GitHub Actions container adds minor startup latency. However, it avoids the ongoing subscription costs of an always-on cloud database instance, dropping vector storage costs to near $0.00.
* GitHub Actions Testing over Dedicated Test Infrastructure: Evaluating agents within GitHub Actions workflows introduces queue delays compared to keeping a persistent cloud testing cluster. However, it utilizes the runner's free pre-installed Docker Compose and compute limits, reducing testing infrastructure costs to zero.
* Edge-Triggered ISR over Live Server-Side Rendering (SSR): Re-rendering flat HTML/JSON corridors directly on CloudFront edge nodes during ingestion updates limits real-time user-context customization on the fly. However, it completely removes the live Node process footprint from the host server, preventing OOM conditions and maintaining high resource efficiency.
* Defense-In-Depth AST Linting over Automated Production Deploys: Restricting the agent's code generation scope to JSON structures and forcing all candidate patches through a manual human engineering code merge adds operational overhead to the pipeline. However, it protects against advanced compiler bypass vectors (like //go:linkname) and resource-exhaustion code loops, materially reducing the attack surface.
* In-Memory Event Evaluation over Direct Database Querying: Processing user threshold triggers in memory requires caching active alert records in RAM. However, it prevents high-frequency read spikes and lock queues on your historical PostgreSQL storage tables during ingestion bursts.
* Standard SQS Queues over internal Memory Channels: Introducing an external messaging queue adds configuration complexity. However, it provides a durable, asynchronous buffer layer that protects the database from connection exhaustion, while relying on idempotent database constraints to handle message sorting cleanly.

------------------------------
© 2026 Oluwasegun Adewola. All rights reserved.
