package slots

import (
	"fmt"
	"math/rand"
	"time"

	"cli-casino/ui"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Wallet     ui.WalletBackend
	Bet        float64
	Spinning   bool
	Reels      Reels
	ReelStates [3]bool
	Result     *Result
	Tick       int
	rng        *rand.Rand
}

func NewModel(wallet ui.WalletBackend) *Model {
	return &Model{
		Wallet:     wallet,
		Bet:        10,
		Spinning:   false,
		Reels:      Reels{},
		ReelStates: [3]bool{false, false, false},
		Result:     nil,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
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

	keyStr := msg.String()
	isUp := msg.Type == tea.KeyUp || keyStr == "up"
	isDown := msg.Type == tea.KeyDown || keyStr == "down"

	switch {
	case msg.Type == tea.KeySpace || keyStr == " ":
		return m.spin()
	case isUp:
		if m.Wallet.GetBalance() >= 10 && m.Bet < 100 && m.Bet < m.Wallet.GetBalance() {
			m.Bet += 10
		}
	case isDown:
		if m.Bet > 10 && m.Wallet.GetBalance() >= 10 {
			m.Bet -= 10
		}
	case msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc || keyStr == "q" || keyStr == "esc":
		m.Wallet.Save()
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) spin() (tea.Model, tea.Cmd) {
	if !m.Wallet.CanAfford(m.Bet) {
		return m, nil
	}

	m.Wallet.Bet(m.Bet)
	m.Spinning = true
	m.ReelStates = [3]bool{true, true, true}
	m.Result = nil

	return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinTick{Tick: 0}
	})
}

type SpinTick struct {
	Tick int
}

func (m *Model) handleSpinTick() (tea.Model, tea.Cmd) {
	m.Tick++

	switch m.Tick {
	case 8:
		m.ReelStates[0] = false
		m.Reels[0] = SpinReels(m.rng)[0]
	case 12:
		m.ReelStates[1] = false
		m.Reels[1] = SpinReels(m.rng)[1]
	case 16:
		m.ReelStates[2] = false
		m.Reels[2] = SpinReels(m.rng)[2]
		m.Result = CalculatePayout(m.Reels, m.Bet)
		if m.Result.Win {
			m.Wallet.Win(m.Result.AmountWon)
		}
		m.Spinning = false
		m.Tick = 0
		return m, nil
	}

	if m.Spinning {
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return SpinTick{Tick: m.Tick}
		})
	}

	return m, nil
}

func (m *Model) View() string {
	s := ui.HeaderStyle.Render("🎰 SLOTS 🎰") + "\n\n"
	s += m.Wallet.Render() + "\n\n"

	if m.Spinning || !allReelsStopped(m.ReelStates) {
		s += RenderSpinning() + "\n\n"
	} else {
		s += RenderReels(m.Reels, m.ReelStates) + "\n\n"

		if m.Result != nil {
			if m.Result.Win {
				s += ui.WinStyle.Render(RenderWinBanner(m.Result.AmountWon)) + "\n\n"
			} else {
				s += ui.LoseStyle.Render(RenderLoseBanner()) + "\n\n"
			}
		}
	}

	s += fmt.Sprintf("Bet: $%.2f  [↑↓ to adjust]\n", m.Bet)
	s += "\n" + ui.SubtleStyle.Render("[SPACE] spin  [q] back to menu")

	return s
}

func allReelsStopped(states [3]bool) bool {
	return !states[0] && !states[1] && !states[2]
}
