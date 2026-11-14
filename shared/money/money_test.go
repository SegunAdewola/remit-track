package money

import (
	"errors"
	"testing"
)

// Money is a value type: an integer amount in the currency's minor units (e.g.
// cents) plus an ISO 4217 currency code. We hold minor units as int64 rather
// than a float so that summing fees and principals is exact — float64 rounding
// silently corrupts money.
//
// Validate only checks the currency code shape. The SIGN of Amount is
// deliberately NOT Money's concern: a fee of 0 and a balance adjustment of -500
// are both well-formed Money; whether a given field may be negative is enforced
// by whoever owns that field (e.g. RateRecord.Validate rejects a non-positive
// ReceiveAmount). Keeping that policy out of Money keeps the type reusable.

func TestMoneyValidate(t *testing.T) {
	tests := []struct {
		name    string
		money   Money
		wantErr error // nil means the value should be valid
	}{
		{name: "valid USD", money: Money{Amount: 1250, Currency: "USD"}, wantErr: nil},
		{name: "valid zero amount", money: Money{Amount: 0, Currency: "GHS"}, wantErr: nil},
		{name: "valid negative amount", money: Money{Amount: -500, Currency: "EUR"}, wantErr: nil},
		{name: "empty currency", money: Money{Amount: 100, Currency: ""}, wantErr: ErrInvalidCurrency},
		{name: "lowercase currency", money: Money{Amount: 100, Currency: "usd"}, wantErr: ErrInvalidCurrency},
		{name: "too short", money: Money{Amount: 100, Currency: "US"}, wantErr: ErrInvalidCurrency},
		{name: "too long", money: Money{Amount: 100, Currency: "USDD"}, wantErr: ErrInvalidCurrency},
		{name: "digit in code", money: Money{Amount: 100, Currency: "US1"}, wantErr: ErrInvalidCurrency},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.money.Validate()
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want errors.Is(err, %v)", err, tt.wantErr)
			}
		})
	}
}
