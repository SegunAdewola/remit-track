package money

// minorUnits maps each supported currency code to its ISO 4217 minor unit,
// the number of decimal places in its minor unit. Stored amounts are read
// through this table, so an existing value must not be changed. Add a row when
// a corridor adds a currency.
var minorUnits = map[string]int{
	"CAD": 2,
	"USD": 2,
	"GBP": 2,
	"EUR": 2,
	"NGN": 2,
	"GHS": 2,
	"KES": 2,
	"XOF": 0,
	"XAF": 0,
}

// MinorUnits returns the number of decimal places in the minor unit of the
// currency with the given ISO 4217 code: 2 for CAD (cents), 0 for XOF. The
// boolean is false if the currency is not supported.
func MinorUnits(code string) (int, bool) {
	units, ok := minorUnits[code]
	return units, ok
}
