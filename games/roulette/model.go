package roulette

import (
	"fmt"
	"math/rand"
	"time"

	"cli-casino/ui"

	tea "github.com/charmbracelet/bubbletea"
)

type Section int

const (
	SectionMain Section = iota
	SectionDozens
	SectionLowHigh
	SectionOutside
)

type Model struct {
	Wallet    ui.WalletBackend
	Bet       float64
	Bets      map[int]float64
	Spinning  bool
	Result    int
	WheelTick int
	Section   Section
	Row       int
	Col       int
	Index     int
	rng       *rand.Rand
}

func NewModel(wallet ui.WalletBackend) *Model {
	return &Model{
		Wallet:   wallet,
		Bet:      10,
		Bets:     make(map[int]float64),
		Spinning: false,
		Result:   -1,
		Section:  SectionMain,
		Row:      0,
		Col:      0,
		Index:    0,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case SpinTick:
		return m.handleSpinTick()
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Spinning {
		return m, nil
	}

	switch msg.String() {
	case "tab":
		m.nextSection()
	case "shift+tab", "btab":
		m.prevSection()
	case "up":
		m.moveUp()
	case "down":
		m.moveDown()
	case "left":
		m.moveLeft()
	case "right":
		m.moveRight()
	case " ":
		return m.placeBet()
	case "c":
		m.clearBets()
	case "s":
		return m.spin()
	case "q", "esc":
		m.Wallet.Save()
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) nextSection() {
	m.Section = (m.Section + 1) % 4
	m.Row = 0
	m.Col = 0
	m.Index = 0
}

func (m *Model) prevSection() {
	m.Section = (m.Section + 3) % 4
	m.Row = 0
	m.Col = 0
	m.Index = 0
}

func (m *Model) moveUp() {
	switch m.Section {
	case SectionMain:
		if m.Row > 0 {
			m.Row--
		}
	case SectionDozens:
		m.Section = SectionMain
		m.Row = 2
	case SectionLowHigh:
		m.Section = SectionDozens
		m.Index = 0
	case SectionOutside:
		m.Section = SectionLowHigh
		m.Index = 0
	}
}

func (m *Model) moveDown() {
	switch m.Section {
	case SectionMain:
		if m.Row < 2 {
			m.Row++
		} else {
			m.Section = SectionDozens
			m.Index = 0
		}
	case SectionDozens:
		m.Section = SectionLowHigh
		m.Index = 0
	case SectionLowHigh:
		m.Section = SectionOutside
		m.Index = 0
	case SectionOutside:
		return
	}
}

func (m *Model) moveLeft() {
	switch m.Section {
	case SectionMain:
		if m.Col > 0 {
			m.Col--
		}
	case SectionDozens:
		if m.Index > 0 {
			m.Index--
		}
	case SectionLowHigh:
		if m.Index > 0 {
			m.Index--
		}
	case SectionOutside:
		if m.Index > 0 {
			m.Index--
		}
	}
}

func (m *Model) moveRight() {
	switch m.Section {
	case SectionMain:
		if m.Col < 11 {
			m.Col++
		}
	case SectionDozens:
		if m.Index < 2 {
			m.Index++
		}
	case SectionLowHigh:
		if m.Index < 2 {
			m.Index++
		}
	case SectionOutside:
		if m.Index < 3 {
			m.Index++
		}
	}
}

func (m *Model) getSelected() int {
	switch m.Section {
	case SectionMain:
		rows := [][]int{
			{1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 34},
			{2, 5, 8, 11, 14, 17, 20, 23, 26, 29, 32, 35},
			{3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36},
		}
		return rows[m.Row][m.Col]
	case SectionDozens:
		return 101 + m.Index
	case SectionLowHigh:
		// 0, 1-18, 19-36
		lowHighIds := []int{0, 104, 105}
		return lowHighIds[m.Index]
	case SectionOutside:
		// EVEN, RED, BLACK, ODD
		outsideIds := []int{106, 108, 109, 107}
		return outsideIds[m.Index]
	}
	return 1
}

func (m *Model) getSelectedName() string {
	switch m.Section {
	case SectionMain:
		return fmt.Sprintf("#%d", m.getSelected())
	case SectionDozens:
		dozens := []string{"1-12", "13-24", "25-36"}
		return dozens[m.Index]
	case SectionLowHigh:
		lowHighNames := []string{"#0", "1-18", "19-36"}
		return lowHighNames[m.Index]
	case SectionOutside:
		outsideNames := []string{"EVEN", "RED", "BLACK", "ODD"}
		return outsideNames[m.Index]
	}
	return "#1"
}

func (m *Model) placeBet() (tea.Model, tea.Cmd) {
	totalBet := m.totalBets()
	if !m.Wallet.CanAfford(m.Bet + totalBet) {
		return m, nil
	}

	m.Bets[m.getSelected()] += m.Bet
	return m, nil
}

func (m *Model) clearBets() {
	m.Bets = make(map[int]float64)
}

func (m *Model) spin() (tea.Model, tea.Cmd) {
	totalBet := m.totalBets()
	if totalBet <= 0 {
		return m, nil
	}

	if !m.Wallet.CanAfford(totalBet) {
		return m, nil
	}

	for _, bet := range m.Bets {
		m.Wallet.Bet(bet)
	}

	m.Spinning = true
	m.Result = -1
	m.WheelTick = 0

	return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinTick{}
	})
}

type SpinTick struct{}

func (m *Model) handleSpinTick() (tea.Model, tea.Cmd) {
	m.WheelTick++

	if m.WheelTick >= 15 {
		m.Result = Spin(m.rng)
		m.Spinning = false
		m.calculateWinnings()
		m.WheelTick = 0
		return m, nil
	}

	return m, tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinTick{}
	})
}

func (m *Model) calculateWinnings() {
	if m.Result < 0 {
		return
	}

	totalWinnings := 0.0

	for number, bet := range m.Bets {
		// Main numbers
		if number >= 0 && number <= 36 && m.Result == number {
			totalWinnings += bet * 35
		}
		// Dozens
		if number >= 101 && number <= 103 {
			dozenNum := number - 101
			start := dozenNum*12 + 1
			end := start + 11
			if m.Result >= start && m.Result <= end {
				totalWinnings += bet * 2
			}
		}
		// Outside
		if number == 104 && m.Result >= 1 && m.Result <= 18 {
			totalWinnings += bet
		}
		if number == 105 && m.Result >= 19 && m.Result <= 36 {
			totalWinnings += bet
		}
		if number == 106 && m.Result%2 == 0 && m.Result != 0 {
			totalWinnings += bet
		}
		if number == 107 && m.Result%2 == 1 {
			totalWinnings += bet
		}
		if number == 108 && isRed(m.Result) {
			totalWinnings += bet
		}
		if number == 109 && isBlack(m.Result) {
			totalWinnings += bet
		}
	}

	if totalWinnings > 0 {
		m.Wallet.Win(totalWinnings)
	}
}

func isRed(n int) bool {
	for _, r := range RedNumbers {
		if n == r {
			return true
		}
	}
	return false
}

func isBlack(n int) bool {
	for _, b := range BlackNumbers {
		if n == b {
			return true
		}
	}
	return false
}

func (m *Model) totalBets() float64 {
	total := 0.0
	for _, bet := range m.Bets {
		total += bet
	}
	return total
}

func (m *Model) View() string {
	s := ui.HeaderStyle.Render("🎡 ROULETTE 🎡") + "\n\n"
	s += m.Wallet.Render() + "\n\n"

	s += RenderWheel(m.Result) + "\n\n"

	selName := m.getSelectedName()

	// Section indicator
	sectionName := "MAIN"
	switch m.Section {
	case SectionDozens:
		sectionName = "DOZENS"
	case SectionLowHigh:
		sectionName = "LOW/HIGH"
	case SectionOutside:
		sectionName = "OUTSIDE"
	}
	s += fmt.Sprintf("[%s] ", sectionName)

	s += RenderBettingTable(m.Bets, m.Section, m.Row, m.Col, m.Index) + "\n\n"

	s += fmt.Sprintf("Selection: %s  Bet: $%.2f  Total: $%.2f\n", selName, m.Bet, m.totalBets())
	s += RenderBets(m.Bets) + "\n\n"

	if m.Result >= 0 && !m.Spinning {
		totalWinnings := m.calculateViewWinnings()
		if totalWinnings > 0 {
			s += ui.WinStyle.Render(RenderResult(m.Result, totalWinnings)) + "\n\n"
		} else {
			s += ui.LoseStyle.Render(RenderResult(m.Result, 0)) + "\n\n"
		}
	}

	s += ui.SubtleStyle.Render("[TAB] sections  [arrows] move  [space] bet  [c] clear  [s] spin  [q] quit")

	return s
}

func (m *Model) calculateViewWinnings() float64 {
	if m.Result < 0 {
		return 0
	}

	totalWinnings := 0.0

	for number, bet := range m.Bets {
		// Main numbers
		if number >= 0 && number <= 36 && m.Result == number {
			totalWinnings += bet * 35
		}
		// Dozens
		if number >= 101 && number <= 103 {
			dozenNum := number - 101
			start := dozenNum*12 + 1
			end := start + 11
			if m.Result >= start && m.Result <= end {
				totalWinnings += bet * 2
			}
		}
		// Outside
		if number == 104 && m.Result >= 1 && m.Result <= 18 {
			totalWinnings += bet
		}
		if number == 105 && m.Result >= 19 && m.Result <= 36 {
			totalWinnings += bet
		}
		if number == 106 && m.Result%2 == 0 && m.Result != 0 {
			totalWinnings += bet
		}
		if number == 107 && m.Result%2 == 1 {
			totalWinnings += bet
		}
		if number == 108 && isRed(m.Result) {
			totalWinnings += bet
		}
		if number == 109 && isBlack(m.Result) {
			totalWinnings += bet
		}
	}

	return totalWinnings
}
