package roulette

import "fmt"

func RenderWheel(result int) string {
	s := "┌─────────────────────────────────┐\n"
	s += "│         ROULETTE WHEEL          │\n"
	s += "├─────────────────────────────────┤\n"

	if result >= 0 {
		s += fmt.Sprintf("│           ║ %2d ║              │\n", result)
	} else {
		s += "│           ║    ║              │\n"
	}

	s += "└─────────────────────────────────┘"
	return s
}

func RenderBettingTable(bets map[int]float64, section Section, index int) string {
	selected := 0
	switch section {
	case SectionMain:
		selected = index
	case SectionDozens:
		selected = 101 + index
	case SectionOutside:
		outsideIds := []int{104, 106, 108, 109, 107, 105}
		selected = outsideIds[index]
	}

	// Main grid: 3 rows x 12 columns
	// Row 0: 1,4,7,10,13,16,19,22,25,28,31,34
	// Row 1: 2,5,8,11,14,17,20,23,26,29,32,35
	// Row 2: 3,6,9,12,15,18,21,24,27,30,33,36
	rows := [][]int{
		{1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 34},
		{2, 5, 8, 11, 14, 17, 20, 23, 26, 29, 32, 35},
		{3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36},
	}

	s := "┌───────┬──────────────────────────────────────────┬───────┐\n"
	s += "│       │  1   4   7  10  13  16  19  22  25  28  31  34│ 1-12  │\n"
	s += "├───────┼──────────────────────────────────────────┼───────┤\n"

	// Draw the 3 rows of numbers
	for rowIdx, row := range rows {
		if rowIdx == 0 {
			s += "│  1st  │"
		} else if rowIdx == 1 {
			s += "│  2nd  │"
		} else {
			s += "│  3rd  │"
		}

		for _, n := range row {
			if selected == n {
				s += fmt.Sprintf("[%2d]", n)
			} else if bets[n] > 0 {
				s += fmt.Sprintf("*%2d*", n)
			} else {
				s += fmt.Sprintf(" %2d ", n)
			}
		}

		// Dozens column
		dozenNum := 101 + rowIdx
		if selected == dozenNum {
			s += fmt.Sprintf("│[%d]│\n", dozenNum)
		} else if bets[dozenNum] > 0 {
			s += fmt.Sprintf("│*%d*│\n", dozenNum)
		} else {
			s += fmt.Sprintf("│ %d  │\n", dozenNum)
		}
	}

	// Zero + Dozens row
	s += "├───────┼──────────────────────────────────────────┼───────┤\n"
	if selected == 0 {
		s += "│[ 0 ] │"
	} else if bets[0] > 0 {
		s += "│* 0 * │"
	} else {
		s += "│  0   │"
	}
	s += "              (1-12) (13-24) (25-36)              │       │\n"

	// Dozens selections (when in dozens section)
	if section == SectionDozens {
		dozens := []string{"1-12", "13-24", "25-36"}
		s += "├───────┼──────────────────────────────────────────┤\n"
		s += "│DOZENS │"
		for i, d := range dozens {
			if index == i {
				s += fmt.Sprintf(" [%s] ", d)
			} else {
				s += fmt.Sprintf("  %s  ", d)
			}
		}
		s += "                                        │       │\n"
	} else {
		s += "├───────┼──────────────────────────────────────────┤\n"
		s += "│       │"
	}

	// Outside bets row
	if section == SectionOutside {
		outs := []string{"1-18", "EVEN", "RED", "BLACK", "ODD", "19-36"}
		for i, o := range outs {
			if index == i {
				s += fmt.Sprintf("[%s]", o)
			} else {
				s += fmt.Sprintf(" %s ", o)
			}
		}
		s += "│\n"
	} else {
		s += " 1-18   EVEN   RED   BLACK   ODD   19-36   │\n"
	}

	s += "└───────┴──────────────────────────────────────────┘"
	return s
}

func RenderBets(bets map[int]float64) string {
	if len(bets) == 0 {
		return "No bets placed"
	}
	total := 0.0
	for _, amt := range bets {
		total += amt
	}
	return fmt.Sprintf("Total bet: $%.2f", total)
}

func RenderResult(result int, winnings float64) string {
	if result < 0 {
		return ""
	}
	if winnings > 0 {
		return fmt.Sprintf("WINNER! #%d — You won $%.2f", result, winnings)
	}
	return fmt.Sprintf("No winner. Ball landed on %d", result)
}
