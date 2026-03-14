package roulette

import "fmt"

func RenderWheel(result int) string {
	s := "┌─────────────────────────────────┐\n"
	s += "│         ROULETTE WHEEL          │\n"
	s += "├─────────────────────────────────┤\n"

	if result >= 0 {
		s += fmt.Sprintf("│           ║ %2d ║                │\n", result)
	} else {
		s += "│             ║    ║              │\n"
	}

	s += "└─────────────────────────────────┘"
	return s
}

func RenderBettingTable(bets map[int]float64, section Section, row, col, index int) string {
	selected := 0
	switch section {
	case SectionMain:
		rows := [][]int{
			{1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 34},
			{2, 5, 8, 11, 14, 17, 20, 23, 26, 29, 32, 35},
			{3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36},
		}
		selected = rows[row][col]
	case SectionDozens:
		selected = 101 + index
	case SectionLowHigh:
		// 0, 1-18, 19-36
		lowHighIds := []int{0, 104, 105}
		selected = lowHighIds[index]
	case SectionOutside:
		// EVEN, RED, BLACK, ODD
		outsideIds := []int{106, 108, 109, 107}
		selected = outsideIds[index]
	}

	rows := [][]int{
		{1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 34},
		{2, 5, 8, 11, 14, 17, 20, 23, 26, 29, 32, 35},
		{3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36},
	}

	cell := func(n int) string {
		s := fmt.Sprintf("%2d", n)
		if selected == n {
			return "[" + s + "]"
		}
		if bets[n] > 0 {
			return "*" + s + "*"
		}
		return " " + s + " "
	}

	dozenCell := func(dozenIdx int) string {
		dozenNum := 101 + dozenIdx
		d := []string{"1-12", "13-24", "25-36"}[dozenIdx]
		if selected == dozenNum {
			return "[" + d + "]"
		}
		if bets[dozenNum] > 0 {
			return "*" + d + "*"
		}
		return " " + d + " "
	}

	// For Low/High section (1-18, 19-36)
	lowHighCell := func(i int) string {
		lowHighIds := []int{104, 105}
		lowHighNames := []string{"1-18", "19-36"}
		id := lowHighIds[i]
		name := lowHighNames[i]
		if selected == id {
			return "[" + name + "]"
		}
		if bets[id] > 0 {
			return "*" + name + "*"
		}
		return " " + name + " "
	}

	// For Outside section (EVEN, RED, BLACK, ODD)
	outsideCell := func(i int) string {
		outsideIds := []int{106, 108, 109, 107}
		outs := []string{"EVEN", "RED", "BLACK", "ODD"}
		id := outsideIds[i]
		o := outs[i]
		if selected == id {
			return "[" + o + "]"
		}
		if bets[id] > 0 {
			return "*" + o + "*"
		}
		return " " + o + " "
	}

	// Build table
	// Each cell is 4 chars wide (including brackets/asterisks)
	// 12 cells = 48 chars, plus borders
	s := "\n"
	s += "┌────────────────────────────────────────────────┐\n"

	// Row 0
	s += "│"
	for _, n := range rows[0] {
		s += cell(n)
	}
	s += "│\n"

	// Row 1
	s += "│"
	for _, n := range rows[1] {
		s += cell(n)
	}
	s += "│\n"

	// Row 2
	s += "│"
	for _, n := range rows[2] {
		s += cell(n)
	}
	s += "│\n"

	// Dozens row
	s += "├────────────────────────────────────────────────┤\n"
	s += "│ " + dozenCell(0) + "  " + dozenCell(1) + "  " + dozenCell(2) + "                       │\n"

	// 0, 1-18 and 19-36 row
	s += "├────────────────────────────────────────────────┤\n"
	var zeroCell string
	if selected == 0 {
		zeroCell = "[0]"
	} else if bets[0] > 0 {
		zeroCell = "*0*"
	} else {
		zeroCell = " 0 "
	}
	s += "│ " + zeroCell + "  " + lowHighCell(0) + "  " + lowHighCell(1) + "                                              │\n"
	
	// Outside bets (EVEN, RED, BLACK, ODD)
	s += "├────────────────────────────────────────────────┤\n"
	s += "│ "
	for i := 0; i < 4; i++ {
		s += outsideCell(i)
		if i < 3 {
			s += "  "
		}
	}
	s += "                                   │\n"

	s += "└────────────────────────────────────────────────┘"
	
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
