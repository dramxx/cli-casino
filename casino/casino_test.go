package casino

import (
	"cli-casino/ui"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewCasino(t *testing.T) {
	model := New()
	
	if model == nil {
		t.Fatal("New() returned nil")
	}
	
	if model.Wallet == nil {
		t.Error("Wallet should not be nil")
	}
	
	if model.ActiveGame != GameNone {
		t.Errorf("Expected ActiveGame to be GameNone, got %v", model.ActiveGame)
	}
	
	if model.GameModel != nil {
		t.Error("GameModel should be nil initially")
	}
}

func TestCasinoUpdateMenuNavigation(t *testing.T) {
	model := New()
	
	// Test up navigation
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	newModel, _ := model.Update(upMsg)
	
	casinoModel := newModel.(*Model)
	if casinoModel.Menu.Selected != len(ui.MenuItems)-1 {
		t.Errorf("Expected selected to wrap to last item, got %d", casinoModel.Menu.Selected)
	}
	
	// Test down navigation
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	newModel, _ = casinoModel.Update(downMsg)
	
	casinoModel = newModel.(*Model)
	if casinoModel.Menu.Selected != 0 {
		t.Errorf("Expected selected to wrap to first item, got %d", casinoModel.Menu.Selected)
	}
}

func TestCasinoUpdateMenuSelection(t *testing.T) {
	model := New()
	
	// Navigate to roulette
	model.Menu.Selected = 1 // Roulette is at index 1
	
	// Select roulette
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	
	casinoModel := newModel.(*Model)
	if casinoModel.ActiveGame != GameRoulette {
		t.Errorf("Expected ActiveGame to be GameRoulette, got %v", casinoModel.ActiveGame)
	}
	
	if casinoModel.GameModel == nil {
		t.Error("GameModel should not be nil after selecting a game")
	}
}

func TestCasinoUpdateQuit(t *testing.T) {
	model := New()
	
	// Test quit with 'q'
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(quitMsg)
	
	if cmd == nil {
		t.Error("Expected quit command")
	}
}

func TestCasinoUpdateRebuy(t *testing.T) {
	model := New()
	
	// Set wallet to low balance
	model.Wallet.Bet(900) // Leave 100 balance
	
	// Navigate to rebuy (should be at index 5 after stats)
	model.Menu.Selected = 5
	model.ActiveGame = GameNone
	
	// Select rebuy
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	
	casinoModel := newModel.(*Model)
	if casinoModel.Wallet.Balance != StartingBalance {
		t.Errorf("Expected balance to be %.2f after rebuy, got %.2f", StartingBalance, casinoModel.Wallet.Balance)
	}
	
	// Should still be in menu (not start a game)
	if casinoModel.ActiveGame != GameNone {
		t.Errorf("Expected to remain in menu after rebuy, got %v", casinoModel.ActiveGame)
	}
}

func TestCasinoUpdateRebuyNoEffect(t *testing.T) {
	model := New()
	
	// Set wallet to full balance
	model.Wallet.Balance = 1500
	
	// Navigate to rebuy
	model.Menu.Selected = 5
	model.ActiveGame = GameNone
	
	initialBalance := model.Wallet.Balance
	
	// Select rebuy
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	
	casinoModel := newModel.(*Model)
	if casinoModel.Wallet.Balance != initialBalance {
		t.Errorf("Expected balance to remain %.2f, got %.2f", initialBalance, casinoModel.Wallet.Balance)
	}
}

func TestCasinoUpdateGameInput(t *testing.T) {
	model := New()
	
	// Start a game
	model.Menu.Selected = 1 // Roulette
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	model = newModel.(*Model)
	
	// Test game input passes through
	gameMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}
	newModel, _ = model.Update(gameMsg)
	
	// Should still be in game (model unchanged by unknown game input)
	if newModel.(*Model).ActiveGame != GameRoulette {
		t.Error("Should remain in game after game input")
	}
}

func TestCasinoUpdateGameExit(t *testing.T) {
	model := New()
	
	// Start a game
	model.Menu.Selected = 1 // Roulette
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	model = newModel.(*Model)
	
	// Exit game with 'q'
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, _ = model.Update(quitMsg)
	
	casinoModel := newModel.(*Model)
	if casinoModel.ActiveGame != GameNone {
		t.Errorf("Expected ActiveGame to be GameNone after exit, got %v", casinoModel.ActiveGame)
	}
	
	if casinoModel.GameModel != nil {
		t.Error("GameModel should be nil after exit")
	}
}

func TestCasinoUpdateGameExitWithEsc(t *testing.T) {
	model := New()
	
	// Start a game
	model.Menu.Selected = 1 // Roulette
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	model = newModel.(*Model)
	
	// Exit game with 'esc'
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ = model.Update(escMsg)
	
	casinoModel := newModel.(*Model)
	if casinoModel.ActiveGame != GameNone {
		t.Errorf("Expected ActiveGame to be GameNone after esc, got %v", casinoModel.ActiveGame)
	}
}

func TestCasinoUpdateWindowSize(t *testing.T) {
	model := New()
	
	windowMsg := tea.WindowSizeMsg{Width: 80, Height: 24}
	newModel, _ := model.Update(windowMsg)
	
	casinoModel := newModel.(*Model)
	if casinoModel.Width != 80 {
		t.Errorf("Expected Width to be 80, got %d", casinoModel.Width)
	}
	
	if casinoModel.Height != 24 {
		t.Errorf("Expected Height to be 24, got %d", casinoModel.Height)
	}
}

func TestCasinoViewMenu(t *testing.T) {
	model := New()
	
	view := model.View()
	if view == "" {
		t.Error("View() should not return empty string")
	}
	
	// Should contain wallet info
	if model.Wallet.Render() == "" {
		t.Error("Wallet render should not be empty")
	}
	
	// Should contain menu items
	if len(ui.MenuItems) == 0 {
		t.Error("MenuItems should not be empty")
	}
}

func TestCasinoViewGame(t *testing.T) {
	model := New()
	
	// Start roulette game
	model.Menu.Selected = 1
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	model = newModel.(*Model)
	
	view := model.View()
	if view == "" {
		t.Error("View() should not return empty string when in game")
	}
}

func TestCasinoViewStats(t *testing.T) {
	model := New()
	
	// Navigate to stats
	model.Menu.Selected = 4 // Stats
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterMsg)
	model = newModel.(*Model)
	
	view := model.View()
	if view == "" {
		t.Error("View() should not return empty string for stats")
	}
	
	// Should contain stats information
	if len(view) < 50 {
		t.Error("Stats view should contain substantial content")
	}
}

func TestCasinoInit(t *testing.T) {
	model := New()
	
	initialSessions := model.Wallet.Sessions
	cmd := model.Init()
	
	if cmd != nil {
		t.Error("Init() should not return a command")
	}
	
	if model.Wallet.Sessions != initialSessions+1 {
		t.Errorf("Expected sessions to increment from %d to %d", initialSessions, initialSessions+1)
	}
}

func TestCasinoMenuGameMapping(t *testing.T) {
	tests := []struct {
		menuIndex   int
		expectedGame GameID
	}{
		{0, GameSlots},
		{1, GameRoulette},
		{2, GameBlackjack},
		{3, GameVideoPoker},
		{4, GameStats},
		{5, GameNone}, // Rebuy
		{6, GameNone}, // Quit
	}
	
	model := New()
	
	for _, tt := range tests {
		t.Run(string(tt.expectedGame), func(t *testing.T) {
			model.Menu.Selected = tt.menuIndex
			model.ActiveGame = GameNone
			
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, _ := model.Update(enterMsg)
			
			casinoModel := newModel.(*Model)
			
			// Special case for quit and rebuy
			if tt.expectedGame == GameNone {
				if tt.menuIndex == 6 { // Quit
					// Quit returns a command, we can't easily test that here
					// but the game should remain GameNone
					if casinoModel.ActiveGame != GameNone {
						t.Errorf("Expected GameNone for quit, got %v", casinoModel.ActiveGame)
					}
				}
				// Rebuy also keeps GameNone
			} else {
				if casinoModel.ActiveGame != tt.expectedGame {
					t.Errorf("Expected %v, got %v", tt.expectedGame, casinoModel.ActiveGame)
				}
			}
		})
	}
}
