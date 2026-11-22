package rates

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/SegunAdewola/remit-track/shared/corridor"
	"github.com/SegunAdewola/remit-track/shared/money"
	"github.com/SegunAdewola/remit-track/shared/raterecord"
)

// cadToNGN100 is the comparison every test asks for: CAD 100.00 sent to NGN.
var cadToNGN100 = Comparison{
	Pair:       corridor.CurrencyPair{Base: "CAD", Quote: "NGN"},
	SendAmount: money.Money{Micros: 100_000_000, Currency: "CAD"},
}

// quote returns a valid CAD 100.00 to NGN quote in which the sender pays
// CAD 100.00 in total, so the effective rate is receiveNaira / 100.
func quote(provider, payIn, payOut string, class raterecord.RateClass, receiveNaira int64) raterecord.RateRecord {
	return raterecord.RateRecord{
		Provider:      provider,
		Pair:          corridor.CurrencyPair{Base: "CAD", Quote: "NGN"},
		SendAmount:    money.Money{Micros: 100_000_000, Currency: "CAD"},
		TotalPaid:     money.Money{Micros: 100_000_000, Currency: "CAD"},
		ReceiveAmount: money.Money{Micros: receiveNaira * money.MicrosPerUnit, Currency: "NGN"},
		Fees:          money.Money{Micros: 0, Currency: "CAD"},
		OfferedRate:   float64(receiveNaira) / 100,
		PayIn:         raterecord.PayMethod(payIn),
		PayOut:        raterecord.PayMethod(payOut),
		Class:         class,
		Timestamp:     time.Unix(1_700_000_000, 0).UTC(),
	}
}

// providers returns the provider of each ranked quote, in order.
func providers(qs []Quote) []string {
	var out []string
	for _, q := range qs {
		out = append(out, q.Record.Provider)
	}
	return out
}

func TestRankOrdersByEffectiveRate(t *testing.T) {
	records := []raterecord.RateRecord{
		quote("lemfi", "bank_transfer", "bank", raterecord.ClassStandard, 105_000),
		quote("sendapp", "bank_transfer", "bank", raterecord.ClassStandard, 110_000),
		quote("remitly", "bank_transfer", "bank", raterecord.ClassStandard, 100_000),
	}
	ranked, skipped, err := Rank(cadToNGN100, records, RankOptions{})
	if err != nil || len(skipped) != 0 {
		t.Fatalf("Rank() skipped = %v, err = %v, want none", skipped, err)
	}
	if got, want := providers(ranked), []string{"sendapp", "lemfi", "remitly"}; !slices.Equal(got, want) {
		t.Fatalf("Rank() order = %v, want %v", got, want)
	}
	for i, q := range ranked {
		if q.Position != i+1 {
			t.Errorf("ranked[%d].Position = %d, want %d", i, q.Position, i+1)
		}
	}
	if ranked[0].EffectiveRate != 1100 {
		t.Errorf("ranked[0].EffectiveRate = %v, want 1100", ranked[0].EffectiveRate)
	}
}

// Equal rates are ordered by provider, then pay-in, then pay-out, whatever the
// input order, so the table never reshuffles between requests.
func TestRankTieBreaksIgnoreInputOrder(t *testing.T) {
	records := []raterecord.RateRecord{
		quote("zeta", "bank_transfer", "bank", raterecord.ClassStandard, 100_000),
		quote("alpha", "card", "mobile_money", raterecord.ClassStandard, 100_000),
		quote("alpha", "bank_transfer", "mobile_money", raterecord.ClassStandard, 100_000),
		quote("alpha", "bank_transfer", "bank", raterecord.ClassStandard, 100_000),
	}
	want := []string{
		"alpha/bank_transfer/bank",
		"alpha/bank_transfer/mobile_money",
		"alpha/card/mobile_money",
		"zeta/bank_transfer/bank",
	}
	reversed := slices.Clone(records)
	slices.Reverse(reversed)
	for _, input := range [][]raterecord.RateRecord{records, reversed} {
		ranked, _, err := Rank(cadToNGN100, input, RankOptions{})
		if err != nil {
			t.Fatalf("Rank() error = %v", err)
		}
		var got []string
		for _, q := range ranked {
			got = append(got, q.Record.Provider+"/"+string(q.Record.PayIn)+"/"+string(q.Record.PayOut))
		}
		if !slices.Equal(got, want) {
			t.Errorf("Rank() order = %v, want %v", got, want)
		}
	}
}

func TestRankPromos(t *testing.T) {
	records := []raterecord.RateRecord{
		quote("sendapp", "bank_transfer", "bank", raterecord.ClassStandard, 100_000),
		quote("sendapp", "bank_transfer", "bank", raterecord.ClassPromoFirstTransfer, 115_000),
		quote("lemfi", "bank_transfer", "bank", raterecord.ClassPromoFirstTransfer, 112_000),
		quote("lemfi", "bank_transfer", "bank", raterecord.ClassStandard, 105_000),
	}
	tests := []struct {
		name string
		opts RankOptions
		want []string // provider + class, best first
	}{
		{
			name: "promos ranked with standard quotes",
			opts: RankOptions{},
			want: []string{"sendapp promo_first_transfer", "lemfi promo_first_transfer", "lemfi standard", "sendapp standard"},
		},
		{
			name: "hide all promos",
			opts: RankOptions{HidePromos: true},
			want: []string{"lemfi standard", "sendapp standard"},
		},
		{
			// The viewer has already sent with sendapp, so its first-transfer
			// promo is not available to them; lemfi's still is.
			name: "promo ineligible for one provider",
			opts: RankOptions{PromoIneligible: map[string]bool{"sendapp": true}},
			want: []string{"lemfi promo_first_transfer", "lemfi standard", "sendapp standard"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranked, _, err := Rank(cadToNGN100, records, tt.opts)
			if err != nil {
				t.Fatalf("Rank() error = %v", err)
			}
			var got []string
			for _, q := range ranked {
				got = append(got, q.Record.Provider+" "+string(q.Record.Class))
				if q.IsPromo() != (q.Record.Class != raterecord.ClassStandard) {
					t.Errorf("%s %s: IsPromo() = %v", q.Record.Provider, q.Record.Class, q.IsPromo())
				}
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("Rank() = %v, want %v", got, tt.want)
			}
		})
	}
}

// An invalid record means a bug upstream, but one bad row must not take down
// the comparison: it is skipped and returned with its error.
func TestRankSkipsInvalidRecords(t *testing.T) {
	bad := quote("brokenpay", "bank_transfer", "bank", raterecord.ClassStandard, 120_000)
	bad.TotalPaid.Micros = 0
	records := []raterecord.RateRecord{
		quote("sendapp", "bank_transfer", "bank", raterecord.ClassStandard, 110_000),
		bad,
	}
	ranked, skipped, err := Rank(cadToNGN100, records, RankOptions{})
	if err != nil {
		t.Fatalf("Rank() error = %v, want nil", err)
	}
	if got := providers(ranked); !slices.Equal(got, []string{"sendapp"}) {
		t.Errorf("Rank() ranked = %v, want [sendapp]", got)
	}
	if len(skipped) != 1 || skipped[0].Record.Provider != "brokenpay" {
		t.Fatalf("Rank() skipped = %+v, want only brokenpay", skipped)
	}
	if !errors.Is(skipped[0].Err, raterecord.ErrNonPositiveAmount) {
		t.Errorf("skipped[0].Err = %v, want errors.Is(err, raterecord.ErrNonPositiveAmount)", skipped[0].Err)
	}
}

// Rank compares quotes for the requested corridor and send amount. A record
// for anything else means the caller queried the wrong data.
func TestRankRejectsMixedComparisons(t *testing.T) {
	base := quote("sendapp", "bank_transfer", "bank", raterecord.ClassStandard, 110_000)

	otherPair := quote("lemfi", "bank_transfer", "bank", raterecord.ClassStandard, 105_000)
	otherPair.Pair.Quote = "GHS"
	otherPair.ReceiveAmount.Currency = "GHS"

	otherAmount := quote("lemfi", "bank_transfer", "bank", raterecord.ClassStandard, 105_000)
	otherAmount.SendAmount.Micros = 200_000_000
	otherAmount.TotalPaid.Micros = 200_000_000

	for name, other := range map[string]raterecord.RateRecord{"different corridor": otherPair, "different send amount": otherAmount} {
		t.Run(name, func(t *testing.T) {
			if _, _, err := Rank(cadToNGN100, []raterecord.RateRecord{base, other}, RankOptions{}); !errors.Is(err, ErrMixedComparison) {
				t.Errorf("Rank() error = %v, want errors.Is(err, ErrMixedComparison)", err)
			}
		})
	}
}

// Every record must match the requested comparison. Inferring the target
// from the records would miss a caller that passed only wrong records.
func TestRankRejectsRecordsAllForAnotherCorridor(t *testing.T) {
	ghs := quote("sendapp", "bank_transfer", "mobile_money", raterecord.ClassStandard, 1_000)
	ghs.Pair.Quote = "GHS"
	ghs.ReceiveAmount.Currency = "GHS"
	if _, _, err := Rank(cadToNGN100, []raterecord.RateRecord{ghs}, RankOptions{}); !errors.Is(err, ErrMixedComparison) {
		t.Errorf("Rank() error = %v, want errors.Is(err, ErrMixedComparison)", err)
	}
}

func TestRankEmpty(t *testing.T) {
	ranked, skipped, err := Rank(cadToNGN100, nil, RankOptions{})
	if err != nil || len(ranked) != 0 || len(skipped) != 0 {
		t.Errorf("Rank(nil) = %v, %v, %v; want empty and nil error", ranked, skipped, err)
	}
}

// FuzzRank generates random quote lists and options from a seed and checks
// properties that must hold for any input: best rate first, positions from 1,
// hidden promos absent, invalid records skipped, every visible valid record
// ranked exactly once, and an order that does not depend on input order.
func FuzzRank(f *testing.F) {
	f.Add(uint64(1), uint8(10))
	f.Add(uint64(42), uint8(0))
	f.Add(uint64(7), uint8(60))
	f.Fuzz(func(t *testing.T, seed uint64, n uint8) {
		rng := rand.New(rand.NewPCG(seed, seed))
		names := []string{"alpha", "beta", "gamma"}
		methods := []string{"bank", "card", "mobile_money"}

		opts := RankOptions{HidePromos: rng.IntN(4) == 0, PromoIneligible: map[string]bool{}}
		for _, p := range names {
			if rng.IntN(2) == 0 {
				opts.PromoIneligible[p] = true
			}
		}

		var records []raterecord.RateRecord
		wantRanked, wantSkipped := 0, 0
		for range int(n % 40) {
			class := raterecord.ClassStandard
			if rng.IntN(3) == 0 {
				class = raterecord.ClassPromoFirstTransfer
			}
			r := quote(names[rng.IntN(len(names))], methods[rng.IntN(len(methods))],
				methods[rng.IntN(len(methods))], class, int64(rng.IntN(200_000)+1))
			hidden := class != raterecord.ClassStandard && (opts.HidePromos || opts.PromoIneligible[r.Provider])
			switch {
			case rng.IntN(10) == 0:
				r.TotalPaid.Micros = 0 // invalid
				wantSkipped++
			case !hidden:
				wantRanked++
			}
			records = append(records, r)
		}

		ranked, skipped, err := Rank(cadToNGN100, records, opts)
		if err != nil {
			t.Fatalf("Rank() error = %v", err)
		}
		if len(ranked) != wantRanked || len(skipped) != wantSkipped {
			t.Fatalf("Rank() ranked %d, skipped %d; want %d and %d", len(ranked), len(skipped), wantRanked, wantSkipped)
		}
		for i, q := range ranked {
			if q.Position != i+1 {
				t.Errorf("ranked[%d].Position = %d, want %d", i, q.Position, i+1)
			}
			if i > 0 && q.EffectiveRate > ranked[i-1].EffectiveRate {
				t.Errorf("ranked[%d] rate %v is better than ranked[%d] rate %v", i, q.EffectiveRate, i-1, ranked[i-1].EffectiveRate)
			}
			if q.IsPromo() && (opts.HidePromos || opts.PromoIneligible[q.Record.Provider]) {
				t.Errorf("ranked[%d] is a promo the viewer should not see: %+v", i, q.Record)
			}
		}
		for _, s := range skipped {
			if s.Err == nil {
				t.Errorf("skipped record %+v has nil Err", s.Record)
			}
		}

		shuffled := slices.Clone(records)
		rng.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		again, _, err := Rank(cadToNGN100, shuffled, opts)
		if err != nil {
			t.Fatalf("Rank(shuffled) error = %v", err)
		}
		if !slices.Equal(again, ranked) {
			t.Errorf("Rank() order depends on input order:\n first  %v\n second %v", providers(ranked), providers(again))
		}
	})
}
