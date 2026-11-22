package raterecord

import (
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"strings"
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
		SendAmount:    money.Money{Micros: 100_000_000, Currency: "CAD"},
		TotalPaid:     money.Money{Micros: 102_990_000, Currency: "CAD"},
		ReceiveAmount: money.Money{Micros: 110_050_000_000, Currency: "NGN"},
		Fees:          money.Money{Micros: 2_990_000, Currency: "CAD"},
		OfferedRate:   1100.50,
		PayIn:         PayMethod("bank_transfer"),
		PayOut:        PayMethod("mobile_money"),
		Class:         ClassStandard,
		Timestamp:     time.Unix(1_700_000_000, 0).UTC(),
	}
}

func TestRateRecordValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*RateRecord)
		wantErr error // nil if the record is valid
	}{
		{name: "valid record", mutate: func(*RateRecord) {}},
		{name: "zero fees", mutate: func(r *RateRecord) { r.Fees.Micros = 0 }},
		{name: "fees in receive currency", mutate: func(r *RateRecord) { r.Fees = money.Money{Micros: 3_300_000_000, Currency: "NGN"} }},
		{name: "promo class", mutate: func(r *RateRecord) { r.Class = ClassPromoFirstTransfer }},
		{
			// Some providers deduct the fee from the amount converted, so
			// TotalPaid equals SendAmount.
			name: "fee deducted, not added",
			mutate: func(r *RateRecord) {
				r.TotalPaid = r.SendAmount
				r.ReceiveAmount.Micros = 106_759_500_000
			},
		},

		{name: "empty provider", mutate: func(r *RateRecord) { r.Provider = "" }, wantErr: ErrEmptyProvider},

		// Validate checks that the pair is well formed, not that it is supported.
		{name: "empty pair", mutate: func(r *RateRecord) { r.Pair = corridor.CurrencyPair{} }, wantErr: ErrMalformedPair},
		{name: "malformed base", mutate: func(r *RateRecord) { r.Pair.Base = "CA" }, wantErr: ErrMalformedPair},
		{name: "malformed quote", mutate: func(r *RateRecord) { r.Pair.Quote = "ngn" }, wantErr: ErrMalformedPair},
		{name: "base equals quote", mutate: func(r *RateRecord) { r.Pair.Quote = "CAD" }, wantErr: ErrMalformedPair},

		{name: "zero send amount", mutate: func(r *RateRecord) { r.SendAmount.Micros = 0 }, wantErr: ErrNonPositiveAmount},
		{name: "negative total paid", mutate: func(r *RateRecord) { r.TotalPaid.Micros = -1 }, wantErr: ErrNonPositiveAmount},
		{name: "zero receive amount", mutate: func(r *RateRecord) { r.ReceiveAmount.Micros = 0 }, wantErr: ErrNonPositiveAmount},
		{name: "negative fees", mutate: func(r *RateRecord) { r.Fees.Micros = -1 }, wantErr: ErrNegativeFees},

		{name: "send amount in quote currency", mutate: func(r *RateRecord) { r.SendAmount.Currency = "NGN" }, wantErr: ErrCurrencyMismatch},
		{name: "total paid in quote currency", mutate: func(r *RateRecord) { r.TotalPaid.Currency = "NGN" }, wantErr: ErrCurrencyMismatch},
		{name: "receive amount in base currency", mutate: func(r *RateRecord) { r.ReceiveAmount.Currency = "CAD" }, wantErr: ErrCurrencyMismatch},
		{name: "fees in unrelated currency", mutate: func(r *RateRecord) { r.Fees.Currency = "EUR" }, wantErr: ErrCurrencyMismatch},
		{name: "malformed money currency", mutate: func(r *RateRecord) { r.SendAmount.Currency = "cad" }, wantErr: ErrCurrencyMismatch},

		// Amounts may not be more precise than their currency's minor units.
		{
			// XOF has no minor unit, so whole francs are valid.
			name: "whole-unit XOF amount",
			mutate: func(r *RateRecord) {
				r.Pair.Quote = "XOF"
				r.ReceiveAmount = money.Money{Micros: 43_000_000_000, Currency: "XOF"}
			},
		},
		{name: "sub-cent send amount", mutate: func(r *RateRecord) { r.SendAmount.Micros = 100_005_000 }, wantErr: ErrExcessPrecision},
		{name: "sub-cent fees", mutate: func(r *RateRecord) { r.Fees.Micros = 2_995_000 }, wantErr: ErrExcessPrecision},
		{name: "float noise in receive amount", mutate: func(r *RateRecord) { r.ReceiveAmount.Micros = 110_050_000_001 }, wantErr: ErrExcessPrecision},
		{
			name: "fractional XOF amount",
			mutate: func(r *RateRecord) {
				r.Pair.Quote = "XOF"
				r.ReceiveAmount = money.Money{Micros: 430_500_000, Currency: "XOF"}
			},
			wantErr: ErrExcessPrecision,
		},
		{
			// JPY is a well-formed code that is not in money.MinorUnits.
			name: "currency outside the minor-units table",
			mutate: func(r *RateRecord) {
				r.Pair.Quote = "JPY"
				r.ReceiveAmount = money.Money{Micros: 15_000_000_000, Currency: "JPY"}
			},
			wantErr: ErrUnknownCurrency,
		},

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

func TestEffectiveRate(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*RateRecord)
		want   float64
	}{
		{
			// NGN 110,050.00 received for CAD 102.99 paid (fee added on top).
			name:   "CAD to NGN",
			mutate: func(*RateRecord) {},
			want:   1068.5503446936596,
		},
		{
			// XOF has no minor unit and CAD has two. Dividing minor units would
			// give 4.3; micros use one scale for both, so the rate is 430.
			name: "CAD to XOF",
			mutate: func(r *RateRecord) {
				r.Pair.Quote = "XOF"
				r.TotalPaid = r.SendAmount
				r.Fees.Micros = 0
				r.ReceiveAmount = money.Money{Micros: 43_000_000_000, Currency: "XOF"}
			},
			want: 430,
		},
		{
			// The CAD 2.99 fee is deducted from the amount converted, so the
			// sender pays CAD 100.00 in total. The rate must use TotalPaid, not
			// SendAmount plus Fees, which would give 1036.60.
			name: "fee deducted",
			mutate: func(r *RateRecord) {
				r.TotalPaid = r.SendAmount
				r.ReceiveAmount.Micros = 106_759_500_000
			},
			want: 1067.595,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validRecord()
			tt.mutate(&r)
			got, err := r.EffectiveRate()
			if err != nil {
				t.Fatalf("EffectiveRate() error = %v, want nil", err)
			}
			// Division rounds, so compare within a relative tolerance.
			if math.Abs(got-tt.want) > 1e-9*tt.want {
				t.Errorf("EffectiveRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEffectiveRateInvalidRecord(t *testing.T) {
	r := validRecord()
	r.TotalPaid.Micros = 0
	if _, err := r.EffectiveRate(); !errors.Is(err, ErrNonPositiveAmount) {
		t.Errorf("EffectiveRate() error = %v, want errors.Is(err, %v)", err, ErrNonPositiveAmount)
	}
}

// validRecordJSON is validRecord in its wire form. Scrapers and the ingestion
// consumer depend on these names, so changing one is a breaking change.
const validRecordJSON = `{
	"provider": "sendapp",
	"pair": {"base": "CAD", "quote": "NGN"},
	"sendAmount": {"micros": "100000000", "currency": "CAD"},
	"totalPaid": {"micros": "102990000", "currency": "CAD"},
	"receiveAmount": {"micros": "110050000000", "currency": "NGN"},
	"fees": {"micros": "2990000", "currency": "CAD"},
	"offeredRate": 1100.5,
	"payIn": "bank_transfer",
	"payOut": "mobile_money",
	"class": "standard",
	"timestamp": "2023-11-14T22:13:20Z"
}`

func TestRateRecordJSONFieldNames(t *testing.T) {
	got, err := json.Marshal(validRecord())
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	// Compare decoded values so the test pins names and values, not key
	// order or whitespace.
	var gotFields, wantFields map[string]any
	if err := json.Unmarshal(got, &gotFields); err != nil {
		t.Fatalf("decoding marshaled record: %v", err)
	}
	if err := json.Unmarshal([]byte(validRecordJSON), &wantFields); err != nil {
		t.Fatalf("decoding validRecordJSON: %v", err)
	}
	if !reflect.DeepEqual(gotFields, wantFields) {
		t.Errorf("json.Marshal(validRecord()) = %s\nwant %s", got, validRecordJSON)
	}
}

func TestRateRecordJSONRoundTrip(t *testing.T) {
	want := validRecord()
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	var got RateRecord
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal(%s) error = %v", data, err)
	}
	if !got.Timestamp.Equal(want.Timestamp) {
		t.Errorf("round-trip Timestamp = %v, want %v", got.Timestamp, want.Timestamp)
	}
	// time.Time must be compared with Equal, so compare the rest separately.
	got.Timestamp, want.Timestamp = time.Time{}, time.Time{}
	if got != want {
		t.Errorf("round-trip = %+v, want %+v", got, want)
	}
}

// Producers may add fields before consumers are updated, so decoding ignores
// fields it does not know.
func TestRateRecordJSONIgnoresUnknownFields(t *testing.T) {
	in := strings.Replace(validRecordJSON, `"provider": "sendapp",`, `"provider": "sendapp", "addedLater": true,`, 1)
	var got RateRecord
	if err := json.Unmarshal([]byte(in), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v, want nil", err)
	}
	if got.Provider != "sendapp" {
		t.Errorf("Provider = %q, want %q", got.Provider, "sendapp")
	}
}
