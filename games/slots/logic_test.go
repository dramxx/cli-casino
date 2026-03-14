package slots

import (
	"math/rand"
	"testing"
)

func TestSpinReels(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	reels := SpinReels(rng)
	
	if len(reels) != 3 {
		t.Errorf("Expected 3 reels, got %d", len(reels))
	}
	
	for i, symbol := range reels {
		if symbol == "" {
			t.Errorf("Reel %d is empty", i)
		}
	}
}

func TestCalculatePayoutThreeMatch(t *testing.T) {
	reels := Reels{SEVEN, SEVEN, SEVEN}
	bet := 10.0
	result := CalculatePayout(reels, bet)
	
	if !result.Win {
		t.Error("Expected win for three matching sevens")
	}
	
	if result.Symbol != SEVEN {
		t.Errorf("Expected symbol SEVEN, got %v", result.Symbol)
	}
	
	if result.MatchCount != 3 {
		t.Errorf("Expected match count 3, got %d", result.MatchCount)
	}
	
	expectedPayout := bet * SymbolTable[SEVEN].Payout
	if result.AmountWon != expectedPayout {
		t.Errorf("Expected payout %.2f, got %.2f", expectedPayout, result.AmountWon)
	}
}

func TestCalculatePayoutTwoMatch(t *testing.T) {
	reels := Reels{BAR, BAR, LEMON}
	bet := 10.0
	result := CalculatePayout(reels, bet)
	
	if !result.Win {
		t.Error("Expected win for two matching bars")
	}
	
	if result.MatchCount != 2 {
		t.Errorf("Expected match count 2, got %d", result.MatchCount)
	}
}

func TestCalculatePayoutNoMatch(t *testing.T) {
	reels := Reels{SEVEN, BAR, LEMON}
	bet := 10.0
	result := CalculatePayout(reels, bet)
	
	if result.Win {
		t.Error("Expected no win for non-matching symbols")
	}
	
	if result.AmountWon != 0 {
		t.Errorf("Expected 0 payout, got %.2f", result.AmountWon)
	}
}

func TestSymbolWeights(t *testing.T) {
	weighted := buildWeightedSymbols()
	
	sevenCount := 0
	grapeCount := 0
	
	for _, sym := range weighted {
		if sym == SEVEN {
			sevenCount++
		}
		if sym == GRAPE {
			grapeCount++
		}
	}
	
	if sevenCount != SymbolTable[SEVEN].Weight {
		t.Errorf("Expected %d sevens, got %d", SymbolTable[SEVEN].Weight, sevenCount)
	}
	
	if grapeCount != SymbolTable[GRAPE].Weight {
		t.Errorf("Expected %d grapes, got %d", SymbolTable[GRAPE].Weight, grapeCount)
	}
}
