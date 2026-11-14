package raterecord

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/SegunAdewola/remit-track/shared/corridor"
	"github.com/SegunAdewola/remit-track/shared/money"
)

// validRecord returns a known-good USD -> GHS quote that each test case mutates
// in exactly one way, so a failure points at a single invariant.
//
// The figures model an "added on top" fee: the sender pays $100 + $2.99 fee,
// and the recipient gets GHS 1,525.00 at an offered rate of 15.25.
func validRecord() RateRecord {
	return RateRecord{
		Provider:      "wise",
		Pair:          corridor.CurrencyPair{Base: "USD", Quote: "GHS"},
		SendAmount:    money.Money{Amount: 10000, Currency: "USD"},
		TotalPaid:     money.Money{Amount: 10299, Currency: "USD"},
		ReceiveAmount: money.Money{Amount: 152500, Currency: "GHS"},
		Fees:          money.Money{Amount: 299, Currency: "USD"},
		OfferedRate:   15.25,
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
		wantErr error // nil means the record should be valid
	}{
		// Valid variants.
		{name: "valid record", mutate: func(*RateRecord) {}},
		{name: "zero fees", mutate: func(r *RateRecord) { r.Fees.Amount = 0 }},
		{name: "fees in receive currency", mutate: func(r *RateRecord) { r.Fees = money.Money{Amount: 4500, Currency: "GHS"} }},
		{name: "promo class", mutate: func(r *RateRecord) { r.Class = ClassPromoFirstTransfer }},
		{
			// Deducted-fee model: TotalPaid equals SendAmount. Validate must not
			// assume TotalPaid == SendAmount + Fees.
			name: "fee deducted, not added",
			mutate: func(r *RateRecord) {
				r.TotalPaid = r.SendAmount
				r.ReceiveAmount.Amount = 148000
			},
		},

		// Provider.
		{name: "empty provider", mutate: func(r *RateRecord) { r.Provider = "" }, wantErr: ErrEmptyProvider},

		// Pair: well-formed only; supported-corridor checks live at ingestion.
		{name: "empty pair", mutate: func(r *RateRecord) { r.Pair = corridor.CurrencyPair{} }, wantErr: ErrMalformedPair},
		{name: "malformed base", mutate: func(r *RateRecord) { r.Pair.Base = "US" }, wantErr: ErrMalformedPair},
		{name: "malformed quote", mutate: func(r *RateRecord) { r.Pair.Quote = "ghs" }, wantErr: ErrMalformedPair},
		{name: "base equals quote", mutate: func(r *RateRecord) { r.Pair.Quote = "USD" }, wantErr: ErrMalformedPair},

		// Amounts must be strictly positive.
		{name: "zero send amount", mutate: func(r *RateRecord) { r.SendAmount.Amount = 0 }, wantErr: ErrNonPositiveAmount},
		{name: "negative total paid", mutate: func(r *RateRecord) { r.TotalPaid.Amount = -1 }, wantErr: ErrNonPositiveAmount},
		{name: "zero receive amount", mutate: func(r *RateRecord) { r.ReceiveAmount.Amount = 0 }, wantErr: ErrNonPositiveAmount},

		// Fees may be zero but never negative.
		{name: "negative fees", mutate: func(r *RateRecord) { r.Fees.Amount = -1 }, wantErr: ErrNegativeFees},

		// Each Money field must be in the currency of its side of the pair.
		{name: "send amount in quote currency", mutate: func(r *RateRecord) { r.SendAmount.Currency = "GHS" }, wantErr: ErrCurrencyMismatch},
		{name: "total paid in quote currency", mutate: func(r *RateRecord) { r.TotalPaid.Currency = "GHS" }, wantErr: ErrCurrencyMismatch},
		{name: "receive amount in base currency", mutate: func(r *RateRecord) { r.ReceiveAmount.Currency = "USD" }, wantErr: ErrCurrencyMismatch},
		{name: "fees in unrelated currency", mutate: func(r *RateRecord) { r.Fees.Currency = "EUR" }, wantErr: ErrCurrencyMismatch},
		{name: "malformed money currency", mutate: func(r *RateRecord) { r.SendAmount.Currency = "usd" }, wantErr: ErrCurrencyMismatch},

		// Offered rate must be a finite positive number.
		{name: "zero rate", mutate: func(r *RateRecord) { r.OfferedRate = 0 }, wantErr: ErrInvalidRate},
		{name: "negative rate", mutate: func(r *RateRecord) { r.OfferedRate = -1 }, wantErr: ErrInvalidRate},
		{name: "NaN rate", mutate: func(r *RateRecord) { r.OfferedRate = math.NaN() }, wantErr: ErrInvalidRate},
		{name: "infinite rate", mutate: func(r *RateRecord) { r.OfferedRate = math.Inf(1) }, wantErr: ErrInvalidRate},

		// Payment methods must be present.
		{name: "empty pay-in", mutate: func(r *RateRecord) { r.PayIn = "" }, wantErr: ErrEmptyPayMethod},
		{name: "empty pay-out", mutate: func(r *RateRecord) { r.PayOut = "" }, wantErr: ErrEmptyPayMethod},

		// Class must be in the closed set.
		{name: "empty class", mutate: func(r *RateRecord) { r.Class = "" }, wantErr: ErrInvalidClass},
		{name: "unknown class", mutate: func(r *RateRecord) { r.Class = RateClass("vip") }, wantErr: ErrInvalidClass},

		// Timestamp.
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

// The wire values of RateClass are persisted (rate_class column) and sent over
// SQS, so they are part of the contract and must not change silently.
func TestRateClassWireValues(t *testing.T) {
	if ClassStandard != "standard" {
		t.Errorf("ClassStandard = %q, want %q", ClassStandard, "standard")
	}
	if ClassPromoFirstTransfer != "promo_first_transfer" {
		t.Errorf("ClassPromoFirstTransfer = %q, want %q", ClassPromoFirstTransfer, "promo_first_transfer")
	}
}
