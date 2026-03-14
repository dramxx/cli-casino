package blackjack

import (
	"math/rand"
)

type Suit string
type Rank string

const (
	Hearts   Suit = "♥"
	Diamonds Suit = "♦"
	Clubs    Suit = "♣"
	Spades   Suit = "♠"
)

const (
	Ace   Rank = "A"
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "10"
	Jack  Rank = "J"
	Queen Rank = "Q"
	King  Rank = "K"
)

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) Value() int {
	switch c.Rank {
	case Ace:
		return 11
	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten, Jack, Queen, King:
		return 10
	}
	return 0
}

type Deck struct {
	Cards []Card
	rng   *rand.Rand
}

func NewDeck(rng *rand.Rand, numDecks int) *Deck {
	suits := []Suit{Hearts, Diamonds, Clubs, Spades}
	ranks := []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	cards := make([]Card, 0, 52*numDecks)
	for i := 0; i < numDecks; i++ {
		for _, suit := range suits {
			for _, rank := range ranks {
				cards = append(cards, Card{Suit: suit, Rank: rank})
			}
		}
	}

	d := &Deck{Cards: cards, rng: rng}
	d.Shuffle()
	return d
}

func (d *Deck) Shuffle() {
	d.rng.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Draw() Card {
	if len(d.Cards) == 0 {
		return Card{}
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}

func (d *Deck) NeedsReshuffle() bool {
	return len(d.Cards) < 52
}

type Hand struct {
	Cards []Card
}

func (h *Hand) Add(card Card) {
	h.Cards = append(h.Cards, card)
}

func (h *Hand) Value() int {
	value := 0
	aces := 0

	for _, card := range h.Cards {
		if card.Rank == Ace {
			aces++
		}
		value += card.Value()
	}

	for aces > 0 && value > 21 {
		value -= 10
		aces--
	}

	return value
}

func (h *Hand) IsBust() bool {
	return h.Value() > 21
}

func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}

func (h *Hand) IsSoft() bool {
	value := 0
	hasAce := false

	for _, card := range h.Cards {
		if card.Rank == Ace {
			hasAce = true
		}
		value += card.Value()
	}

	return hasAce && value > 21 && h.Value() <= 21
}

func (h *Hand) CanSplit() bool {
	return len(h.Cards) == 2 && h.Cards[0].Rank == h.Cards[1].Rank
}
