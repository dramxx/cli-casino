package blackjack

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
	StateDealing
	StatePlayerTurn
	StateDealerTurn
	StateResolution
)

type GameResult int

const (
	ResultNone GameResult = iota
	ResultPlayerBlackjack
	ResultPlayerWin
	ResultDealerWin
	ResultDealerBust
	ResultPlayerBust
	ResultPush
)

type Model struct {
	Wallet     ui.WalletBackend
	Bet        float64
	MinBet     float64
	MaxBet     float64
	Deck       *Deck
	PlayerHand Hand
	DealerHand Hand
	State      GameState
	Result     GameResult
	Message    string
	CanDouble  bool
	CanSplit   bool
	rng        *rand.Rand
}

func NewModel(wallet ui.WalletBackend) *Model {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Model{
		Wallet:  wallet,
		Bet:     10,
		MinBet:  10,
		MaxBet:  100,
		Deck:    NewDeck(rng, 6),
		State:   StateBetting,
		Result:  ResultNone,
		Message: "Place your bet and press ENTER to deal",
		rng:     rng,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case DealerDrawMsg:
		return m.handleDealerDraw()
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.State {
	case StateBetting:
		return m.handleBettingKey(msg)
	case StatePlayerTurn:
		return m.handlePlayerTurnKey(msg)
	case StateResolution:
		return m.handleResolutionKey(msg)
	}

	if msg.String() == "q" || msg.String() == "esc" {
		m.Wallet.Save()
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) handleBettingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()
	isUp := msg.Type == tea.KeyUp || keyStr == "up"
	isDown := msg.Type == tea.KeyDown || keyStr == "down"

	switch {
	case isUp:
		if m.Wallet.GetBalance() >= m.MinBet {
			m.Bet += 10
			if m.Bet > m.MaxBet {
				m.Bet = m.MaxBet
			}
			if m.Bet > m.Wallet.GetBalance() {
				m.Bet = m.Wallet.GetBalance()
			}
		}
	case isDown:
		if m.Bet > m.MinBet && m.Wallet.GetBalance() >= m.MinBet {
			m.Bet -= 10
		}
	case msg.Type == tea.KeyEnter || keyStr == "enter":
		return m.deal()
	case msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc || keyStr == "q" || keyStr == "esc":
		m.Wallet.Save()
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handlePlayerTurnKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "h":
		return m.hit()
	case "s":
		return m.stand()
	case "d":
		if m.CanDouble {
			return m.double()
		}
	case "q", "esc":
		m.Wallet.Save()
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleResolutionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

	if m.Deck.NeedsReshuffle() {
		m.Deck = NewDeck(m.rng, 6)
		m.Message = "Shuffling new deck..."
	}

	m.Wallet.Bet(m.Bet)
	m.PlayerHand = Hand{}
	m.DealerHand = Hand{}
	m.State = StateDealing

	m.PlayerHand.Add(m.Deck.Draw())
	m.DealerHand.Add(m.Deck.Draw())
	m.PlayerHand.Add(m.Deck.Draw())
	m.DealerHand.Add(m.Deck.Draw())

	m.CanDouble = len(m.PlayerHand.Cards) == 2 && m.Wallet.CanAfford(m.Bet)
	m.CanSplit = m.PlayerHand.CanSplit() && m.Wallet.CanAfford(m.Bet)

	if m.PlayerHand.IsBlackjack() {
		if m.DealerHand.IsBlackjack() {
			m.Result = ResultPush
			m.Wallet.Win(m.Bet)
			m.Message = "Both blackjack - Push!"
		} else {
			m.Result = ResultPlayerBlackjack
			payout := m.Bet * 2.5
			m.Wallet.Win(payout)
			m.Message = fmt.Sprintf("Blackjack! You win $%.2f", payout-m.Bet)
		}
		m.State = StateResolution
		return m, nil
	}

	m.State = StatePlayerTurn
	m.Message = "Your turn: (h)it, (s)tand"
	if m.CanDouble {
		m.Message += ", (d)ouble"
	}

	return m, nil
}

func (m *Model) hit() (tea.Model, tea.Cmd) {
	m.PlayerHand.Add(m.Deck.Draw())
	m.CanDouble = false

	if m.PlayerHand.IsBust() {
		m.Result = ResultPlayerBust
		m.State = StateResolution
		m.Message = fmt.Sprintf("Bust! You lose $%.2f", m.Bet)
		return m, nil
	}

	m.Message = "Your turn: (h)it, (s)tand"
	return m, nil
}

func (m *Model) stand() (tea.Model, tea.Cmd) {
	m.State = StateDealerTurn
	m.Message = "Dealer's turn..."
	return m, m.dealerPlay()
}

func (m *Model) double() (tea.Model, tea.Cmd) {
	if !m.Wallet.CanAfford(m.Bet) {
		m.Message = "Insufficient funds to double!"
		return m, nil
	}

	m.Wallet.Bet(m.Bet)
	m.Bet *= 2
	m.PlayerHand.Add(m.Deck.Draw())
	m.CanDouble = false

	if m.PlayerHand.IsBust() {
		m.Result = ResultPlayerBust
		m.State = StateResolution
		m.Message = fmt.Sprintf("Bust! You lose $%.2f", m.Bet)
		return m, nil
	}

	return m.stand()
}

type DealerDrawMsg struct{}

func (m *Model) dealerPlay() tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return DealerDrawMsg{}
	})
}

func (m *Model) handleDealerDraw() (tea.Model, tea.Cmd) {
	dealerValue := m.DealerHand.Value()

	if dealerValue < 17 || (dealerValue == 17 && m.DealerHand.IsSoft()) {
		m.DealerHand.Add(m.Deck.Draw())
		m.Message = fmt.Sprintf("Dealer draws... (%d)", m.DealerHand.Value())
		return m, m.dealerPlay()
	}

	m.State = StateResolution
	m.resolveGame()
	return m, nil
}

func (m *Model) resolveGame() {
	playerValue := m.PlayerHand.Value()
	dealerValue := m.DealerHand.Value()

	if m.DealerHand.IsBust() {
		m.Result = ResultDealerBust
		payout := m.Bet * 2
		m.Wallet.Win(payout)
		m.Message = fmt.Sprintf("Dealer bust! You win $%.2f", payout-m.Bet)
	} else if playerValue > dealerValue {
		m.Result = ResultPlayerWin
		payout := m.Bet * 2
		m.Wallet.Win(payout)
		m.Message = fmt.Sprintf("You win! You win $%.2f", payout-m.Bet)
	} else if playerValue < dealerValue {
		m.Result = ResultDealerWin
		m.Message = fmt.Sprintf("Dealer wins. You lose $%.2f", m.Bet)
	} else {
		m.Result = ResultPush
		m.Wallet.Win(m.Bet)
		m.Message = "Push - Bet returned"
	}
}

func (m *Model) newRound() (tea.Model, tea.Cmd) {
	m.Bet = 10
	if m.Bet > m.Wallet.GetBalance() {
		m.Bet = m.Wallet.GetBalance()
	}
	m.State = StateBetting
	m.Result = ResultNone
	m.Message = "Place your bet and press ENTER to deal"
	m.PlayerHand = Hand{}
	m.DealerHand = Hand{}
	m.CanDouble = false
	m.CanSplit = false
	return m, nil
}

func (m *Model) View() string {
	s := ui.HeaderStyle.Render("🃏 BLACKJACK 🃏") + "\n\n"
	s += m.Wallet.Render() + "\n\n"

	if m.State != StateBetting {
		s += ui.SubtleStyle.Render("Dealer's Hand") + " "
		hideFirst := m.State == StatePlayerTurn || m.State == StateDealing
		s += ui.SubtleStyle.Render(fmt.Sprintf("[%s]", RenderHandValue(&m.DealerHand, hideFirst))) + "\n"
		s += RenderHand(m.DealerHand.Cards, hideFirst) + "\n\n"

		s += ui.SubtleStyle.Render("Your Hand") + " "
		s += ui.SubtleStyle.Render(fmt.Sprintf("[%d]", m.PlayerHand.Value())) + "\n"
		s += RenderHand(m.PlayerHand.Cards, false) + "\n\n"
	}

	s += fmt.Sprintf("Bet: $%.2f\n", m.Bet)
	s += m.Message + "\n\n"

	switch m.State {
	case StateBetting:
		s += ui.SubtleStyle.Render("[↑↓] adjust bet  [enter] deal  [q] quit")
	case StatePlayerTurn:
		s += ui.SubtleStyle.Render("[h] hit  [s] stand")
		if m.CanDouble {
			s += ui.SubtleStyle.Render("  [d] double")
		}
		s += ui.SubtleStyle.Render("  [q] quit")
	case StateResolution:
		s += ui.SubtleStyle.Render("[enter] new round  [q] quit")
	default:
		s += ui.SubtleStyle.Render("[q] quit")
	}

	return s
}
