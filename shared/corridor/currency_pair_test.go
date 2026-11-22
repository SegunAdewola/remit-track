package corridor

import (
	"encoding/json"
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
		{name: "valid uppercase pair", input: "CAD:NGN", want: CurrencyPair{Base: "CAD", Quote: "NGN"}},
		{name: "empty string", input: "", wantErr: true},
		{name: "missing colon", input: "CAD", wantErr: true},
		{name: "too many parts", input: "CAD:NGN:EUR", wantErr: true},
		{name: "base too short", input: "CA:NGN", wantErr: true},
		{name: "quote too long", input: "CAD:NGNX", wantErr: true},
		{name: "non-letter in code", input: "CA1:NGN", wantErr: true},
		{name: "lowercase rejected", input: "cad:ngn", wantErr: true},
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
	c := CurrencyPair{Base: "CAD", Quote: "NGN"}
	if got := c.String(); got != "CAD:NGN" {
		t.Errorf("String() = %q, want %q", got, "CAD:NGN")
	}
}

func TestParseStringRoundTrip(t *testing.T) {
	const s = "CAD:NGN"
	c, err := Parse(s)
	if err != nil {
		t.Fatalf("Parse(%q) unexpected error: %v", s, err)
	}
	if got := c.String(); got != s {
		t.Errorf("round-trip: Parse(%q).String() = %q, want %q", s, got, s)
	}
}

// The JSON form is a wire contract, so the field names are pinned.
func TestCurrencyPairJSON(t *testing.T) {
	c := CurrencyPair{Base: "CAD", Quote: "NGN"}
	const want = `{"base":"CAD","quote":"NGN"}`

	got, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("json.Marshal(%+v) error = %v", c, err)
	}
	if string(got) != want {
		t.Errorf("json.Marshal(%+v) = %s, want %s", c, got, want)
	}

	var back CurrencyPair
	if err := json.Unmarshal([]byte(want), &back); err != nil {
		t.Fatalf("json.Unmarshal(%s) error = %v", want, err)
	}
	if back != c {
		t.Errorf("json.Unmarshal(%s) = %+v, want %+v", want, back, c)
	}
}
