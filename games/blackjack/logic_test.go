package blackjack

import (
	"math/rand"
	"testing"
)

func TestCardValue(t *testing.T) {
	tests := []struct {
		card     Card
		expected int
	}{
		{Card{Rank: Ace, Suit: Hearts}, 11},
		{Card{Rank: Two, Suit: Spades}, 2},
		{Card{Rank: Ten, Suit: Clubs}, 10},
		{Card{Rank: Jack, Suit: Diamonds}, 10},
		{Card{Rank: Queen, Suit: Hearts}, 10},
		{Card{Rank: King, Suit: Spades}, 10},
	}
	
	for _, tt := range tests {
		if got := tt.card.Value(); got != tt.expected {
			t.Errorf("Card %v %v: expected %d, got %d", tt.card.Rank, tt.card.Suit, tt.expected, got)
		}
	}
}

func TestHandValue(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected int
	}{
		{
			name:     "Simple hand",
			cards:    []Card{{Rank: Ten, Suit: Hearts}, {Rank: Five, Suit: Spades}},
			expected: 15,
		},
		{
			name:     "Blackjack",
			cards:    []Card{{Rank: Ace, Suit: Hearts}, {Rank: King, Suit: Spades}},
			expected: 21,
		},
		{
			name:     "Soft hand",
			cards:    []Card{{Rank: Ace, Suit: Hearts}, {Rank: Six, Suit: Spades}},
			expected: 17,
		},
		{
			name:     "Multiple aces",
			cards:    []Card{{Rank: Ace, Suit: Hearts}, {Rank: Ace, Suit: Spades}, {Rank: Nine, Suit: Clubs}},
			expected: 21,
		},
		{
			name:     "Bust with ace adjustment",
			cards:    []Card{{Rank: Ace, Suit: Hearts}, {Rank: King, Suit: Spades}, {Rank: Five, Suit: Clubs}},
			expected: 16,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := Hand{Cards: tt.cards}
			if got := hand.Value(); got != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestHandIsBust(t *testing.T) {
	hand := Hand{Cards: []Card{
		{Rank: King, Suit: Hearts},
		{Rank: Queen, Suit: Spades},
		{Rank: Five, Suit: Clubs},
	}}
	
	if !hand.IsBust() {
		t.Error("Expected hand to be bust")
	}
}

func TestHandIsBlackjack(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected bool
	}{
		{
			name:     "Blackjack",
			cards:    []Card{{Rank: Ace, Suit: Hearts}, {Rank: King, Suit: Spades}},
			expected: true,
		},
		{
			name:     "21 but not blackjack",
			cards:    []Card{{Rank: Seven, Suit: Hearts}, {Rank: Seven, Suit: Spades}, {Rank: Seven, Suit: Clubs}},
			expected: false,
		},
		{
			name:     "Not 21",
			cards:    []Card{{Rank: Ten, Suit: Hearts}, {Rank: Nine, Suit: Spades}},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := Hand{Cards: tt.cards}
			if got := hand.IsBlackjack(); got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestDeckCreation(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	deck := NewDeck(rng, 6)
	
	expectedCards := 52 * 6
	if len(deck.Cards) != expectedCards {
		t.Errorf("Expected %d cards, got %d", expectedCards, len(deck.Cards))
	}
}

func TestDeckDraw(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	deck := NewDeck(rng, 1)
	
	initialCount := len(deck.Cards)
	card := deck.Draw()
	
	if len(deck.Cards) != initialCount-1 {
		t.Errorf("Expected %d cards after draw, got %d", initialCount-1, len(deck.Cards))
	}
	
	if card.Rank == "" {
		t.Error("Drew empty card")
	}
}

func TestHandCanSplit(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected bool
	}{
		{
			name:     "Can split",
			cards:    []Card{{Rank: Eight, Suit: Hearts}, {Rank: Eight, Suit: Spades}},
			expected: true,
		},
		{
			name:     "Cannot split different ranks",
			cards:    []Card{{Rank: Eight, Suit: Hearts}, {Rank: Nine, Suit: Spades}},
			expected: false,
		},
		{
			name:     "Cannot split with more than 2 cards",
			cards:    []Card{{Rank: Eight, Suit: Hearts}, {Rank: Eight, Suit: Spades}, {Rank: Eight, Suit: Clubs}},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hand := Hand{Cards: tt.cards}
			if got := hand.CanSplit(); got != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
		})
	}
}
