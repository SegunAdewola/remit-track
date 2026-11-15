// Package corridor provides types that identify remittance corridors.
package corridor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/SegunAdewola/remit-track/shared/money"
)

// CurrencyPair identifies a corridor by the currency sent (Base) and the
// currency received (Quote). Its string form is "BASE:QUOTE", such as
// "CAD:NGN".
type CurrencyPair struct {
	Base  string
	Quote string
}

// ErrInvalidPair indicates a string that is not a valid "BASE:QUOTE" pair.
var ErrInvalidPair = errors.New("invalid currency pair")

// Parse parses s in the form "BASE:QUOTE", such as "CAD:NGN". Each code must
// satisfy [money.IsValidCurrency]. If s is malformed, Parse returns an error
// wrapping ErrInvalidPair. Parse does not check that the codes name real
// currencies.
func Parse(s string) (CurrencyPair, error) {
	base, quote, found := strings.Cut(s, ":")
	if !found || !money.IsValidCurrency(base) || !money.IsValidCurrency(quote) {
		return CurrencyPair{}, fmt.Errorf("corridor: %q: %w", s, ErrInvalidPair)
	}
	return CurrencyPair{Base: base, Quote: quote}, nil
}

// String returns the pair in the form "BASE:QUOTE".
func (c CurrencyPair) String() string {
	return c.Base + ":" + c.Quote
}
