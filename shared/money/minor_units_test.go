package money

import "testing"

// Each value is pinned to ISO 4217, so changing the table also requires
// changing this test.
func TestMinorUnits(t *testing.T) {
	tests := []struct {
		code   string
		want   int
		wantOK bool
	}{
		{code: "CAD", want: 2, wantOK: true},
		{code: "USD", want: 2, wantOK: true},
		{code: "GBP", want: 2, wantOK: true},
		{code: "EUR", want: 2, wantOK: true},
		{code: "NGN", want: 2, wantOK: true},
		{code: "GHS", want: 2, wantOK: true},
		{code: "KES", want: 2, wantOK: true},
		{code: "XOF", want: 0, wantOK: true},
		{code: "XAF", want: 0, wantOK: true},

		// Not in the table.
		{code: "JPY", wantOK: false},
		{code: "ZZZ", wantOK: false},
		{code: "cad", wantOK: false},
		{code: "", wantOK: false},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			got, ok := MinorUnits(tt.code)
			if ok != tt.wantOK {
				t.Fatalf("MinorUnits(%q) ok = %v, want %v", tt.code, ok, tt.wantOK)
			}
			if ok && got != tt.want {
				t.Errorf("MinorUnits(%q) = %d, want %d", tt.code, got, tt.want)
			}
		})
	}
}
