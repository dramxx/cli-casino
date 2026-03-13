package ui

type WalletBackend interface {
	Bet(amount float64) bool
	Win(amount float64)
	Lose(amount float64)
	CanAfford(amount float64) bool
	Save() error
	GetBalance() float64
	Render() string
}
