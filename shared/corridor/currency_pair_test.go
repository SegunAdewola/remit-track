package corridor

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    CurrencyPair
		wantErr bool
	}{
		{name: "valid uppercase pair", input: "USD:GHS", want: CurrencyPair{Base: "USD", Quote: "GHS"}},
		{name: "empty string", input: "", wantErr: true},
		{name: "missing colon", input: "USD", wantErr: true},
		{name: "too many parts", input: "USD:GHS:EUR", wantErr: true},
		{name: "base too short", input: "US:GHS", wantErr: true},
		{name: "quote too long", input: "USD:GHSX", wantErr: true},
		{name: "non-letter in code", input: "US1:GHS", wantErr: true},
		{name: "lowercase rejected", input: "usd:ghs", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q) = %+v, want error", tt.input, got)
				}
				if !errors.Is(err, ErrInvalidPair) {
					t.Errorf("Parse(%q) error = %v, want errors.Is(err, ErrInvalidPair)", tt.input, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestCurrencyPairString(t *testing.T) {
	c := CurrencyPair{Base: "USD", Quote: "GHS"}
	if got := c.String(); got != "USD:GHS" {
		t.Errorf("String() = %q, want %q", got, "USD:GHS")
	}
}

func TestParseStringRoundTrip(t *testing.T) {
	const s = "USD:GHS"
	c, err := Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", s, err)
	}
	if got := c.String(); got != s {
		t.Errorf("round-trip: Parse(%q).String() = %q, want %q", s, got, s)
	}
}
