// Package postgres provides the shared Postgres connection pool, its
// configuration from the environment, and the migration runner. It holds no
// domain logic or tables; each module owns its own schema.
package postgres
