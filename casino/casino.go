package casino

import (
	"fmt"

	"cli-casino/ui"
	tea "github.com/charmbracelet/bubbletea"
)

type GameID string

const (
	GameNone       GameID = ""
	GameSlots      GameID = "slots"
	GameRoulette   GameID = "roulette"
	GameBlackjack  GameID = "blackjack"
	GameVideoPoker GameID = "videopoker"
	GameStats      GameID = "stats"
)

type Model struct {
	Wallet     *Wallet
	Menu       ui.MenuModel
	ActiveGame GameID
	GameModel  tea.Model
	Width      int
	Height     int
}

func New() *Model {
	return &Model{
		Wallet:     NewWallet(),
		Menu:       ui.NewMenuModel(),
		ActiveGame: GameNone,
	}
}

func (m *Model) Init() tea.Cmd {
	m.Wallet.IncrementSession()
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.ActiveGame == GameNone {
			return m.updateMenu(msg)
		}
		return m.updateGame(msg)
	}

	if m.ActiveGame != GameNone && m.GameModel != nil {
		var cmd tea.Cmd
		m.GameModel, cmd = m.GameModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		m.Menu = m.Menu.Up()
	case "down":
		m.Menu = m.Menu.Down()
	case "enter":
		selected := m.Menu.GetSelected()
		switch selected {
		case "quit":
			m.Wallet.Save()
			return m, tea.Quit
		case "slots":
			m.ActiveGame = GameSlots
			return m, nil
		case "roulette":
			m.ActiveGame = GameRoulette
			return m, nil
		case "blackjack":
			m.ActiveGame = GameBlackjack
			return m, nil
		case "videopoker":
			m.ActiveGame = GameVideoPoker
			return m, nil
		case "stats":
			m.ActiveGame = GameStats
			return m, nil
		}
	case "q":
		m.Wallet.Save()
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) updateGame(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "q" || msg.String() == "esc" {
		m.ActiveGame = GameNone
		m.GameModel = nil
		return m, nil
	}

	if m.GameModel != nil {
		var cmd tea.Cmd
		m.GameModel, cmd = m.GameModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) View() string {
	if m.ActiveGame == GameNone {
		return m.Menu.View(m.Wallet.Render())
	}

	if m.ActiveGame == GameStats {
		return m.viewStats()
	}

	if m.GameModel != nil {
		return m.GameModel.View()
	}

	return m.viewComingSoon()
}

func (m *Model) viewComingSoon() string {
	gameName := string(m.ActiveGame)
	return ui.BoxStyle.Render(
		ui.TitleStyle.Render(gameName+"\n\n") +
			ui.SubtleStyle.Render("Coming soon!\n\nPress q to go back"),
	)
}

func (m *Model) viewStats() string {
	s := ui.HeaderStyle.Render("📊 STATS") + "\n\n"
	s += ui.BoxStyle.Render(fmt.Sprintf(
		"Balance:      $%.2f\nTotal Won:    $%.2f\nTotal Lost:   $%.2f\nSessions:     %d\nBiggest Win:  $%.2f",
		m.Wallet.Balance,
		m.Wallet.TotalWon,
		m.Wallet.TotalLost,
		m.Wallet.Sessions,
		m.Wallet.BiggestWin,
	))
	s += "\n\n" + ui.SubtleStyle.Render("Press q to go back")
	return s
}
