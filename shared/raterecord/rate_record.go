// Package raterecord provides RateRecord, the quote type that scrapers produce
// and the backend stores and ranks.
package raterecord

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/SegunAdewola/remit-track/shared/corridor"
	"github.com/SegunAdewola/remit-track/shared/money"
)

// RateRecord is a quote observed from a provider for one transfer: the amount
// sent, the total the sender pays, the amount the recipient receives, and the
// fees, for a given corridor, pair of payment methods, and rate class.
//
// Fields hold the values the provider reported. RateRecord does not require
// them to agree with each other. For example, TotalPaid need not equal
// SendAmount plus Fees, since some providers deduct fees from the amount
// converted instead of adding them.
//
// The JSON form is the message scrapers send to ingestion. Its field names
// are a contract between separately deployed programs and must not change.
// Decoding ignores unknown fields and does not validate; call Validate after
// decoding.
type RateRecord struct {
	Provider      string                `json:"provider"`      // provider slug, such as "sendapp"
	Pair          corridor.CurrencyPair `json:"pair"`          // Base is sent, Quote is received
	SendAmount    money.Money           `json:"sendAmount"`    // amount quoted for, in Base
	TotalPaid     money.Money           `json:"totalPaid"`     // amount the sender pays, in Base
	ReceiveAmount money.Money           `json:"receiveAmount"` // amount the recipient receives, in Quote
	Fees          money.Money           `json:"fees"`          // total fees, in Base or Quote
	OfferedRate   float64               `json:"offeredRate"`   // advertised exchange rate, for display only
	PayIn         PayMethod             `json:"payIn"`         // how the sender pays
	PayOut        PayMethod             `json:"payOut"`        // how the recipient is paid
	Class         RateClass             `json:"class"`         // pricing the quote applies to
	Timestamp     time.Time             `json:"timestamp"`     // when the quote was observed
}

// Errors wrapped by the error that [RateRecord.Validate] returns.
var (
	ErrEmptyProvider     = errors.New("empty provider")
	ErrMalformedPair     = errors.New("malformed currency pair")
	ErrNonPositiveAmount = errors.New("non-positive amount")
	ErrNegativeFees      = errors.New("negative fees")
	ErrCurrencyMismatch  = errors.New("currency mismatch")
	ErrInvalidRate       = errors.New("invalid rate")
	ErrEmptyPayMethod    = errors.New("empty pay method")
	ErrInvalidClass      = errors.New("invalid class")
	ErrZeroTimestamp     = errors.New("zero timestamp")
	ErrExcessPrecision   = errors.New("excess precision")
	ErrUnknownCurrency   = errors.New("currency not in minor-units table")
)

// Validate checks that r is well formed. If it is not, Validate returns an
// error that wraps every problem found, each matching one of the Err values in
// this package under errors.Is. A malformed Pair is reported on its own, since
// the currency checks depend on it.
//
// Validate does not check whether r.Pair is a supported corridor.
func (r RateRecord) Validate() error {
	if err := checkPair(r.Pair); err != nil {
		return fmt.Errorf("raterecord: quote from %q: %w", r.Provider, err)
	}
	base, quote := r.Pair.Base, r.Pair.Quote
	err := errors.Join(
		checkNonEmpty("provider", r.Provider, ErrEmptyProvider),
		checkPositive("send amount", r.SendAmount),
		checkPositive("total paid", r.TotalPaid),
		checkPositive("receive amount", r.ReceiveAmount),
		checkNonNegative("fees", r.Fees),
		checkCurrency("send amount", r.SendAmount, base),
		checkCurrency("total paid", r.TotalPaid, base),
		checkCurrency("receive amount", r.ReceiveAmount, quote),
		checkCurrency("fees", r.Fees, base, quote),
		checkPrecision("send amount", r.SendAmount),
		checkPrecision("total paid", r.TotalPaid),
		checkPrecision("receive amount", r.ReceiveAmount),
		checkPrecision("fees", r.Fees),
		checkRate(r.OfferedRate),
		checkNonEmpty("pay-in method", string(r.PayIn), ErrEmptyPayMethod),
		checkNonEmpty("pay-out method", string(r.PayOut), ErrEmptyPayMethod),
		checkClass(r.Class),
		checkTimestamp(r.Timestamp),
	)
	if err != nil {
		return fmt.Errorf("raterecord: quote from %q: %w", r.Provider, err)
	}
	return nil
}

// checkPair reports a pair whose codes are malformed or identical.
func checkPair(p corridor.CurrencyPair) error {
	if !money.IsValidCurrency(p.Base) || !money.IsValidCurrency(p.Quote) || p.Base == p.Quote {
		return fmt.Errorf("pair %q: %w", p, ErrMalformedPair)
	}
	return nil
}

// checkNonEmpty returns sentinel, wrapped with name, if s is empty.
func checkNonEmpty(name, s string, sentinel error) error {
	if s == "" {
		return fmt.Errorf("%s: %w", name, sentinel)
	}
	return nil
}

func checkPositive(name string, m money.Money) error {
	if m.Micros <= 0 {
		return fmt.Errorf("%s %d %s micros: %w", name, m.Micros, m.Currency, ErrNonPositiveAmount)
	}
	return nil
}

func checkNonNegative(name string, m money.Money) error {
	if m.Micros < 0 {
		return fmt.Errorf("%s %d %s micros: %w", name, m.Micros, m.Currency, ErrNegativeFees)
	}
	return nil
}

// checkPrecision reports m if its currency is not in money.MinorUnits, or if m
// has more decimal places than the currency allows.
func checkPrecision(name string, m money.Money) error {
	digits, ok := money.MinorUnits(m.Currency)
	if !ok {
		return fmt.Errorf("%s in %q: %w", name, m.Currency, ErrUnknownCurrency)
	}
	// step is the smallest amount the currency can express, in micros:
	// 10_000 for 2 digits (one cent), 1_000_000 for 0 digits (one franc).
	step := int64(money.MicrosPerUnit)
	for range digits {
		step /= 10
	}
	if m.Micros%step != 0 {
		return fmt.Errorf("%s %d %s micros, want a multiple of %d: %w",
			name, m.Micros, m.Currency, step, ErrExcessPrecision)
	}
	return nil
}

// checkCurrency reports m if its currency is not one of allowed.
func checkCurrency(name string, m money.Money, allowed ...string) error {
	if !slices.Contains(allowed, m.Currency) {
		return fmt.Errorf("%s in %q, want one of %q: %w", name, m.Currency, allowed, ErrCurrencyMismatch)
	}
	return nil
}

// checkRate reports a rate that is not a finite positive number. The NaN test
// is needed because every comparison with NaN, including rate <= 0, is false.
func checkRate(rate float64) error {
	if math.IsNaN(rate) || math.IsInf(rate, 0) || rate <= 0 {
		return fmt.Errorf("offered rate %v: %w", rate, ErrInvalidRate)
	}
	return nil
}

func checkClass(c RateClass) error {
	if !c.Valid() {
		return fmt.Errorf("class %q: %w", c, ErrInvalidClass)
	}
	return nil
}

func checkTimestamp(t time.Time) error {
	if t.IsZero() {
		return fmt.Errorf("timestamp: %w", ErrZeroTimestamp)
	}
	return nil
}

// PayMethod identifies how money is paid in or out, such as "bank_transfer" or
// "mobile_money". Any non-empty value is valid.
type PayMethod string

// RateClass identifies the pricing a quote applies to. The values are stored
// and sent between services, so existing values must not be changed.
type RateClass string

// Rate classes. Only ClassStandard quotes are ranked.
const (
	ClassStandard           RateClass = "standard"             // available to every customer
	ClassPromoFirstTransfer RateClass = "promo_first_transfer" // first transfer only
)

// Valid reports whether c is a defined rate class.
func (c RateClass) Valid() bool {
	switch c {
	case ClassStandard, ClassPromoFirstTransfer:
		return true
	default:
		return false
	}
}

// EffectiveRate returns the amount the recipient receives per unit the sender
// pays, fees included: ReceiveAmount divided by TotalPaid, in Quote per Base.
// Quotes are ranked by this value, higher being better, because it accounts
// for fees that the advertised OfferedRate leaves out.
//
// EffectiveRate validates r first and returns the error from
// [RateRecord.Validate] if r is not well formed.
func (r RateRecord) EffectiveRate() (float64, error) {
	if err := r.Validate(); err != nil {
		return 0, fmt.Errorf("raterecord: effective rate: %w", err)
	}

	return float64(r.ReceiveAmount.Micros) / (float64(r.TotalPaid.Micros)), nil
}
