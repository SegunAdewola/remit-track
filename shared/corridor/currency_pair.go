// Package corridor defines value types for remittance corridors, starting with
// the currency pair that a rate applies to.
package corridor

import (
	"errors"
	"fmt"
	"strings"
)

// CurrencyPair identifies a remittance corridor by its base and quote ISO 4217
// currency codes, e.g. USD -> GHS.
type CurrencyPair struct {
	Base  string
	Quote string
}

// ErrInvalidPair is returned by Parse when a string is not a well-formed
// "BASE:QUOTE" currency pair.
var ErrInvalidPair = errors.New("invalid currency pair")

// Parse reads a canonical "BASE:QUOTE" currency pair, e.g. "USD:GHS", where
// each side is exactly three uppercase ASCII letters. It returns ErrInvalidPair
// if s is malformed. Whether the codes name real currencies is not checked here.
func Parse(s string) (CurrencyPair, error) {
	base, quote, found := strings.Cut(s, ":")
	if !found || !isValidCurrency(base) || !isValidCurrency(quote) {
		return CurrencyPair{}, fmt.Errorf("corridor: %q: %w", s, ErrInvalidPair)
	}
	return CurrencyPair{Base: base, Quote: quote}, nil
}

// isValidCurrency reports whether s is exactly three uppercase ASCII letters.
func isValidCurrency(s string) bool {
	if len(s) != 3 {
		return false
	}
	return s[0] >= 'A' && s[0] <= 'Z' &&
		s[1] >= 'A' && s[1] <= 'Z' &&
		s[2] >= 'A' && s[2] <= 'Z'
}

// String returns the canonical "BASE:QUOTE" form of the pair.
func (c CurrencyPair) String() string {
	return c.Base + ":" + c.Quote
}
