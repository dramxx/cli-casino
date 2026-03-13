# CLI Casino 🎰

### Development Plan — Go + Bubbletea

---

## Vision

A terminal casino with real personality. ASCII art everywhere, satisfying animations,
a persistent wallet that carries across sessions. Feels like a place you'd want to hang
out in, not just a demo project. Built for Linux, playable anywhere with a decent terminal.

---

## Tech Stack

| Concern       | Choice                                     | Why                                    |
| ------------- | ------------------------------------------ | -------------------------------------- |
| Language      | Go 1.22+                                   | Learn Go, single binary output         |
| TUI framework | Bubbletea                                  | Elm architecture, great for game state |
| UI components | Bubbles + Lipgloss                         | Bubbletea's own component/styling libs |
| Persistence   | JSON flat file (`~/.cli-casino/save.json`) | Simple, no DB dependency               |
| ASCII art     | Hand-crafted + embedded in code            | Part of the fun                        |
| RNG           | `math/rand` with seeded source             | Fine for a casino game                 |

---

## Project Structure

```
cli-casino/
├── main.go
├── go.mod
├── go.sum
│
├── casino/
│   ├── casino.go          # Root model, menu, navigation
│   ├── wallet.go          # Balance, bet management, persistence
│   └── save.go            # Load/save player state to disk
│
├── games/
│   ├── slots/
│   │   ├── model.go       # Bubbletea model
│   │   ├── logic.go       # Spin logic, paylines, weights
│   │   └── art.go         # ASCII art, reel symbols, animations
│   ├── roulette/
│   │   ├── model.go
│   │   ├── logic.go       # Bet types, payout calculation
│   │   └── art.go         # Wheel, table layout
│   ├── blackjack/
│   │   ├── model.go
│   │   ├── logic.go       # Deck, dealer AI, hand evaluation
│   │   └── art.go         # Card ASCII art
│   └── videopoker/
│       ├── model.go
│       ├── logic.go       # Hand ranking, hold/draw mechanics
│       └── art.go         # Card art, hand name display
│
└── ui/
    ├── menu.go            # Main menu model
    ├── theme.go           # Lipgloss styles, color palette
    └── common.go          # Shared UI helpers (borders, headers, etc.)
```

---

## Architecture — Bubbletea Pattern

Each game is a **self-contained Bubbletea model** implementing `Init / Update / View`.
The root casino model owns navigation and passes control to the active game model.

```
CasinoModel
  └── activeGame (interface)
        ├── SlotsModel
        ├── RouletteModel
        ├── BlackjackModel
        └── VideoPokerModel
```

When a game ends, it emits a `GameResult` message back to the root model which
updates the wallet and returns to menu.

---

## Shared Systems (build these first)

### Wallet

- Starting balance: **$1000**
- Persisted to `~/.cli-casino/save.json`
- Tracks: balance, total won, total lost, sessions played, biggest win
- Displayed in header on every screen
- "Broke" state — player hits $0, offer to restart with fresh $1000

### Theme

Single color palette used everywhere via Lipgloss. Suggest:

- Background: terminal default (don't fight it)
- Primary accent: Gold `#FFD700`
- Win color: Green `#00FF88`
- Loss color: Red `#FF4444`
- Muted: Gray `#666666`
- Card/border: White `#FFFFFF`

### Main Menu ASCII Header

```
 ██████╗██╗     ██╗     ██████╗ █████╗ ███████╗██╗███╗   ██╗ ██████╗
██╔════╝██║     ██║    ██╔════╝██╔══██╗██╔════╝██║████╗  ██║██╔═══██╗
██║     ██║     ██║    ██║     ███████║███████╗██║██╔██╗ ██║██║   ██║
██║     ██║     ██║    ██║     ██╔══██║╚════██║██║██║╚██╗██║██║   ██║
╚██████╗███████╗██║    ╚██████╗██║  ██║███████║██║██║ ╚████║╚██████╔╝
 ╚═════╝╚══════╝╚═╝     ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝╚═╝  ╚═══╝ ╚═════╝
```

---

## Game 1 — Slots 🎰

**Why first:** Simplest logic, lets you nail the ASCII animation system early.

### Symbols & Weights

```
Symbol    Display    Weight    Payout (x bet)
SEVEN       7️        1          50x
BAR        BAR        3          10x
BELL       🔔         5           5x
CHERRY     🍒         8           3x
LEMON      🍋        10           2x
GRAPE      🍇        12           1.5x
```

(Use text representations in actual ASCII, not emoji — terminal compatibility)

### Reels ASCII

```
╔═══════════════════════╗
║   ┌─────────────────┐ ║
║   │  LEMON  │ BELL  │ CHERRY │
║   │ >>>>  7  │ >>>>  7  │ >>>>  7  │   ← JACKPOT
║   │  BAR   │ GRAPE │ LEMON  │
║   └─────────────────┘ ║
╚═══════════════════════╝
```

### Spin Animation

- Reels spin independently, stop left-to-right with a small delay between each
- Implement as a ticker: each tick scrolls the reel by one symbol
- Reel 1 stops at tick 8, Reel 2 at tick 12, Reel 3 at tick 16
- Flash the win line on match

### Controls

- `SPACE` — spin (deducts bet)
- `↑↓` — adjust bet (min $10, max $100 or 10% of balance)
- `q` — back to menu

---

## Game 2 — Roulette 🎡

**Why second:** Betting system is the interesting part, wheel is mostly art.

### Bet Types to Implement

- Straight up (single number) — 35:1
- Red / Black — 1:1
- Odd / Even — 1:1
- Low (1-18) / High (19-36) — 1:1
- Dozens (1-12, 13-24, 25-36) — 2:1

### The Wheel ASCII

Not a full circle — render as a spinning strip of numbers that wraps:

```
┌──────────────────────────────────────────┐
│ .. 26 │ 3 │ 35 │ 12 │ 28 │ 7 │ 29 │ 18 ..│
│              ▲ BALL                       │
└──────────────────────────────────────────┘
```

Ball position slides across ticks, decelerates, lands.

### Betting Table

Render a simplified version of the felt layout — numbers 0-36 in a grid,
player navigates with arrow keys and places chips.

### Controls

- Arrow keys — navigate betting grid
- `SPACE` / `ENTER` — place chip on selection
- `c` — clear bets
- `s` — spin
- `q` — back to menu

---

## Game 3 — Blackjack 🃏

**The meat.** Most interesting state machine of the bunch.

### Card ASCII Art

```
┌─────────┐   ┌─────────┐   ┌─────────┐
│ A       │   │ K       │   │░░░░░░░░░│
│         │   │         │   │░░░░░░░░░│
│    ♠    │   │    ♥    │   │░░░░░░░░░│  ← face down
│         │   │         │   │░░░░░░░░░│
│       A │   │       K │   │░░░░░░░░░│
└─────────┘   └─────────┘   └─────────┘
```

### Game States

```
Betting → Dealing → PlayerTurn → DealerTurn → Resolution → Betting
```

### Rules

- Standard 6-deck shoe, reshuffled when < 25% remains
- Dealer hits on soft 16, stands on soft 17
- Blackjack pays 3:2
- Double down allowed on first two cards
- Split allowed on pairs (keep it to one split for now)
- No insurance (adds complexity, low fun return)

### Dealer AI

Simple rule-based, no card counting. Just follows house rules.
Display dealer's logic as flavor text: _"Dealer hits on 14..."_

### Controls

- `h` — hit
- `s` — stand
- `d` — double down
- `p` — split (when available)
- `ENTER` — confirm bet and deal
- `q` — back to menu

---

## Game 4 — Video Poker 🎲

**The hardest.** Hand evaluation is non-trivial.

### Variant: Jacks or Better

Standard, well-understood, good payouts to implement.

### Paytable

```
Royal Flush        800x
Straight Flush     50x
Four of a Kind     25x
Full House          9x
Flush               6x
Straight            4x
Three of a Kind     3x
Two Pair            2x
Jacks or Better     1x
```

### Hand Evaluation

This is the interesting engineering problem. You need to classify any 5-card hand.
Suggested approach: sort by rank, then pattern match on rank frequencies + flush/straight checks.

### Hold/Draw UI

```
┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
│ A       │  │ A       │  │ K       │  │ 3       │  │ 7       │
│    ♠    │  │    ♥    │  │    ♠    │  │    ♣    │  │    ♦    │
│       A │  │       A │  │       K │  │       3 │  │       7 │
└─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘
  [HOLD]       [HOLD]                                            ← cursor
```

### Controls

- `1-5` or arrow keys — select card to toggle hold
- `SPACE` — toggle hold on selected card
- `ENTER` — draw (replace non-held cards)
- `q` — back to menu

---

## Development Phases

### Phase 0 — Scaffolding (30 min)

- [ ] `go mod init cli-casino`
- [ ] Install bubbletea, bubbles, lipgloss
- [ ] Main menu with navigation (no games yet)
- [ ] Wallet system with persistence
- [ ] Theme and shared UI components
- [ ] "Coming soon" placeholder for each game

### Phase 1 — Slots (1 hour)

- [ ] Reel data structure and spin logic
- [ ] Payout calculation
- [ ] Spinning animation with ticker
- [ ] Win/lose feedback
- [ ] Bet adjustment UI
- [ ] Wire up to wallet

### Phase 2 — Roulette (1 hour)

- [ ] Bet types and payout math
- [ ] Betting table navigation
- [ ] Wheel animation
- [ ] Multiple simultaneous bets
- [ ] Wire up to wallet

### Phase 3 — Blackjack (1.5 hours)

- [ ] Deck and card types
- [ ] Deal, hit, stand, bust logic
- [ ] Dealer AI
- [ ] Double down + split
- [ ] Full ASCII card rendering
- [ ] Wire up to wallet

### Phase 4 — Video Poker (1.5 hours)

- [ ] Hand evaluator (the fun part)
- [ ] Hold/draw mechanics
- [ ] Paytable and payout calculation
- [ ] Hold UI with card selection
- [ ] Wire up to wallet

### Phase 5 — Polish (whenever)

- [ ] Stats screen (biggest win, total played, etc.)
- [ ] Sounds via `beep` if feeling adventurous
- [ ] Smooth transition animations between screens
- [ ] `--help` and `--reset` CLI flags
- [ ] README with demo gif

---

## Nice-to-haves (don't plan for these, just keep in mind)

- Daily bonus chips (check last session timestamp)
- High score / leaderboard stored locally
- Different "rooms" with higher stakes
- Slot machine themes (Vegas, Fruit, Sci-Fi)
- `CASINO_SEED` env var for reproducible RNG (debugging)

---

## Linux-first, but Windows-friendly rules

- Use `lipgloss` for all styling — it handles terminal capability detection
- Avoid ANSI codes directly, let lipgloss/bubbletea handle it
- Use `os.UserHomeDir()` for save path, not hardcoded `~`
- Box-drawing chars (─ │ ╔ etc.) are fine — they work in Windows Terminal
- Avoid terminal size assumptions — query with `bubbletea.WindowSizeMsg`
- Don't use `clear` shell command — use bubbletea's full-screen alt buffer mode

---

## Getting Started Commands

```bash
mkdir cli-casino && cd cli-casino
go mod init cli-casino
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
```

---

_Let's go. Start with Phase 0 and don't touch slots until the menu feels good._
