package money

import (
	"errors"
	"testing"
)

// Stored amounts are read at this scale, so changing it would misread every
// stored amount.
func TestMicrosPerUnit(t *testing.T) {
	if MicrosPerUnit != 1_000_000 {
		t.Errorf("MicrosPerUnit = %d, want 1000000", MicrosPerUnit)
	}
}

func TestMoneyValidate(t *testing.T) {
	tests := []struct {
		name    string
		money   Money
		wantErr error // nil if the value is valid
	}{
		{name: "valid CAD", money: Money{Micros: 12_500_000, Currency: "CAD"}},
		{name: "valid zero amount", money: Money{Micros: 0, Currency: "NGN"}},
		{name: "valid negative amount", money: Money{Micros: -5_000_000, Currency: "CAD"}},
		{name: "empty currency", money: Money{Micros: 1_000_000, Currency: ""}, wantErr: ErrInvalidCurrency},
		{name: "lowercase currency", money: Money{Micros: 1_000_000, Currency: "cad"}, wantErr: ErrInvalidCurrency},
		{name: "too short", money: Money{Micros: 1_000_000, Currency: "CA"}, wantErr: ErrInvalidCurrency},
		{name: "too long", money: Money{Micros: 1_000_000, Currency: "CADD"}, wantErr: ErrInvalidCurrency},
		{name: "digit in code", money: Money{Micros: 1_000_000, Currency: "CA1"}, wantErr: ErrInvalidCurrency},
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
