package raterecord

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/SegunAdewola/remit-track/shared/corridor"
	"github.com/SegunAdewola/remit-track/shared/money"
)

// validRecord returns a valid CAD to NGN quote. Each test case changes one
// field so that a failure points at one rule.
//
// The fee is added on top: the sender pays CAD 100.00 plus a CAD 2.99 fee, and
// the recipient receives NGN 110,050.00 at a rate of 1100.50.
func validRecord() RateRecord {
	return RateRecord{
		Provider:      "sendapp",
		Pair:          corridor.CurrencyPair{Base: "CAD", Quote: "NGN"},
		SendAmount:    money.Money{Amount: 10000, Currency: "CAD"},
		TotalPaid:     money.Money{Amount: 10299, Currency: "CAD"},
		ReceiveAmount: money.Money{Amount: 11005000, Currency: "NGN"},
		Fees:          money.Money{Amount: 299, Currency: "CAD"},
		OfferedRate:   1100.50,
		PayIn:         PayMethod("bank_transfer"),
		PayOut:        PayMethod("mobile_money"),
		Class:         ClassStandard,
		Timestamp:     time.Unix(1_700_000_000, 0),
	}
}

func TestRateRecordValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*RateRecord)
		wantErr error // nil if the record is valid
	}{
		{name: "valid record", mutate: func(*RateRecord) {}},
		{name: "zero fees", mutate: func(r *RateRecord) { r.Fees.Amount = 0 }},
		{name: "fees in receive currency", mutate: func(r *RateRecord) { r.Fees = money.Money{Amount: 330000, Currency: "NGN"} }},
		{name: "promo class", mutate: func(r *RateRecord) { r.Class = ClassPromoFirstTransfer }},
		{
			// Some providers deduct the fee from the amount converted, so
			// TotalPaid equals SendAmount.
			name: "fee deducted, not added",
			mutate: func(r *RateRecord) {
				r.TotalPaid = r.SendAmount
				r.ReceiveAmount.Amount = 10675950
			},
		},

		{name: "empty provider", mutate: func(r *RateRecord) { r.Provider = "" }, wantErr: ErrEmptyProvider},

		// Validate checks that the pair is well formed, not that it is supported.
		{name: "empty pair", mutate: func(r *RateRecord) { r.Pair = corridor.CurrencyPair{} }, wantErr: ErrMalformedPair},
		{name: "malformed base", mutate: func(r *RateRecord) { r.Pair.Base = "CA" }, wantErr: ErrMalformedPair},
		{name: "malformed quote", mutate: func(r *RateRecord) { r.Pair.Quote = "ngn" }, wantErr: ErrMalformedPair},
		{name: "base equals quote", mutate: func(r *RateRecord) { r.Pair.Quote = "CAD" }, wantErr: ErrMalformedPair},

		{name: "zero send amount", mutate: func(r *RateRecord) { r.SendAmount.Amount = 0 }, wantErr: ErrNonPositiveAmount},
		{name: "negative total paid", mutate: func(r *RateRecord) { r.TotalPaid.Amount = -1 }, wantErr: ErrNonPositiveAmount},
		{name: "zero receive amount", mutate: func(r *RateRecord) { r.ReceiveAmount.Amount = 0 }, wantErr: ErrNonPositiveAmount},
		{name: "negative fees", mutate: func(r *RateRecord) { r.Fees.Amount = -1 }, wantErr: ErrNegativeFees},

		{name: "send amount in quote currency", mutate: func(r *RateRecord) { r.SendAmount.Currency = "NGN" }, wantErr: ErrCurrencyMismatch},
		{name: "total paid in quote currency", mutate: func(r *RateRecord) { r.TotalPaid.Currency = "NGN" }, wantErr: ErrCurrencyMismatch},
		{name: "receive amount in base currency", mutate: func(r *RateRecord) { r.ReceiveAmount.Currency = "CAD" }, wantErr: ErrCurrencyMismatch},
		{name: "fees in unrelated currency", mutate: func(r *RateRecord) { r.Fees.Currency = "EUR" }, wantErr: ErrCurrencyMismatch},
		{name: "malformed money currency", mutate: func(r *RateRecord) { r.SendAmount.Currency = "cad" }, wantErr: ErrCurrencyMismatch},

		{name: "zero rate", mutate: func(r *RateRecord) { r.OfferedRate = 0 }, wantErr: ErrInvalidRate},
		{name: "negative rate", mutate: func(r *RateRecord) { r.OfferedRate = -1 }, wantErr: ErrInvalidRate},
		{name: "NaN rate", mutate: func(r *RateRecord) { r.OfferedRate = math.NaN() }, wantErr: ErrInvalidRate},
		{name: "infinite rate", mutate: func(r *RateRecord) { r.OfferedRate = math.Inf(1) }, wantErr: ErrInvalidRate},

		{name: "empty pay-in", mutate: func(r *RateRecord) { r.PayIn = "" }, wantErr: ErrEmptyPayMethod},
		{name: "empty pay-out", mutate: func(r *RateRecord) { r.PayOut = "" }, wantErr: ErrEmptyPayMethod},

		{name: "empty class", mutate: func(r *RateRecord) { r.Class = "" }, wantErr: ErrInvalidClass},
		{name: "unknown class", mutate: func(r *RateRecord) { r.Class = RateClass("vip") }, wantErr: ErrInvalidClass},

		{name: "zero timestamp", mutate: func(r *RateRecord) { r.Timestamp = time.Time{} }, wantErr: ErrZeroTimestamp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validRecord()
			tt.mutate(&r)
			err := r.Validate()
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

func TestRateClassValid(t *testing.T) {
	tests := []struct {
		class RateClass
		want  bool
	}{
		{class: ClassStandard, want: true},
		{class: ClassPromoFirstTransfer, want: true},
		{class: "", want: false},
		{class: "STANDARD", want: false},
		{class: "promo", want: false},
	}
	for _, tt := range tests {
		t.Run(string(tt.class), func(t *testing.T) {
			if got := tt.class.Valid(); got != tt.want {
				t.Errorf("RateClass(%q).Valid() = %v, want %v", tt.class, got, tt.want)
			}
		})
	}
}

// Rate class values are stored and sent between services, so changing one is
// a breaking change.
func TestRateClassWireValues(t *testing.T) {
	if ClassStandard != "standard" {
		t.Errorf("ClassStandard = %q, want %q", ClassStandard, "standard")
	}
	if ClassPromoFirstTransfer != "promo_first_transfer" {
		t.Errorf("ClassPromoFirstTransfer = %q, want %q", ClassPromoFirstTransfer, "promo_first_transfer")
	}
}
