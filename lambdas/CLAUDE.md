# lambdas/ (stateless Go scrapers)

Each scraper is an independent deployment archive. Depends on ../shared, never on backend/internal.

- Use shared/scrapeharness for the uTLS transport and proxy rotation. Do not hand-roll TLS per scraper.
- Normalize every provider into shared/raterecord. A scraper's only output is a RateRecord pushed to SQS.
- Keep provider-specific parsing isolated per dir (wise/, sendwave/). No shared mutable state.
- Tests: table-driven against saved sample payloads (fixtures), not live network. Ask before adding deps.