package ui

type MenuItem struct {
	Name        string
	Description string
	ID          string
}

var MenuItems = []MenuItem{
	{Name: "Slots", Description: "Spin the reels!", ID: "slots"},
	{Name: "Roulette", Description: "Place your bets", ID: "roulette"},
	{Name: "Blackjack", Description: "Beat the dealer", ID: "blackjack"},
	{Name: "Video Poker", Description: "Draw to win", ID: "videopoker"},
	{Name: "Stats", Description: "View your stats", ID: "stats"},
	{Name: "Quit", Description: "Exit the casino", ID: "quit"},
}

type MenuModel struct {
	Selected int
}

func NewMenuModel() MenuModel {
	return MenuModel{
		Selected: 0,
	}
}

func (m MenuModel) Up() MenuModel {
	m.Selected--
	if m.Selected < 0 {
		m.Selected = len(MenuItems) - 1
	}
	return m
}

func (m MenuModel) Down() MenuModel {
	m.Selected = (m.Selected + 1) % len(MenuItems)
	return m
}

func (m MenuModel) GetSelected() string {
	if m.Selected >= 0 && m.Selected < len(MenuItems) {
		return MenuItems[m.Selected].ID
	}
	return ""
}

func (m MenuModel) View(wallet string) string {
	var s string

	s += HeaderStyle.Render(`  
 ██████╗██╗     ██╗      ██████╗ █████╗ ███████╗██╗███╗   ██╗ ██████╗
██╔════╝██║     ██║     ██╔════╝██╔══██╗██╔════╝██║████╗  ██║██╔═══██╗
██║     ██║     ██║     ██║     ███████║███████╗██║██╔██╗ ██║██║   ██║
██║     ██║     ██║     ██║     ██╔══██║╚════██║██║██║╚██╗██║██║   ██║
╚██████╗███████╗██║     ╚██████╗██║  ██║███████║██║██║ ╚████║╚██████╔╝
 ╚═════╝╚══════╝╚═╝      ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝╚═╝  ╚═══╝ ╚═════╝
`) + "\n\n"
	s += wallet + "\n\n"

	for i, item := range MenuItems {
		if i == m.Selected {
			s += "  " + MenuItemSelectedStyle.Render("▶ "+item.Name)
		} else {
			s += "   " + MenuItemStyle.Render(item.Name)
		}
		s += "\n"
		if i == m.Selected {
			s += SubtleStyle.Render("    "+item.Description) + "\n"
		} else {
			s += "\n"
		}
	}

	s += "\n" + SubtleStyle.Render("  ↑↓ navigate • enter select • q quit")

	return s
}
