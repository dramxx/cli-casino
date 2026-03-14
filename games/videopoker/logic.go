package videopoker

import (
	"math/rand"
	"sort"
)

type Suit string
type Rank int

const (
	Hearts   Suit = "♥"
	Diamonds Suit = "♦"
	Clubs    Suit = "♣"
	Spades   Suit = "♠"
)

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case King:
		return "K"
	case Queen:
		return "Q"
	case Jack:
		return "J"
	case Ten:
		return "10"
	default:
		return string(rune('0' + int(r)))
	}
}

type Card struct {
	Suit Suit
	Rank Rank
}

type Deck struct {
	Cards []Card
	rng   *rand.Rand
}

func NewDeck(rng *rand.Rand) *Deck {
	suits := []Suit{Hearts, Diamonds, Clubs, Spades}
	ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}

	cards := make([]Card, 0, 52)
	for _, suit := range suits {
		for _, rank := range ranks {
			cards = append(cards, Card{Suit: suit, Rank: rank})
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

type HandRank int

const (
	HighCard HandRank = iota
	JacksOrBetter
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (hr HandRank) String() string {
	switch hr {
	case RoyalFlush:
		return "Royal Flush"
	case StraightFlush:
		return "Straight Flush"
	case FourOfAKind:
		return "Four of a Kind"
	case FullHouse:
		return "Full House"
	case Flush:
		return "Flush"
	case Straight:
		return "Straight"
	case ThreeOfAKind:
		return "Three of a Kind"
	case TwoPair:
		return "Two Pair"
	case JacksOrBetter:
		return "Jacks or Better"
	default:
		return "High Card"
	}
}

func (hr HandRank) Payout(bet float64) float64 {
	multipliers := map[HandRank]float64{
		RoyalFlush:     800,
		StraightFlush:  50,
		FourOfAKind:    25,
		FullHouse:      9,
		Flush:          6,
		Straight:       4,
		ThreeOfAKind:   3,
		TwoPair:        2,
		JacksOrBetter:  1,
	}
	if mult, ok := multipliers[hr]; ok {
		return bet * mult
	}
	return 0
}

func EvaluateHand(cards []Card) HandRank {
	if len(cards) != 5 {
		return HighCard
	}

	ranks := make([]Rank, len(cards))
	for i, card := range cards {
		ranks[i] = card.Rank
	}
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] < ranks[j]
	})

	isFlush := true
	firstSuit := cards[0].Suit
	for _, card := range cards {
		if card.Suit != firstSuit {
			isFlush = false
			break
		}
	}

	isStraight := true
	for i := 1; i < len(ranks); i++ {
		if ranks[i] != ranks[i-1]+1 {
			isStraight = false
			break
		}
	}

	if !isStraight && ranks[0] == Two && ranks[1] == Three && ranks[2] == Four && ranks[3] == Five && ranks[4] == Ace {
		isStraight = true
	}

	rankCounts := make(map[Rank]int)
	for _, rank := range ranks {
		rankCounts[rank]++
	}

	counts := make([]int, 0, len(rankCounts))
	for _, count := range rankCounts {
		counts = append(counts, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	if isStraight && isFlush {
		if ranks[0] == Ten && ranks[4] == Ace {
			return RoyalFlush
		}
		return StraightFlush
	}

	if counts[0] == 4 {
		return FourOfAKind
	}

	if counts[0] == 3 && counts[1] == 2 {
		return FullHouse
	}

	if isFlush {
		return Flush
	}

	if isStraight {
		return Straight
	}

	if counts[0] == 3 {
		return ThreeOfAKind
	}

	if counts[0] == 2 && counts[1] == 2 {
		return TwoPair
	}

	if counts[0] == 2 {
		for rank, count := range rankCounts {
			if count == 2 && rank >= Jack {
				return JacksOrBetter
			}
		}
	}

	return HighCard
}
