package roulette

import (
	"fmt"
	"testing"
)

// MockWallet implements WalletBackend for testing
type MockWallet struct {
	Balance float64
}

func (m *MockWallet) Bet(amount float64) bool {
	if amount > m.Balance {
		return false
	}
	m.Balance -= amount
	return true
}

func (m *MockWallet) Win(amount float64) {
	m.Balance += amount
}

func (m *MockWallet) Lose(amount float64) {
	m.Balance -= amount
}

func (m *MockWallet) CanAfford(amount float64) bool {
	return m.Balance >= amount
}

func (m *MockWallet) Save() error {
	return nil
}

func (m *MockWallet) GetBalance() float64 {
	return m.Balance
}

func (m *MockWallet) Render() string {
	return fmt.Sprintf("💰 $%.2f", m.Balance)
}

func TestRouletteNewNavigation(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	
	// Test initial state
	if model.Section != SectionMain {
		t.Errorf("Expected initial section to be SectionMain, got %v", model.Section)
	}
	
	if model.Row != 0 {
		t.Errorf("Expected initial row to be 0, got %d", model.Row)
	}
	
	if model.Col != 0 {
		t.Errorf("Expected initial col to be 0, got %d", model.Col)
	}
	
	if model.Index != 0 {
		t.Errorf("Expected initial index to be 0, got %d", model.Index)
	}
}

func TestRouletteSectionNavigation(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	
	// Test tab navigation (next section)
	model.nextSection()
	if model.Section != SectionDozens {
		t.Errorf("Expected SectionDozens after first tab, got %v", model.Section)
	}
	
	model.nextSection()
	if model.Section != SectionLowHigh {
		t.Errorf("Expected SectionLowHigh after second tab, got %v", model.Section)
	}
	
	model.nextSection()
	if model.Section != SectionOutside {
		t.Errorf("Expected SectionOutside after third tab, got %v", model.Section)
	}
	
	model.nextSection()
	if model.Section != SectionMain {
		t.Errorf("Expected SectionMain after fourth tab (wrap), got %v", model.Section)
	}
	
	// Test shift+tab navigation (previous section)
	model.prevSection()
	if model.Section != SectionOutside {
		t.Errorf("Expected SectionOutside after shift+tab, got %v", model.Section)
	}
}

func TestRouletteMainSectionNavigation(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	
	// Test right navigation
	model.moveRight()
	if model.Col != 1 {
		t.Errorf("Expected col to be 1 after moveRight, got %d", model.Col)
	}
	
	// Test left navigation
	model.moveLeft()
	if model.Col != 0 {
		t.Errorf("Expected col to be 0 after moveLeft, got %d", model.Col)
	}
	
	// Test down navigation
	model.moveDown()
	if model.Row != 1 {
		t.Errorf("Expected row to be 1 after moveDown, got %d", model.Row)
	}
	
	// Test up navigation
	model.moveUp()
	if model.Row != 0 {
		t.Errorf("Expected row to be 0 after moveUp, got %d", model.Row)
	}
}

func TestRouletteMainSectionBoundaries(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	
	// Test right boundary (should not go beyond col 11)
	model.Col = 11
	model.moveRight()
	if model.Col != 11 {
		t.Errorf("Expected col to remain 11 at boundary, got %d", model.Col)
	}
	
	// Test left boundary (should not go below col 0)
	model.Col = 0
	model.moveLeft()
	if model.Col != 0 {
		t.Errorf("Expected col to remain 0 at boundary, got %d", model.Col)
	}
	
	// Test down boundary (should move to Dozens when at row 2)
	model.Row = 2
	model.moveDown()
	if model.Section != SectionDozens {
		t.Errorf("Expected to move to Dozens section, got %v", model.Section)
	}
	
	// Reset to Main for up test
	model.Section = SectionMain
	model.Row = 0
	model.moveUp()
	if model.Row != 0 {
		t.Errorf("Expected row to remain 0 at boundary, got %d", model.Row)
	}
}

func TestRouletteDozensSectionNavigation(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	model.Section = SectionDozens
	
	// Test right navigation
	model.moveRight()
	if model.Index != 1 {
		t.Errorf("Expected index to be 1 after moveRight in dozens, got %d", model.Index)
	}
	
	model.moveRight()
	if model.Index != 2 {
		t.Errorf("Expected index to be 2 after second moveRight in dozens, got %d", model.Index)
	}
	
	// Test right boundary
	model.moveRight()
	if model.Index != 2 {
		t.Errorf("Expected index to remain 2 at boundary in dozens, got %d", model.Index)
	}
	
	// Test left navigation
	model.moveLeft()
	if model.Index != 1 {
		t.Errorf("Expected index to be 1 after moveLeft in dozens, got %d", model.Index)
	}
	
	// Test left boundary
	model.Index = 0
	model.moveLeft()
	if model.Index != 0 {
		t.Errorf("Expected index to remain 0 at boundary in dozens, got %d", model.Index)
	}
}

func TestRouletteLowHighSectionNavigation(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	model.Section = SectionLowHigh
	
	// Test right navigation (0 -> 1-18 -> 19-36)
	model.moveRight()
	if model.Index != 1 {
		t.Errorf("Expected index to be 1 after moveRight in low/high, got %d", model.Index)
	}
	
	model.moveRight()
	if model.Index != 2 {
		t.Errorf("Expected index to be 2 after second moveRight in low/high, got %d", model.Index)
	}
	
	// Test right boundary
	model.moveRight()
	if model.Index != 2 {
		t.Errorf("Expected index to remain 2 at boundary in low/high, got %d", model.Index)
	}
	
	// Test left navigation
	model.moveLeft()
	if model.Index != 1 {
		t.Errorf("Expected index to be 1 after moveLeft in low/high, got %d", model.Index)
	}
	
	// Test left boundary
	model.Index = 0
	model.moveLeft()
	if model.Index != 0 {
		t.Errorf("Expected index to remain 0 at boundary in low/high, got %d", model.Index)
	}
}

func TestRouletteOutsideSectionNavigation(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	model.Section = SectionOutside
	
	// Test right navigation (EVEN -> RED -> BLACK -> ODD)
	model.moveRight()
	if model.Index != 1 {
		t.Errorf("Expected index to be 1 after moveRight in outside, got %d", model.Index)
	}
	
	model.moveRight()
	if model.Index != 2 {
		t.Errorf("Expected index to be 2 after second moveRight in outside, got %d", model.Index)
	}
	
	model.moveRight()
	if model.Index != 3 {
		t.Errorf("Expected index to be 3 after third moveRight in outside, got %d", model.Index)
	}
	
	// Test right boundary
	model.moveRight()
	if model.Index != 3 {
		t.Errorf("Expected index to remain 3 at boundary in outside, got %d", model.Index)
	}
	
	// Test left navigation
	model.moveLeft()
	if model.Index != 2 {
		t.Errorf("Expected index to be 2 after moveLeft in outside, got %d", model.Index)
	}
	
	// Test left boundary
	model.Index = 0
	model.moveLeft()
	if model.Index != 0 {
		t.Errorf("Expected index to remain 0 at boundary in outside, got %d", model.Index)
	}
}

func TestRouletteVerticalNavigationBetweenSections(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	
	// Test moving down from Main to Dozens
	model.Row = 2
	model.moveDown()
	if model.Section != SectionDozens {
		t.Errorf("Expected SectionDozens when moving down from Main, got %v", model.Section)
	}
	if model.Index != 0 {
		t.Errorf("Expected index to be 0 when moving to Dozens, got %d", model.Index)
	}
	
	// Test moving up from Dozens to Main
	model.moveUp()
	if model.Section != SectionMain {
		t.Errorf("Expected SectionMain when moving up from Dozens, got %v", model.Section)
	}
	if model.Row != 2 {
		t.Errorf("Expected row to be 2 when moving to Main, got %d", model.Row)
	}
	
	// Test moving down from Main to Dozens again
	model.Row = 2
	model.moveDown()
	if model.Section != SectionDozens {
		t.Errorf("Expected SectionDozens when moving down from Main, got %v", model.Section)
	}
	
	// Test moving down from Dozens to LowHigh
	model.moveDown()
	if model.Section != SectionLowHigh {
		t.Errorf("Expected SectionLowHigh when moving down from Dozens, got %v", model.Section)
	}
	
	// Test moving down from LowHigh to Outside
	model.moveDown()
	if model.Section != SectionOutside {
		t.Errorf("Expected SectionOutside when moving down from LowHigh, got %v", model.Section)
	}
	
	// Test moving up from Outside to LowHigh
	model.moveUp()
	if model.Section != SectionLowHigh {
		t.Errorf("Expected SectionLowHigh when moving up from Outside, got %v", model.Section)
	}
}

func TestRouletteGetSelectedMain(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	
	// Test selection in main section
	model.Row = 0
	model.Col = 0
	selected := model.getSelected()
	if selected != 1 {
		t.Errorf("Expected selection 1 at row 0, col 0, got %d", selected)
	}
	
	model.Row = 1
	model.Col = 11
	selected = model.getSelected()
	if selected != 35 {
		t.Errorf("Expected selection 35 at row 1, col 11, got %d", selected)
	}
	
	model.Row = 2
	model.Col = 5
	selected = model.getSelected()
	if selected != 18 {
		t.Errorf("Expected selection 18 at row 2, col 5, got %d", selected)
	}
}

func TestRouletteGetSelectedDozens(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	model.Section = SectionDozens
	
	// Test selection in dozens section
	model.Index = 0
	selected := model.getSelected()
	if selected != 101 {
		t.Errorf("Expected selection 101 for first dozen, got %d", selected)
	}
	
	model.Index = 1
	selected = model.getSelected()
	if selected != 102 {
		t.Errorf("Expected selection 102 for second dozen, got %d", selected)
	}
	
	model.Index = 2
	selected = model.getSelected()
	if selected != 103 {
		t.Errorf("Expected selection 103 for third dozen, got %d", selected)
	}
}

func TestRouletteGetSelectedLowHigh(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	model.Section = SectionLowHigh
	
	// Test selection in low/high section
	model.Index = 0
	selected := model.getSelected()
	if selected != 0 {
		t.Errorf("Expected selection 0 for index 0 in low/high, got %d", selected)
	}
	
	model.Index = 1
	selected = model.getSelected()
	if selected != 104 {
		t.Errorf("Expected selection 104 for 1-18, got %d", selected)
	}
	
	model.Index = 2
	selected = model.getSelected()
	if selected != 105 {
		t.Errorf("Expected selection 105 for 19-36, got %d", selected)
	}
}

func TestRouletteGetSelectedOutside(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	model.Section = SectionOutside
	
	// Test selection in outside section
	model.Index = 0
	selected := model.getSelected()
	if selected != 106 {
		t.Errorf("Expected selection 106 for EVEN, got %d", selected)
	}
	
	model.Index = 1
	selected = model.getSelected()
	if selected != 108 {
		t.Errorf("Expected selection 108 for RED, got %d", selected)
	}
	
	model.Index = 2
	selected = model.getSelected()
	if selected != 109 {
		t.Errorf("Expected selection 109 for BLACK, got %d", selected)
	}
	
	model.Index = 3
	selected = model.getSelected()
	if selected != 107 {
		t.Errorf("Expected selection 107 for ODD, got %d", selected)
	}
}

func TestRouletteGetSelectedName(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	
	// Test main section names
	model.Row = 0
	model.Col = 0
	name := model.getSelectedName()
	if name != "#1" {
		t.Errorf("Expected name '#1', got '%s'", name)
	}
	
	// Test dozens section names
	model.Section = SectionDozens
	model.Index = 0
	name = model.getSelectedName()
	if name != "1-12" {
		t.Errorf("Expected name '1-12', got '%s'", name)
	}
	
	// Test low/high section names
	model.Section = SectionLowHigh
	model.Index = 0
	name = model.getSelectedName()
	if name != "#0" {
		t.Errorf("Expected name '#0', got '%s'", name)
	}
	
	model.Index = 1
	name = model.getSelectedName()
	if name != "1-18" {
		t.Errorf("Expected name '1-18', got '%s'", name)
	}
	
	// Test outside section names
	model.Section = SectionOutside
	model.Index = 0
	name = model.getSelectedName()
	if name != "EVEN" {
		t.Errorf("Expected name 'EVEN', got '%s'", name)
	}
}

func TestRouletteZeroBettingInLowHigh(t *testing.T) {
	wallet := &MockWallet{Balance: 1000}
	model := NewModel(wallet)
	model.Section = SectionLowHigh
	model.Index = 0 // 0 position
	
	selected := model.getSelected()
	if selected != 0 {
		t.Errorf("Expected to select 0 in low/high section, got %d", selected)
	}
	
	name := model.getSelectedName()
	if name != "#0" {
		t.Errorf("Expected name '#0' for zero selection, got '%s'", name)
	}
}
