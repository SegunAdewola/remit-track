# frontend/ (Next.js App Router, TypeScript)

Server and build-time rendering, not a SPA. The SEO strategy depends on it (docs/sadd.md Section A).

- Routes under app/[locale]/exchange-rates/[corridor]/. ISR pages revalidate via the ingestion webhook.
- JSON-LD builders (ExchangeRateSpecification, FAQPage) live in lib/schema/. Keep them pure and unit-tested.
- Data fetching in lib/api/. Components stay presentational; no fetch calls inside components.
- Tests: unit-test lib/; use Playwright only for true end-to-end. Ask before adding libraries.
- Run npm test, npm run lint, npm run build before every commit.