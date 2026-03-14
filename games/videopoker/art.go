package videopoker

import (
	"fmt"
	"strings"
)

func RenderCard(card Card) string {
	lines := make([]string, 6)
	lines[0] = "┌─────────┐"
	lines[1] = fmt.Sprintf("│ %-2s      │", card.Rank.String())
	lines[2] = "│         │"
	lines[3] = fmt.Sprintf("│    %s    │", card.Suit)
	lines[4] = fmt.Sprintf("│      %2s │", card.Rank.String())
	lines[5] = "└─────────┘"

	return strings.Join(lines, "\n")
}

func RenderHand(cards []Card, held []bool) string {
	if len(cards) == 0 {
		return ""
	}

	cardLines := make([][]string, len(cards))
	for i, card := range cards {
		cardLines[i] = strings.Split(RenderCard(card), "\n")
	}

	result := make([]string, 6)
	for row := 0; row < 6; row++ {
		for col := 0; col < len(cards); col++ {
			result[row] += cardLines[col][row]
			if col < len(cards)-1 {
				result[row] += " "
			}
		}
	}

	holdLine := ""
	for i := 0; i < len(cards); i++ {
		if held[i] {
			holdLine += "  [HOLD]   "
		} else {
			holdLine += "           "
		}
		if i < len(cards)-1 {
			holdLine += " "
		}
	}

	result = append(result, holdLine)
	return strings.Join(result, "\n")
}

func RenderPaytable(bet float64, currentRank HandRank) string {
	ranks := []HandRank{
		RoyalFlush,
		StraightFlush,
		FourOfAKind,
		FullHouse,
		Flush,
		Straight,
		ThreeOfAKind,
		TwoPair,
		JacksOrBetter,
	}

	s := "┌────────────────────────┬──────────┐\n"
	s += "│       PAYTABLE         │  PAYOUT  │\n"
	s += "├────────────────────────┼──────────┤\n"

	for _, rank := range ranks {
		payout := rank.Payout(bet)
		marker := " "
		if rank == currentRank {
			marker = "►"
		}
		s += fmt.Sprintf("│%s %-21s │ $%-7.0f │\n", marker, rank.String(), payout)
	}

	s += "└────────────────────────┴──────────┘"
	return s
}
