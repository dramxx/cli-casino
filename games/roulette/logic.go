package roulette

import "math/rand"

type Number struct {
	Value   int
	IsZero  bool
	IsRed   bool
	IsBlack bool
	IsEven  bool
	IsOdd   bool
	IsLow   bool
	IsHigh  bool
}

var Numbers []Number

var RedNumbers = []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
var BlackNumbers = []int{2, 4, 6, 8, 10, 11, 13, 15, 17, 20, 22, 24, 26, 28, 29, 31, 33, 35}

func init() {
	Numbers = make([]Number, 37)
	Numbers[0] = Number{Value: 0, IsZero: true}

	for i := 1; i <= 36; i++ {
		n := Number{Value: i}
		n.IsRed = contains(RedNumbers, i)
		n.IsBlack = contains(BlackNumbers, i)
		n.IsEven = i%2 == 0
		n.IsOdd = i%2 == 1
		n.IsLow = i >= 1 && i <= 18
		n.IsHigh = i >= 19 && i <= 36
		Numbers[i] = n
	}
}

func contains(slice []int, val int) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

type BetType string

const (
	BetStraight BetType = "straight"
	BetRed      BetType = "red"
	BetBlack    BetType = "black"
	BetEven     BetType = "even"
	BetOdd      BetType = "odd"
	BetLow      BetType = "low"
	BetHigh     BetType = "high"
	BetDozen1   BetType = "dozen1"
	BetDozen2   BetType = "dozen2"
	BetDozen3   BetType = "dozen3"
)

type Bet struct {
	Type   BetType
	Amount float64
	Number int
}

var Payouts = map[BetType]float64{
	BetStraight: 35,
	BetRed:      1,
	BetBlack:    1,
	BetEven:     1,
	BetOdd:      1,
	BetLow:      1,
	BetHigh:     1,
	BetDozen1:   2,
	BetDozen2:   2,
	BetDozen3:   2,
}

func (b *Bet) Wins(number int) bool {
	switch b.Type {
	case BetStraight:
		return b.Number == number
	case BetRed:
		return Numbers[number].IsRed
	case BetBlack:
		return Numbers[number].IsBlack
	case BetEven:
		return Numbers[number].IsEven && !Numbers[number].IsZero
	case BetOdd:
		return Numbers[number].IsOdd && !Numbers[number].IsZero
	case BetLow:
		return Numbers[number].IsLow && !Numbers[number].IsZero
	case BetHigh:
		return Numbers[number].IsHigh && !Numbers[number].IsZero
	case BetDozen1:
		return number >= 1 && number <= 12
	case BetDozen2:
		return number >= 13 && number <= 24
	case BetDozen3:
		return number >= 25 && number <= 36
	}
	return false
}

func (b *Bet) Payout() float64 {
	return b.Amount * Payouts[b.Type]
}

func Spin(rng *rand.Rand) int {
	return rng.Intn(37)
}
