package casino

import (
	"fmt"
	"time"

	"cli-casino/ui"

	"github.com/charmbracelet/lipgloss"
)

const StartingBalance = 1000.0

type WalletBackend interface {
	Bet(amount float64) bool
	Win(amount float64)
	Lose(amount float64)
	CanAfford(amount float64) bool
	Save() error
}

type Wallet struct {
	Balance    float64
	TotalWon   float64
	TotalLost  float64
	Sessions   int
	BiggestWin float64
}

func NewWallet() *Wallet {
	saveData, err := LoadSaveData()
	if err != nil {
		return &Wallet{
			Balance:    StartingBalance,
			TotalWon:   0,
			TotalLost:  0,
			Sessions:   0,
			BiggestWin: 0,
		}
	}

	wallet := &Wallet{
		Balance:    saveData.Balance,
		TotalWon:   saveData.TotalWon,
		TotalLost:  saveData.TotalLost,
		Sessions:   saveData.Sessions,
		BiggestWin: saveData.BiggestWin,
	}

	return wallet
}

func (w *Wallet) Bet(amount float64) bool {
	if amount > w.Balance {
		return false
	}
	w.Balance -= amount
	return true
}

func (w *Wallet) Win(amount float64) {
	w.Balance += amount
	w.TotalWon += amount
	if amount > w.BiggestWin {
		w.BiggestWin = amount
	}
}

func (w *Wallet) Loss(amount float64) {
	w.TotalLost += amount
}

func (w *Wallet) Lose(amount float64) {
	w.TotalLost += amount
}

func (w *Wallet) IsBroke() bool {
	return w.Balance <= 0
}

func (w *Wallet) CanAfford(amount float64) bool {
	return w.Balance >= amount
}

func (w *Wallet) GetBalance() float64 {
	return w.Balance
}

func (w *Wallet) Save() error {
	saveData := &SaveData{
		Balance:    w.Balance,
		TotalWon:   w.TotalWon,
		TotalLost:  w.TotalLost,
		Sessions:   w.Sessions,
		BiggestWin: w.BiggestWin,
		LastPlayed: time.Now().Format(time.RFC3339),
	}
	return SaveSaveData(saveData)
}

func (w *Wallet) IncrementSession() {
	w.Sessions++
}

func (w *Wallet) Reset() {
	w.Balance = StartingBalance
	w.TotalWon = 0
	w.TotalLost = 0
	w.Sessions = 0
	w.BiggestWin = 0
}

func (w *Wallet) Rebuy() bool {
	if w.Balance < StartingBalance {
		w.Balance = StartingBalance
		return true
	}
	return false
}

func (w *Wallet) Render() string {
	balanceStyle := ui.WalletStyle
	if w.IsBroke() {
		balanceStyle = balanceStyle.Foreground(ui.Lose)
	}
	return balanceStyle.Render(fmt.Sprintf("💰 $%.2f", w.Balance))
}

var WalletStyle = lipgloss.NewStyle().
	Foreground(ui.Primary).
	Padding(0, 1)
