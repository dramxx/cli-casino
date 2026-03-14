# CLI Casino 🎰

A fully-featured terminal casino with real personality. Built with Go and Bubbletea, featuring ASCII art, smooth animations, and persistent wallet state across sessions.

## Features

- **4 Complete Games**
  - 🎰 **Slots** - Classic 3-reel slot machine with weighted symbols
  - 🎡 **Roulette** - Full European roulette with multiple bet types
  - 🃏 **Blackjack** - Standard rules with double down and split
  - 🎲 **Video Poker** - Jacks or Better variant with proper hand evaluation

- **Persistent Wallet**
  - Starting balance: $1000
  - Tracks total won, total lost, sessions played, and biggest win
  - Saves automatically to `~/.cli-casino/save.json`

- **Beautiful TUI**
  - ASCII art everywhere
  - Smooth animations
  - Color-coded wins/losses
  - Responsive keyboard controls

## Installation

### Prerequisites
- Go 1.22 or higher

### Build from Source

```bash
git clone <repository-url>
cd cli-casino
go build -o cli-casino
```

### Run

```bash
./cli-casino
```

On Windows:
```bash
.\cli-casino.exe
```

## How to Play

### Main Menu
- `↑↓` - Navigate menu
- `Enter` - Select game
- `q` - Quit

### Slots 🎰
- `Space` - Spin the reels
- `↑↓` - Adjust bet ($10-$100)
- `q` - Back to menu

**Payouts:**
- 7️⃣ (Three 7s) - 50x bet
- BAR (Three BARs) - 10x bet
- BELL (Three BELLs) - 5x bet
- CHERRY (Three CHERRYs) - 3x bet
- LEMON (Three LEMONs) - 2x bet
- GRAPE (Three GRAPEs) - 1.5x bet

### Roulette 🎡
- `Tab` - Switch between betting sections (Main, Zero, Dozens, Outside)
- `Arrow Keys` - Navigate within section (edges transition to neighboring sections)
- `Space` - Place bet on selected number/area
- `c` - Clear all bets
- `s` - Spin the wheel
- `q` - Back to menu

**Bet Types:**
- Straight up (single number) - 35:1
- Dozens (1-12, 13-24, 25-36) - 2:1
- Red/Black - 1:1
- Odd/Even - 1:1
- Low (1-18) / High (19-36) - 1:1

**Navigation:**
- In Main grid: arrows move through 3x12 grid
- At grid edges: arrows transition to neighboring sections
- Tab/Shift+Tab: cycle through sections

### Blackjack 🃏
- `↑↓` - Adjust bet during betting phase
- `Enter` - Deal cards
- `h` - Hit (take another card)
- `s` - Stand (end your turn)
- `d` - Double down (double bet, take one card, then stand)
- `q` - Back to menu

**Rules:**
- 6-deck shoe, reshuffled when < 25% remains
- Dealer hits on soft 16, stands on soft 17
- Blackjack pays 3:2
- Double down allowed on first two cards
- Split allowed on pairs (single split only)

### Video Poker 🎲
- `↑↓` - Adjust bet during betting phase ($5-$25)
- `Enter` - Deal cards / Draw new cards
- `←→` or `1-5` - Select card
- `Space` - Toggle hold on selected card
- `q` - Back to menu

**Jacks or Better Paytable:**
- Royal Flush - 800x
- Straight Flush - 50x
- Four of a Kind - 25x
- Full House - 9x
- Flush - 6x
- Straight - 4x
- Three of a Kind - 3x
- Two Pair - 2x
- Jacks or Better - 1x

## Project Structure

```
cli-casino/
├── main.go                 # Entry point
├── casino/                 # Root model and wallet
│   ├── casino.go          # Main casino controller
│   ├── wallet.go          # Balance management
│   └── save.go            # Persistence layer
├── games/                 # Individual game implementations
│   ├── slots/
│   │   ├── model.go       # Bubbletea model
│   │   ├── logic.go       # Game logic
│   │   └── art.go         # ASCII rendering
│   ├── roulette/
│   ├── blackjack/
│   └── videopoker/
└── ui/                    # Shared UI components
    ├── menu.go            # Main menu
    ├── theme.go           # Color palette and styles
    ├── common.go          # Shared helpers
    └── wallet.go          # Wallet interface
```

## Architecture

Built using the Elm architecture via [Bubbletea](https://github.com/charmbracelet/bubbletea):
- Each game is a self-contained model implementing `Init / Update / View`
- Root casino model handles navigation and wallet state
- Games communicate results back to update the wallet
- All styling via [Lipgloss](https://github.com/charmbracelet/lipgloss)

## Testing

Run all tests:
```bash
go test ./...
```

Run tests for a specific game:
```bash
go test ./games/slots -v
go test ./games/blackjack -v
go test ./games/videopoker -v
go test ./games/roulette -v
```

## Development

### Dependencies
```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
```

### Adding a New Game
1. Create a new directory under `games/`
2. Implement `model.go` with Bubbletea interface
3. Add game logic in `logic.go`
4. Create ASCII art in `art.go`
5. Wire up in `casino/casino.go`
6. Add menu item in `ui/menu.go`

## Save File Location

The wallet state is saved to:
- **Linux/macOS**: `~/.cli-casino/save.json`
- **Windows**: `%USERPROFILE%\.cli-casino\save.json`

## Tips

- Start with small bets to learn each game
- In Roulette, you can place multiple bets before spinning
- In Blackjack, double down on 10 or 11 for best odds
- In Video Poker, always hold pairs of Jacks or better
- Check the Stats screen to track your performance

## Credits

Built with:
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## License

MIT License - See LICENSE file for details

---

**Good luck and have fun! 🎰🎡🃏🎲**
