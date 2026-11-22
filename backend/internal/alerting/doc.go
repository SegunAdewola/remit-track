// Package alerting evaluates user rate alerts against each new batch of
// quotes and delivers notifications through a transactional outbox. It owns
// the alerting Postgres schema (alert configs, trigger state, notifications).
package alerting
