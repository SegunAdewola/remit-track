// Package money provides a value type for an amount of money in one currency.
package money

import (
	"errors"
	"fmt"
)

// Money is an amount in the minor units of a currency, such as cents for CAD
// or kobo for NGN. For example, CAD 12.50 is Money{Amount: 1250, Currency:
// "CAD"}. [MinorUnits] gives the number of decimal places for a currency.
//
// Amount is an integer so that arithmetic on money is exact. Money places no
// constraint on the sign of Amount; callers decide whether zero or negative
// amounts are allowed.
type Money struct {
	Amount   int64
	Currency string // ISO 4217 code, such as "CAD"
}

// ErrInvalidCurrency indicates a malformed currency code.
var ErrInvalidCurrency = errors.New("invalid currency code")

// Validate returns an error wrapping ErrInvalidCurrency if m.Currency is not a
// well-formed currency code. It does not check Amount.
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
