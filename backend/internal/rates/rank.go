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
	// 1. Every record must be for the requested comparison.
	for _, r := range records {
		if r.Pair != want.Pair || r.SendAmount != want.SendAmount {
			return nil, nil, fmt.Errorf("rates: quote from %q for %s: %w", r.Provider, r.Pair, ErrMixedComparison)
		}
	}

	// 2. Build quotes: skip invalid records, drop promos the viewer should not see.
	var quotes []Quote
	var skipped []Skipped
	for _, r := range records {
		rate, err := r.EffectiveRate()
		if err != nil {
			skipped = append(skipped, Skipped{Record: r, Err: err})
			continue
		}
		if opts.HidePromos && r.Class != raterecord.ClassStandard {
			continue
		}
		if r.Class == raterecord.ClassPromoFirstTransfer && opts.PromoIneligible[r.Provider] {
			continue
		}
		quotes = append(quotes, Quote{Record: r, EffectiveRate: rate})
	}

	// 3. Best rate first; ties by provider, then pay-in, then pay-out.
	slices.SortFunc(quotes, func(a, b Quote) int {
		return cmp.Or(
			cmp.Compare(b.EffectiveRate, a.EffectiveRate),
			cmp.Compare(a.Record.Provider, b.Record.Provider),
			cmp.Compare(a.Record.PayIn, b.Record.PayIn),
			cmp.Compare(a.Record.PayOut, b.Record.PayOut),
		)
	})

	// 4. Position starts at 1.
	for i := range quotes {
		quotes[i].Position = i + 1
	}

	return quotes, skipped, nil
}
