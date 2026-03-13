package slots

import "fmt"

func SymbolDisplay(s Symbol) string {
	switch s {
	case SEVEN:
		return "7"
	case BAR:
		return "BAR"
	case BELL:
		return "BELL"
	case CHERRY:
		return "CHERRY"
	case LEMON:
		return "LEMON"
	case GRAPE:
		return "GRAPE"
	default:
		return "?"
	}
}

func padCenter(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	left := (width - len(s)) / 2
	right := width - len(s) - left
	return fmt.Sprintf("%s%s%s", spaces(left), s, spaces(right))
}

func spaces(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += " "
	}
	return result
}

// Inner content width breakdown:
//   ║  │ 7chars │  │ 7chars │  │ 7chars │   ║
//   1 + 2 + 9 + 2 + 9 + 2 + 9 + 3 + 1 = 38 total
//   => 36 chars between the two ║, so outer border = 36x═

func RenderReels(symbols Reels, reelStates [3]bool) string {
	s := "╔════════════════════════════════════╗\n"
	s += "║  ┌───────┐  ┌───────┐  ┌───────┐   ║\n"
	for row := 0; row < 3; row++ {
		s += "║  "
		for r := 0; r < 3; r++ {
			s += "│"
			if reelStates[r] {
				s += " ????? "
			} else {
				sym := symbols[r]
				display := SymbolDisplay(sym)
				if row == 1 {
					s += padCenter(display, 7)
				} else {
					s += "       "
				}
			}
			if r < 2 {
				s += "│  "
			} else {
				s += "│"
			}
		}
		s += "   ║\n"
	}
	s += "║  └───────┘  └───────┘  └───────┘   ║\n"
	s += "╚════════════════════════════════════╝"
	return s
}

func RenderSpinning() string {
	s := "╔════════════════════════════════════╗\n"
	s += "║  ┌───────┐  ┌───────┐  ┌───────┐   ║\n"
	for row := 0; row < 3; row++ {
		s += "║  "
		for r := 0; r < 3; r++ {
			s += "│"
			if row == 1 {
				s += " ????? "
			} else {
				s += "       "
			}
			if r < 2 {
				s += "│  "
			} else {
				s += "│"
			}
		}
		s += "   ║\n"
	}
	s += "║  └───────┘  └───────┘  └───────┘   ║\n"
	s += "╚════════════════════════════════════╝"
	return s
}

func RenderWinBanner(amount float64) string {
	return fmt.Sprintf(`
    ╔═══════════════════════╗
    ║    YOU WON!           ║
    ║    $%.2f              ║
    ╚═══════════════════════╝
`, amount)
}

func RenderLoseBanner() string {
	return `
    ╔═══════════════════════╗
    ║    NEVER LUCKY..      ║
    ║    TRY AGAIN!         ║
    ╚═══════════════════════╝
`
}
