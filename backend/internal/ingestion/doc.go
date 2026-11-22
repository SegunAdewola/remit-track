// Package ingestion consumes scraped quotes from SQS, validates and batches
// them, sends invalid messages to the dead-letter queue, and rejects corridors
// RemitTrack does not serve. It owns no tables. It declares the interfaces it
// needs (a quote writer and a batch observer); cmd/api-server supplies them.
package ingestion
