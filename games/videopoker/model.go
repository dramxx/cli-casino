package videopoker

import (
	"fmt"
	"math/rand"
	"time"

	"cli-casino/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type GameState int

const (
	StateBetting GameState = iota
	StateHolding
	StateDrawing
	StateResult
)

type Model struct {
	Wallet   ui.WalletBackend
	Bet      float64
	MinBet   float64
	MaxBet   float64
	Deck     *Deck
	Hand     []Card
	Held     []bool
	Selected int
	State    GameState
	Result   HandRank
	Message  string
	rng      *rand.Rand
}

func NewModel(wallet ui.WalletBackend) *Model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Model{
		Wallet:   wallet,
		Bet:      5,
		MinBet:   5,
		MaxBet:   25,
		State:    StateBetting,
		Message:  "Place your bet and press ENTER to deal",
		Selected: 0,
		rng:      rng,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.State {
	case StateBetting:
		return m.handleBettingKey(msg)
	case StateHolding:
		return m.handleHoldingKey(msg)
	case StateResult:
		return m.handleResultKey(msg)
	}

	if msg.String() == "q" || msg.String() == "esc" {
		m.Wallet.Save()
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) handleBettingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		m.Bet += 5
		if m.Bet > m.MaxBet {
			m.Bet = m.MaxBet
		}
		if m.Bet > m.Wallet.GetBalance() {
			m.Bet = m.Wallet.GetBalance()
		}
	case "down":
		m.Bet -= 5
		if m.Bet < m.MinBet {
			m.Bet = m.MinBet
		}
	case "enter":
		return m.deal()
	case "q", "esc":
		m.Wallet.Save()
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleHoldingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left":
		m.Selected--
		if m.Selected < 0 {
			m.Selected = 4
		}
	case "right":
		m.Selected = (m.Selected + 1) % 5
	case "1":
		m.Selected = 0
	case "2":
		m.Selected = 1
	case "3":
		m.Selected = 2
	case "4":
		m.Selected = 3
	case "5":
		m.Selected = 4
	case " ":
		m.Held[m.Selected] = !m.Held[m.Selected]
	case "enter":
		return m.draw()
	case "q", "esc":
		m.Wallet.Save()
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleResultKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		return m.newRound()
	case "q", "esc":
		m.Wallet.Save()
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) deal() (tea.Model, tea.Cmd) {
	if !m.Wallet.CanAfford(m.Bet) {
		m.Message = "Insufficient funds!"
		return m, nil
	}

	m.Wallet.Bet(m.Bet)
	m.Deck = NewDeck(m.rng)
	m.Hand = make([]Card, 5)
	m.Held = make([]bool, 5)
	m.Selected = 0

	for i := 0; i < 5; i++ {
		m.Hand[i] = m.Deck.Draw()
	}

	m.State = StateHolding
	m.Message = "Select cards to hold, then press ENTER to draw"
	return m, nil
}

func (m *Model) draw() (tea.Model, tea.Cmd) {
	for i := 0; i < 5; i++ {
		if !m.Held[i] {
			m.Hand[i] = m.Deck.Draw()
		}
	}

	m.State = StateResult
	m.Result = EvaluateHand(m.Hand)

	payout := m.Result.Payout(m.Bet)
	if payout > 0 {
		m.Wallet.Win(payout)
		m.Message = fmt.Sprintf("%s! You win $%.0f", m.Result.String(), payout)
	} else {
		m.Message = fmt.Sprintf("%s - No payout", m.Result.String())
	}

	return m, nil
}

func (m *Model) newRound() (tea.Model, tea.Cmd) {
	m.Bet = 5
	if m.Bet > m.Wallet.GetBalance() {
		m.Bet = m.Wallet.GetBalance()
	}
	m.State = StateBetting
	m.Result = HighCard
	m.Message = "Place your bet and press ENTER to deal"
	m.Hand = nil
	m.Held = nil
	m.Selected = 0
	return m, nil
}

func (m *Model) View() string {
	s := ui.HeaderStyle.Render("🎲 VIDEO POKER 🎲") + "\n\n"
	s += m.Wallet.Render() + "\n\n"

	if m.State != StateBetting {
		s += RenderHand(m.Hand, m.Held) + "\n\n"

		if m.State == StateHolding {
			s += "Selected: Card " + fmt.Sprintf("%d", m.Selected+1) + "\n\n"
		}
	}

	if m.State == StateResult {
		s += RenderPaytable(m.Bet, m.Result) + "\n\n"
	} else {
		s += RenderPaytable(m.Bet, HighCard) + "\n\n"
	}

	s += fmt.Sprintf("Bet: $%.0f\n", m.Bet)
	s += m.Message + "\n\n"

	switch m.State {
	case StateBetting:
		s += ui.SubtleStyle.Render("[↑↓] adjust bet  [enter] deal  [q] quit")
	case StateHolding:
		s += ui.SubtleStyle.Render("[←→] or [1-5] select  [space] toggle hold  [enter] draw  [q] quit")
	case StateResult:
		s += ui.SubtleStyle.Render("[enter] new round  [q] quit")
	}

	return s
}
