package slots

import "math/rand"

type Symbol string

const (
	SEVEN  Symbol = "7"
	BAR    Symbol = "BAR"
	BELL   Symbol = "BELL"
	CHERRY Symbol = "CHERRY"
	LEMON  Symbol = "LEMON"
	GRAPE  Symbol = "GRAPE"
)

var Symbols = []Symbol{SEVEN, BAR, BELL, CHERRY, LEMON, GRAPE}

type SymbolInfo struct {
	Symbol  Symbol
	Display string
	Weight  int
	Payout  float64
}

var SymbolTable = map[Symbol]SymbolInfo{
	SEVEN:  {Symbol: SEVEN, Display: "7", Weight: 1, Payout: 50},
	BAR:    {Symbol: BAR, Display: "BAR", Weight: 3, Payout: 10},
	BELL:   {Symbol: BELL, Display: "BELL", Weight: 5, Payout: 5},
	CHERRY: {Symbol: CHERRY, Display: "CHERRY", Weight: 8, Payout: 3},
	LEMON:  {Symbol: LEMON, Display: "LEMON", Weight: 10, Payout: 2},
	GRAPE:  {Symbol: GRAPE, Display: "GRAPE", Weight: 12, Payout: 1.5},
}

func buildWeightedSymbols() []Symbol {
	var weighted []Symbol
	for _, s := range Symbols {
		info := SymbolTable[s]
		for i := 0; i < info.Weight; i++ {
			weighted = append(weighted, s)
		}
	}
	return weighted
}

var weightedSymbols = buildWeightedSymbols()

type Reels [3]Symbol

func SpinReels(rng *rand.Rand) Reels {
	return Reels{
		weightedSymbols[rng.Intn(len(weightedSymbols))],
		weightedSymbols[rng.Intn(len(weightedSymbols))],
		weightedSymbols[rng.Intn(len(weightedSymbols))],
	}
}

type Result struct {
	Reels      Reels
	Win        bool
	Symbol     Symbol
	MatchCount int
	Payout     float64
	AmountWon  float64
}

func CalculatePayout(reels Reels, bet float64) *Result {
	result := &Result{Reels: reels}

	if reels[0] == reels[1] && reels[1] == reels[2] {
		result.Win = true
		result.Symbol = reels[0]
		result.MatchCount = 3
		info := SymbolTable[reels[0]]
		result.Payout = info.Payout
		result.AmountWon = bet * info.Payout
		return result
	}

	if reels[0] == reels[1] || reels[1] == reels[2] || reels[0] == reels[2] {
		if reels[0] == reels[1] {
			result.Symbol = reels[0]
			result.MatchCount = 2
		} else if reels[1] == reels[2] {
			result.Symbol = reels[1]
			result.MatchCount = 2
		} else {
			result.Symbol = reels[0]
			result.MatchCount = 2
		}
		info := SymbolTable[result.Symbol]
		if info.Payout >= 2 {
			result.Payout = info.Payout / 2
			result.AmountWon = bet * result.Payout
			result.Win = true
		}
	}

	return result
}
