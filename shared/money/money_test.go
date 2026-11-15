package money

import (
	"errors"
	"testing"
)

func TestMoneyValidate(t *testing.T) {
	tests := []struct {
		name    string
		money   Money
		wantErr error // nil if the value is valid
	}{
		{name: "valid CAD", money: Money{Amount: 1250, Currency: "CAD"}},
		{name: "valid zero amount", money: Money{Amount: 0, Currency: "NGN"}},
		{name: "valid negative amount", money: Money{Amount: -500, Currency: "CAD"}},
		{name: "empty currency", money: Money{Amount: 100, Currency: ""}, wantErr: ErrInvalidCurrency},
		{name: "lowercase currency", money: Money{Amount: 100, Currency: "cad"}, wantErr: ErrInvalidCurrency},
		{name: "too short", money: Money{Amount: 100, Currency: "CA"}, wantErr: ErrInvalidCurrency},
		{name: "too long", money: Money{Amount: 100, Currency: "CADD"}, wantErr: ErrInvalidCurrency},
		{name: "digit in code", money: Money{Amount: 100, Currency: "CA1"}, wantErr: ErrInvalidCurrency},
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
