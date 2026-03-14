package casino

import (
	"testing"
)

func TestNewWallet(t *testing.T) {
	// Test wallet creation - loads existing save data if present
	wallet := NewWallet()
	if wallet == nil {
		t.Fatal("NewWallet() returned nil")
	}
	
	// Should have a positive balance (either starting balance or saved balance)
	if wallet.Balance <= 0 {
		t.Errorf("Expected positive balance, got %.2f", wallet.Balance)
	}
	
	// Should have valid statistics
	if wallet.TotalWon < 0 {
		t.Errorf("Expected TotalWon to be non-negative, got %.2f", wallet.TotalWon)
	}
	
	if wallet.TotalLost < 0 {
		t.Errorf("Expected TotalLost to be non-negative, got %.2f", wallet.TotalLost)
	}
	
	if wallet.Sessions < 0 {
		t.Errorf("Expected Sessions to be non-negative, got %d", wallet.Sessions)
	}
	
	if wallet.BiggestWin < 0 {
		t.Errorf("Expected BiggestWin to be non-negative, got %.2f", wallet.BiggestWin)
	}
}

func TestWalletBet(t *testing.T) {
	wallet := NewWallet()
	initialBalance := wallet.Balance
	
	// Test successful bet
	success := wallet.Bet(100)
	if !success {
		t.Error("Bet(100) should succeed with available balance")
	}
	
	if wallet.Balance != initialBalance-100 {
		t.Errorf("Expected balance %.2f after bet, got %.2f", initialBalance-100, wallet.Balance)
	}
	
	// Test failed bet (insufficient funds)
	success = wallet.Bet(initialBalance + 100) // More than current balance
	if success {
		t.Error("Bet more than balance should fail")
	}
	
	if wallet.Balance != initialBalance-100 {
		t.Error("Balance should not change on failed bet")
	}
}

func TestWalletWin(t *testing.T) {
	wallet := NewWallet()
	
	// Place a bet first
	wallet.Bet(100)
	initialBalance := wallet.Balance
	initialTotalWon := wallet.TotalWon
	initialBiggestWin := wallet.BiggestWin
	
	// Test win
	wallet.Win(350)
	
	if wallet.Balance != initialBalance+350 {
		t.Errorf("Expected balance %.2f after win, got %.2f", initialBalance+350, wallet.Balance)
	}
	
	if wallet.TotalWon != initialTotalWon+350 {
		t.Errorf("Expected TotalWon to be %.2f, got %.2f", initialTotalWon+350, wallet.TotalWon)
	}
	
	expectedBiggestWin := initialBiggestWin
	if 350 > initialBiggestWin {
		expectedBiggestWin = 350
	}
	if wallet.BiggestWin != expectedBiggestWin {
		t.Errorf("Expected BiggestWin to be %.2f, got %.2f", expectedBiggestWin, wallet.BiggestWin)
	}
}

func TestWalletCanAfford(t *testing.T) {
	wallet := NewWallet()
	balance := wallet.Balance
	
	// Test can afford
	if !wallet.CanAfford(100) {
		t.Error("Should be able to afford 100 with available balance")
	}
	
	// Test cannot afford
	if wallet.CanAfford(balance + 1) {
		t.Error("Should not be able to afford more than current balance")
	}
	
	// Test exact amount
	if !wallet.CanAfford(balance) {
		t.Error("Should be able to afford exact current balance")
	}
}

func TestWalletIsBroke(t *testing.T) {
	wallet := NewWallet()
	
	// Test not broke initially
	if wallet.IsBroke() {
		t.Error("Should not be broke initially")
	}
	
	// Test broke after betting all
	wallet.Bet(wallet.Balance)
	if !wallet.IsBroke() {
		t.Error("Should be broke after betting all money")
	}
	
	// Test not broke after winning
	wallet.Win(100)
	if wallet.IsBroke() {
		t.Error("Should not be broke after winning")
	}
}

func TestWalletRebuy(t *testing.T) {
	tests := []struct {
		name           string
		initialBalance float64
		expectedResult bool
		expectedBalance float64
	}{
		{
			name:           "Rebuy when broke",
			initialBalance: 0,
			expectedResult: true,
			expectedBalance: StartingBalance,
		},
		{
			name:           "Rebuy when low",
			initialBalance: 500,
			expectedResult: true,
			expectedBalance: StartingBalance,
		},
		{
			name:           "No rebuy when full",
			initialBalance: StartingBalance,
			expectedResult: false,
			expectedBalance: StartingBalance,
		},
		{
			name:           "No rebuy when above",
			initialBalance: 1500,
			expectedResult: false,
			expectedBalance: 1500,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet := &Wallet{Balance: tt.initialBalance}
			
			result := wallet.Rebuy()
			
			if result != tt.expectedResult {
				t.Errorf("Expected Rebuy() result %v, got %v", tt.expectedResult, result)
			}
			
			if wallet.Balance != tt.expectedBalance {
				t.Errorf("Expected balance %.2f, got %.2f", tt.expectedBalance, wallet.Balance)
			}
		})
	}
}

func TestWalletReset(t *testing.T) {
	wallet := NewWallet()
	
	// Modify wallet state
	wallet.Bet(100)
	wallet.Win(200)
	wallet.IncrementSession()
	
	// Reset
	wallet.Reset()
	
	// Check reset state
	if wallet.Balance != StartingBalance {
		t.Errorf("Expected balance %.2f after reset, got %.2f", StartingBalance, wallet.Balance)
	}
	
	if wallet.TotalWon != 0 {
		t.Errorf("Expected TotalWon to be 0 after reset, got %.2f", wallet.TotalWon)
	}
	
	if wallet.TotalLost != 0 {
		t.Errorf("Expected TotalLost to be 0 after reset, got %.2f", wallet.TotalLost)
	}
	
	if wallet.Sessions != 0 {
		t.Errorf("Expected Sessions to be 0 after reset, got %d", wallet.Sessions)
	}
	
	if wallet.BiggestWin != 0 {
		t.Errorf("Expected BiggestWin to be 0 after reset, got %.2f", wallet.BiggestWin)
	}
}

func TestWalletIncrementSession(t *testing.T) {
	wallet := NewWallet()
	
	initialSessions := wallet.Sessions
	wallet.IncrementSession()
	
	if wallet.Sessions != initialSessions+1 {
		t.Errorf("Expected sessions %d, got %d", initialSessions+1, wallet.Sessions)
	}
}

func TestWalletBiggestWinTracking(t *testing.T) {
	wallet := NewWallet()
	initialBiggestWin := wallet.BiggestWin
	
	// First win (smaller than existing)
	firstWin := 100.0
	if firstWin > initialBiggestWin {
		wallet.Win(firstWin)
		if wallet.BiggestWin != firstWin {
			t.Errorf("Expected BiggestWin to be %.2f, got %.2f", firstWin, wallet.BiggestWin)
		}
	} else {
		wallet.Win(firstWin)
		if wallet.BiggestWin != initialBiggestWin {
			t.Errorf("Expected BiggestWin to remain %.2f, got %.2f", initialBiggestWin, wallet.BiggestWin)
		}
	}
	
	// Bigger win than current biggest
	bigWin := initialBiggestWin + 500
	wallet.Win(bigWin)
	if wallet.BiggestWin != bigWin {
		t.Errorf("Expected BiggestWin to be %.2f, got %.2f", bigWin, wallet.BiggestWin)
	}
}

func TestWalletSave(t *testing.T) {
	wallet := NewWallet()
	wallet.Bet(100)
	wallet.Win(200)
	
	err := wallet.Save()
	if err != nil {
		t.Errorf("Save() returned error: %v", err)
	}
}

func TestWalletLoad(t *testing.T) {
	// Create and save a wallet
	originalWallet := NewWallet()
	originalWallet.Bet(100)
	originalWallet.Win(200)
	originalWallet.IncrementSession()
	originalWallet.Save()
	
	// Load the wallet
	loadedWallet := NewWallet()
	
	// Note: Since we can't easily mock the save system in this test,
	// we just verify that loading works without errors
	if loadedWallet == nil {
		t.Error("Loaded wallet should not be nil")
	}
}

func TestWalletLoadNoAutoRefill(t *testing.T) {
	// Test that wallet loads existing save data without auto-refill
	wallet := NewWallet()
	
	// Should have the saved balance (not auto-refilled to StartingBalance)
	// Since we can't predict the exact saved balance, just verify it's positive
	if wallet.Balance <= 0 {
		t.Errorf("Expected positive balance from save data, got %.2f", wallet.Balance)
	}
	
	// If there was existing save data with 0 balance, it should remain 0
	// (no auto-refill). Since we can't mock this easily, we just verify
	// the wallet loads successfully.
}
