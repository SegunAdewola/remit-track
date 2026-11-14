// Package raterecord defines RateRecord, the canonical quote a provider returns
// for one transfer, shared by the scrapers that produce it and the backend that
// stores and ranks it.
package raterecord

import (
	"errors"
	"math"
	"time"

	"github.com/SegunAdewola/remit-track/shared/corridor"
	"github.com/SegunAdewola/remit-track/shared/money"
)

// RateRecord is one quote observed from a provider: what a sender pays and
// what the recipient receives for a given corridor, amount, payment methods
// and rate class at an instant.
//
// Every value is the provider's own output, captured as observed. Nothing here
// is derived from a fee formula, because a miscomputed price looks exactly like
// a correct one. For that reason RateRecord holds no cross-field invariants
// such as TotalPaid == SendAmount + Fees: providers differ in whether fees are
// added on top or deducted from the converted amount.
type RateRecord struct {
	// Provider is the provider's slug, e.g. "wise".
	Provider string
	// Pair is the corridor: Base is the send currency, Quote the receive currency.
	Pair corridor.CurrencyPair
	// SendAmount is the amount the quote was requested for, in Base.
	SendAmount money.Money
	// TotalPaid is what the sender actually pays, fees included, in Base.
	TotalPaid money.Money
	// ReceiveAmount is what the recipient gets, in Quote.
	ReceiveAmount money.Money
	// Fees is the total fee charged, in either currency of the pair. May be zero.
	Fees money.Money
	// OfferedRate is the provider's advertised exchange rate. It is for display
	// only; rankings use amounts, never this float.
	OfferedRate float64
	// PayIn is how the sender funds the transfer, e.g. "bank_transfer".
	PayIn PayMethod
	// PayOut is how the recipient is paid, e.g. "mobile_money".
	PayOut PayMethod
	// Class distinguishes the standard price from promotional ones.
	Class RateClass
	// Timestamp is when the quote was observed.
	Timestamp time.Time
}

// Errors returned by [RateRecord.Validate], wrapped with the offending field
// and value. Match them with errors.Is.
var (
	// ErrEmptyProvider means Provider is empty.
	ErrEmptyProvider = errors.New("empty provider")
	// ErrMalformedPair means Pair has a malformed currency code or Base == Quote.
	ErrMalformedPair = errors.New("malformed currency pair")
	// ErrNonPositiveAmount means SendAmount, TotalPaid or ReceiveAmount is not > 0.
	ErrNonPositiveAmount = errors.New("non-positive amount")
	// ErrNegativeFees means Fees is below zero.
	ErrNegativeFees = errors.New("negative fees")
	// ErrCurrencyMismatch means a money field is not in the currency of its side
	// of the pair.
	ErrCurrencyMismatch = errors.New("currency mismatch")
	// ErrInvalidRate means OfferedRate is zero, negative, NaN or infinite.
	ErrInvalidRate = errors.New("invalid rate")
	// ErrEmptyPayMethod means PayIn or PayOut is empty.
	ErrEmptyPayMethod = errors.New("empty pay method")
	// ErrInvalidClass means Class is not one of the defined RateClass values.
	ErrInvalidClass = errors.New("invalid class")
	// ErrZeroTimestamp means Timestamp is the zero time.
	ErrZeroTimestamp = errors.New("zero timestamp")
)

// Validate reports whether r is a well-formed quote. It returns nil if so, or
// an error wrapping one or more of the Err* values above.
//
// Validate checks shape only. It does not check that r.Pair is a corridor
// RemitTrack supports; that allowlist is configuration and is enforced at
// ingestion.
func (r RateRecord) Validate() error {
	if r.Provider == "" {
		return ErrEmptyProvider
	}
	if r.Pair == (corridor.CurrencyPair{}) {
		return ErrMalformedPair
	}
	if !money.IsValidCurrency(r.Pair.Base) {
		return ErrMalformedPair
	}
	if !money.IsValidCurrency(r.Pair.Quote) {
		return ErrMalformedPair
	}
	if r.Pair.Base == r.Pair.Quote {
		return ErrMalformedPair
	}
	if r.SendAmount.Amount <= 0 {
		return ErrNonPositiveAmount
	}
	if r.TotalPaid.Amount <= 0 {
		return ErrNonPositiveAmount
	}
	if r.ReceiveAmount.Amount <= 0 {
		return ErrNonPositiveAmount
	}
	if r.Fees.Amount < 0 {
		return ErrNegativeFees
	}
	if r.SendAmount.Currency != r.Pair.Base {
		return ErrCurrencyMismatch
	}
	if r.TotalPaid.Currency != r.Pair.Base {
		return ErrCurrencyMismatch
	}
	if r.ReceiveAmount.Currency != r.Pair.Quote {
		return ErrCurrencyMismatch
	}
	if r.Fees.Currency != r.Pair.Base && r.Fees.Currency != r.Pair.Quote {
		return ErrCurrencyMismatch
	}
	if r.OfferedRate <= 0 {
		return ErrInvalidRate
	}
	if math.IsNaN(r.OfferedRate) {
		return ErrInvalidRate
	}
	if math.IsInf(r.OfferedRate, 0) {
		return ErrInvalidRate
	}
	if r.PayIn == "" {
		return ErrEmptyPayMethod
	}
	if r.PayOut == "" {
		return ErrEmptyPayMethod
	}
	if !r.Class.Valid() {
		return ErrInvalidClass
	}
	if r.Timestamp.IsZero() {
		return ErrZeroTimestamp
	}
	return nil
}

// PayMethod names how money enters or leaves a transfer, e.g. "bank_transfer"
// or "mobile_money". The set is open: any non-empty value is accepted until
// provider parsers settle a shared vocabulary.
type PayMethod string

// RateClass distinguishes a provider's standard price from promotional ones.
// Its values are stored and sent over the wire, so they must never change.
type RateClass string

// The defined rate classes. Only ClassStandard is ranked; promotional classes
// are shown as a badge so returning users are never shown a price they cannot
// get.
const (
	// ClassStandard is the price any customer gets.
	ClassStandard RateClass = "standard"
	// ClassPromoFirstTransfer is a discounted price for a customer's first
	// transfer only.
	ClassPromoFirstTransfer RateClass = "promo_first_transfer"
)

// Valid reports whether c is one of the defined rate classes.
func (c RateClass) Valid() bool {
	switch c {
	case ClassStandard, ClassPromoFirstTransfer:
		return true
	default:
		return false
	}
}
