package videopoker

import (
	"math/rand"
	"testing"
)

func TestEvaluateHandRoyalFlush(t *testing.T) {
	cards := []Card{
		{Rank: Ten, Suit: Hearts},
		{Rank: Jack, Suit: Hearts},
		{Rank: Queen, Suit: Hearts},
		{Rank: King, Suit: Hearts},
		{Rank: Ace, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != RoyalFlush {
		t.Errorf("Expected RoyalFlush, got %v", rank)
	}
}

func TestEvaluateHandStraightFlush(t *testing.T) {
	cards := []Card{
		{Rank: Five, Suit: Spades},
		{Rank: Six, Suit: Spades},
		{Rank: Seven, Suit: Spades},
		{Rank: Eight, Suit: Spades},
		{Rank: Nine, Suit: Spades},
	}
	
	rank := EvaluateHand(cards)
	if rank != StraightFlush {
		t.Errorf("Expected StraightFlush, got %v", rank)
	}
}

func TestEvaluateHandFourOfAKind(t *testing.T) {
	cards := []Card{
		{Rank: Seven, Suit: Hearts},
		{Rank: Seven, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
		{Rank: Seven, Suit: Spades},
		{Rank: King, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != FourOfAKind {
		t.Errorf("Expected FourOfAKind, got %v", rank)
	}
}

func TestEvaluateHandFullHouse(t *testing.T) {
	cards := []Card{
		{Rank: Three, Suit: Hearts},
		{Rank: Three, Suit: Diamonds},
		{Rank: Three, Suit: Clubs},
		{Rank: Six, Suit: Spades},
		{Rank: Six, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != FullHouse {
		t.Errorf("Expected FullHouse, got %v", rank)
	}
}

func TestEvaluateHandFlush(t *testing.T) {
	cards := []Card{
		{Rank: Two, Suit: Clubs},
		{Rank: Four, Suit: Clubs},
		{Rank: Seven, Suit: Clubs},
		{Rank: Nine, Suit: Clubs},
		{Rank: King, Suit: Clubs},
	}
	
	rank := EvaluateHand(cards)
	if rank != Flush {
		t.Errorf("Expected Flush, got %v", rank)
	}
}

func TestEvaluateHandStraight(t *testing.T) {
	cards := []Card{
		{Rank: Five, Suit: Hearts},
		{Rank: Six, Suit: Diamonds},
		{Rank: Seven, Suit: Clubs},
		{Rank: Eight, Suit: Spades},
		{Rank: Nine, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != Straight {
		t.Errorf("Expected Straight, got %v", rank)
	}
}

func TestEvaluateHandThreeOfAKind(t *testing.T) {
	cards := []Card{
		{Rank: Queen, Suit: Hearts},
		{Rank: Queen, Suit: Diamonds},
		{Rank: Queen, Suit: Clubs},
		{Rank: Five, Suit: Spades},
		{Rank: Nine, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != ThreeOfAKind {
		t.Errorf("Expected ThreeOfAKind, got %v", rank)
	}
}

func TestEvaluateHandTwoPair(t *testing.T) {
	cards := []Card{
		{Rank: Jack, Suit: Hearts},
		{Rank: Jack, Suit: Diamonds},
		{Rank: Four, Suit: Clubs},
		{Rank: Four, Suit: Spades},
		{Rank: Ace, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != TwoPair {
		t.Errorf("Expected TwoPair, got %v", rank)
	}
}

func TestEvaluateHandJacksOrBetter(t *testing.T) {
	cards := []Card{
		{Rank: King, Suit: Hearts},
		{Rank: King, Suit: Diamonds},
		{Rank: Three, Suit: Clubs},
		{Rank: Seven, Suit: Spades},
		{Rank: Nine, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != JacksOrBetter {
		t.Errorf("Expected JacksOrBetter, got %v", rank)
	}
}

func TestEvaluateHandHighCard(t *testing.T) {
	cards := []Card{
		{Rank: Two, Suit: Hearts},
		{Rank: Two, Suit: Diamonds},
		{Rank: Five, Suit: Clubs},
		{Rank: Eight, Suit: Spades},
		{Rank: King, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != HighCard {
		t.Errorf("Expected HighCard, got %v", rank)
	}
}

func TestHandRankPayout(t *testing.T) {
	bet := 5.0
	
	tests := []struct {
		rank     HandRank
		expected float64
	}{
		{RoyalFlush, 4000.0},
		{StraightFlush, 250.0},
		{FourOfAKind, 125.0},
		{FullHouse, 45.0},
		{Flush, 30.0},
		{Straight, 20.0},
		{ThreeOfAKind, 15.0},
		{TwoPair, 10.0},
		{JacksOrBetter, 5.0},
		{HighCard, 0.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.rank.String(), func(t *testing.T) {
			payout := tt.rank.Payout(bet)
			if payout != tt.expected {
				t.Errorf("Expected payout %.2f, got %.2f", tt.expected, payout)
			}
		})
	}
}

func TestDeckCreation(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	deck := NewDeck(rng)
	
	if len(deck.Cards) != 52 {
		t.Errorf("Expected 52 cards, got %d", len(deck.Cards))
	}
}

func TestDeckDraw(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	deck := NewDeck(rng)
	
	card := deck.Draw()
	if len(deck.Cards) != 51 {
		t.Errorf("Expected 51 cards after draw, got %d", len(deck.Cards))
	}
	
	if card.Rank == 0 {
		t.Error("Drew invalid card")
	}
}

func TestWheelStraight(t *testing.T) {
	cards := []Card{
		{Rank: Ace, Suit: Hearts},
		{Rank: Two, Suit: Diamonds},
		{Rank: Three, Suit: Clubs},
		{Rank: Four, Suit: Spades},
		{Rank: Five, Suit: Hearts},
	}
	
	rank := EvaluateHand(cards)
	if rank != Straight {
		t.Errorf("Expected Straight (wheel), got %v", rank)
	}
}
