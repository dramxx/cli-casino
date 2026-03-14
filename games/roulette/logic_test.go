package roulette

import (
	"math/rand"
	"testing"
)

func TestSpin(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	
	for i := 0; i < 100; i++ {
		result := Spin(rng)
		if result < 0 || result > 36 {
			t.Errorf("Spin result %d out of range [0, 36]", result)
		}
	}
}

func TestNumberProperties(t *testing.T) {
	if Numbers[0].IsZero != true {
		t.Error("Number 0 should be marked as zero")
	}
	
	if Numbers[0].IsRed || Numbers[0].IsBlack {
		t.Error("Number 0 should not be red or black")
	}
	
	if Numbers[1].IsRed != true {
		t.Error("Number 1 should be red")
	}
	
	if Numbers[2].IsBlack != true {
		t.Error("Number 2 should be black")
	}
	
	if Numbers[1].IsLow != true {
		t.Error("Number 1 should be low")
	}
	
	if Numbers[36].IsHigh != true {
		t.Error("Number 36 should be high")
	}
	
	if Numbers[2].IsEven != true {
		t.Error("Number 2 should be even")
	}
	
	if Numbers[1].IsOdd != true {
		t.Error("Number 1 should be odd")
	}
}

func TestBetWinsStraight(t *testing.T) {
	bet := Bet{Type: BetStraight, Number: 17, Amount: 10}
	
	if !bet.Wins(17) {
		t.Error("Straight bet on 17 should win on 17")
	}
	
	if bet.Wins(18) {
		t.Error("Straight bet on 17 should not win on 18")
	}
}

func TestBetWinsRed(t *testing.T) {
	bet := Bet{Type: BetRed, Amount: 10}
	
	if !bet.Wins(1) {
		t.Error("Red bet should win on 1 (red)")
	}
	
	if bet.Wins(2) {
		t.Error("Red bet should not win on 2 (black)")
	}
	
	if bet.Wins(0) {
		t.Error("Red bet should not win on 0")
	}
}

func TestBetWinsBlack(t *testing.T) {
	bet := Bet{Type: BetBlack, Amount: 10}
	
	if !bet.Wins(2) {
		t.Error("Black bet should win on 2 (black)")
	}
	
	if bet.Wins(1) {
		t.Error("Black bet should not win on 1 (red)")
	}
}

func TestBetWinsEven(t *testing.T) {
	bet := Bet{Type: BetEven, Amount: 10}
	
	if !bet.Wins(2) {
		t.Error("Even bet should win on 2")
	}
	
	if bet.Wins(1) {
		t.Error("Even bet should not win on 1 (odd)")
	}
	
	if bet.Wins(0) {
		t.Error("Even bet should not win on 0")
	}
}

func TestBetWinsOdd(t *testing.T) {
	bet := Bet{Type: BetOdd, Amount: 10}
	
	if !bet.Wins(1) {
		t.Error("Odd bet should win on 1")
	}
	
	if bet.Wins(2) {
		t.Error("Odd bet should not win on 2 (even)")
	}
}

func TestBetWinsLow(t *testing.T) {
	bet := Bet{Type: BetLow, Amount: 10}
	
	if !bet.Wins(1) {
		t.Error("Low bet should win on 1")
	}
	
	if !bet.Wins(18) {
		t.Error("Low bet should win on 18")
	}
	
	if bet.Wins(19) {
		t.Error("Low bet should not win on 19")
	}
	
	if bet.Wins(0) {
		t.Error("Low bet should not win on 0")
	}
}

func TestBetWinsHigh(t *testing.T) {
	bet := Bet{Type: BetHigh, Amount: 10}
	
	if !bet.Wins(19) {
		t.Error("High bet should win on 19")
	}
	
	if !bet.Wins(36) {
		t.Error("High bet should win on 36")
	}
	
	if bet.Wins(18) {
		t.Error("High bet should not win on 18")
	}
}

func TestBetWinsDozens(t *testing.T) {
	bet1 := Bet{Type: BetDozen1, Amount: 10}
	bet2 := Bet{Type: BetDozen2, Amount: 10}
	bet3 := Bet{Type: BetDozen3, Amount: 10}
	
	if !bet1.Wins(1) || !bet1.Wins(12) {
		t.Error("Dozen1 bet should win on 1-12")
	}
	
	if bet1.Wins(13) {
		t.Error("Dozen1 bet should not win on 13")
	}
	
	if !bet2.Wins(13) || !bet2.Wins(24) {
		t.Error("Dozen2 bet should win on 13-24")
	}
	
	if !bet3.Wins(25) || !bet3.Wins(36) {
		t.Error("Dozen3 bet should win on 25-36")
	}
}

func TestBetPayout(t *testing.T) {
	tests := []struct {
		betType  BetType
		amount   float64
		expected float64
	}{
		{BetStraight, 10, 350},
		{BetRed, 10, 10},
		{BetBlack, 10, 10},
		{BetEven, 10, 10},
		{BetOdd, 10, 10},
		{BetLow, 10, 10},
		{BetHigh, 10, 10},
		{BetDozen1, 10, 20},
		{BetDozen2, 10, 20},
		{BetDozen3, 10, 20},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.betType), func(t *testing.T) {
			bet := Bet{Type: tt.betType, Amount: tt.amount}
			payout := bet.Payout()
			if payout != tt.expected {
				t.Errorf("Expected payout %.2f, got %.2f", tt.expected, payout)
			}
		})
	}
}

func TestRedNumbers(t *testing.T) {
	expectedReds := []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
	
	if len(RedNumbers) != len(expectedReds) {
		t.Errorf("Expected %d red numbers, got %d", len(expectedReds), len(RedNumbers))
	}
	
	for _, red := range expectedReds {
		found := false
		for _, r := range RedNumbers {
			if r == red {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Number %d should be red", red)
		}
	}
}

func TestBlackNumbers(t *testing.T) {
	expectedBlacks := []int{2, 4, 6, 8, 10, 11, 13, 15, 17, 20, 22, 24, 26, 28, 29, 31, 33, 35}
	
	if len(BlackNumbers) != len(expectedBlacks) {
		t.Errorf("Expected %d black numbers, got %d", len(expectedBlacks), len(BlackNumbers))
	}
	
	for _, black := range expectedBlacks {
		found := false
		for _, b := range BlackNumbers {
			if b == black {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Number %d should be black", black)
		}
	}
}
