// Package money defines an exact monetary value type: an integer amount in a
// currency's minor units plus its ISO 4217 currency code.
package money

import (
	"errors"
	"fmt"
)

// Money is an amount of a single currency, held as an integer count of the
// currency's minor units (e.g. cents for USD, so $12.50 is Amount 1250).
// Integers keep sums of fees and principals exact; float64 rounding would
// silently corrupt money.
//
// Money does not restrict the sign of Amount. Whether a value may be zero or
// negative depends on what it represents, so that rule belongs to the caller.
type Money struct {
	Amount   int64
	Currency string
}

// ErrInvalidCurrency is returned, wrapped with the offending code, when a
// currency code is not well-formed. Match it with errors.Is.
var ErrInvalidCurrency = errors.New("invalid currency code")

// Validate reports whether m has a well-formed currency code, returning an
// error wrapping ErrInvalidCurrency if not. It does not check Amount.
func (m Money) Validate() error {
	if !IsValidCurrency(m.Currency) {
		return fmt.Errorf("money: %q: %w", m.Currency, ErrInvalidCurrency)
	}
	return nil
}

// IsValidCurrency reports whether code has the shape of an ISO 4217 currency
// code: exactly three uppercase ASCII letters. It does not check that the code
// names a real currency.
func IsValidCurrency(code string) bool {
	if len(code) != 3 {
		return false
	}
	return code[0] >= 'A' && code[0] <= 'Z' &&
		code[1] >= 'A' && code[1] <= 'Z' &&
		code[2] >= 'A' && code[2] <= 'Z'
}
