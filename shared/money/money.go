// Package money provides a value type for an amount of money in one currency.
package money

import (
	"errors"
	"fmt"
)

// Money is an amount in micros, millionths of one whole unit, of a currency.
// The scale is the same for every currency: CAD 12.50 is Money{Micros:
// 12_500_000, Currency: "CAD"} and XOF 430 is Money{Micros: 430_000_000,
// Currency: "XOF"}.
//
// Micros is an integer so that arithmetic on money is exact. Money places no
// constraint on its sign; callers decide whether zero or negative amounts are
// allowed.
type Money struct {
	Micros   int64
	Currency string // ISO 4217 code, such as "CAD"
}

// MicrosPerUnit is the number of micros in one whole unit of any currency.
const MicrosPerUnit = 1_000_000

// ErrInvalidCurrency indicates a malformed currency code.
var ErrInvalidCurrency = errors.New("invalid currency code")

// Validate returns an error wrapping ErrInvalidCurrency if m.Currency is not a
// well-formed currency code. It does not check Micros.
func (m Money) Validate() error {
	if !IsValidCurrency(m.Currency) {
		return fmt.Errorf("money: %q: %w", m.Currency, ErrInvalidCurrency)
	}
	return nil
}

// IsValidCurrency reports whether code is three uppercase ASCII letters, the
// form of an ISO 4217 code. It does not check that the currency exists.
func IsValidCurrency(code string) bool {
	if len(code) != 3 {
		return false
	}
	return code[0] >= 'A' && code[0] <= 'Z' &&
		code[1] >= 'A' && code[1] <= 'Z' &&
		code[2] >= 'A' && code[2] <= 'Z'
}
