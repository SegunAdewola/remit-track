package rates

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/SegunAdewola/remit-track/shared/corridor"
	"github.com/SegunAdewola/remit-track/shared/money"
	"github.com/SegunAdewola/remit-track/shared/raterecord"
)

// Quote is one provider quote in a ranked comparison.
type Quote struct {
	Record        raterecord.RateRecord // the quote as the provider reported it
	EffectiveRate float64               // ReceiveAmount / TotalPaid; higher is better
	Position      int                   // 1 for the best quote
}

// IsPromo reports whether q is a promotional quote rather than the standard
// price, so it can be labeled.
func (q Quote) IsPromo() bool {
	return q.Record.Class != raterecord.ClassStandard
}

// RankOptions controls which quotes a viewer sees. The zero value shows every
// valid quote, promotions included.
type RankOptions struct {
	// HidePromos leaves out every promotional quote. Crawler-facing pages
	// set it because they cannot know whether a viewer qualifies.
	HidePromos bool
	// PromoIneligible is the set of provider slugs whose first-transfer
	// promotion the viewer can no longer get, because they have already
	// sent with that provider. A nil map means eligible for all.
	PromoIneligible map[string]bool
}

// Skipped is a record Rank left out because it failed validation. Such records
// should have been rejected at ingestion, so callers should log them.
type Skipped struct {
	Record raterecord.RateRecord
	Err    error
}

// Comparison identifies what is being compared: quotes for one corridor and
// one send amount.
type Comparison struct {
	Pair       corridor.CurrencyPair
	SendAmount money.Money
}

// ErrMixedComparison means a record passed to Rank is for a different corridor
// or send amount than the one requested, which indicates a bug in the caller.
var ErrMixedComparison = errors.New("quote does not match the requested comparison")

// Rank orders the quotes for want, best effective rate first, and numbers
// them from 1. Promotions are ranked alongside standard quotes unless opts
// hides them. Equal rates are ordered by provider, then pay-in method, then
// pay-out method, so the result does not depend on the order of records.
//
// Invalid records are left out and returned in skipped. If any record is for
// a different corridor or send amount than want, Rank returns an error
// wrapping ErrMixedComparison.
func Rank(want Comparison, records []raterecord.RateRecord, opts RankOptions) ([]Quote, []Skipped, error) {
	if err := checkComparison(want, records); err != nil {
		return nil, nil, err
	}
	quotes, skipped := buildQuotes(records, opts)
	slices.SortFunc(quotes, compareQuotes)
	for i := range quotes {
		quotes[i].Position = i + 1
	}
	return quotes, skipped, nil
}

// checkComparison reports the first record that is not for want.
func checkComparison(want Comparison, records []raterecord.RateRecord) error {
	for _, r := range records {
		if r.Pair != want.Pair || r.SendAmount != want.SendAmount {
			return fmt.Errorf("rates: quote from %q for %s: %w", r.Provider, r.Pair, ErrMixedComparison)
		}
	}
	return nil
}

// buildQuotes turns the valid records the viewer may see into unranked quotes,
// and returns the invalid records separately.
func buildQuotes(records []raterecord.RateRecord, opts RankOptions) ([]Quote, []Skipped) {
	var quotes []Quote
	var skipped []Skipped
	for _, r := range records {
		rate, err := r.EffectiveRate() // also validates r
		if err != nil {
			skipped = append(skipped, Skipped{Record: r, Err: err})
			continue
		}
		if opts.hides(r) {
			continue
		}
		quotes = append(quotes, Quote{Record: r, EffectiveRate: rate})
	}
	return quotes, skipped
}

// hides reports whether the viewer described by o should not see r: every
// promotion when HidePromos is set, and first-transfer promotions from
// providers the viewer has already used.
func (o RankOptions) hides(r raterecord.RateRecord) bool {
	if r.Class == raterecord.ClassStandard {
		return false
	}
	if o.HidePromos {
		return true
	}
	return r.Class == raterecord.ClassPromoFirstTransfer && o.PromoIneligible[r.Provider]
}

// compareQuotes orders quotes best first: higher effective rate, then
// provider, pay-in method and pay-out method alphabetically, so equal rates
// always come out in the same order.
func compareQuotes(a, b Quote) int {
	return cmp.Or(
		cmp.Compare(b.EffectiveRate, a.EffectiveRate), // b first: descending
		cmp.Compare(a.Record.Provider, b.Record.Provider),
		cmp.Compare(a.Record.PayIn, b.Record.PayIn),
		cmp.Compare(a.Record.PayOut, b.Record.PayOut),
	)
}
