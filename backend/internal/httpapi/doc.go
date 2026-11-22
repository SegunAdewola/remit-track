// Package httpapi is the HTTP transport for the backend: routing, middleware,
// and rate limiting. Handlers hold no business logic; they call other modules
// through interfaces that cmd/api-server wires in.
package httpapi
