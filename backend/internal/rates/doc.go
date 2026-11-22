// Package rates stores provider quotes and serves them ranked by effective
// rate. It owns the rates Postgres schema (historical_rates, current_rates)
// and its read cache. Other modules reach it only through interfaces that
// cmd/api-server wires in.
package rates
