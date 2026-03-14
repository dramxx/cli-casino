package blackjack

import (
	"fmt"
	"strings"
)

func RenderCard(card Card, faceDown bool) string {
	if faceDown {
		return "┌─────────┐\n│░░░░░░░░░│\n│░░░░░░░░░│\n│░░░░░░░░░│\n│░░░░░░░░░│\n└─────────┘"
	}

	lines := make([]string, 6)
	lines[0] = "┌─────────┐"
	lines[1] = fmt.Sprintf("│ %-2s      │", card.Rank)
	lines[2] = "│         │"
	lines[3] = fmt.Sprintf("│    %s    │", card.Suit)
	lines[4] = fmt.Sprintf("│      %2s │", card.Rank)
	lines[5] = "└─────────┘"

	return strings.Join(lines, "\n")
}

func RenderHand(cards []Card, hideFirst bool) string {
	if len(cards) == 0 {
		return ""
	}

	cardLines := make([][]string, len(cards))
	for i, card := range cards {
		faceDown := hideFirst && i == 0
		cardLines[i] = strings.Split(RenderCard(card, faceDown), "\n")
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

	return strings.Join(result, "\n")
}

func RenderHandValue(hand *Hand, hideFirst bool) string {
	if hideFirst {
		return "?"
	}
	value := hand.Value()
	if hand.IsSoft() {
		return fmt.Sprintf("Soft %d", value)
	}
	return fmt.Sprintf("%d", value)
}

