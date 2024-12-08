package renderer

import "fmt"

// handleDrawCharacterCommand handles the draw character command (0x2)
func (tr *TerminalRenderer) handleDrawCharacterCommand(cmd *Command) error {
	// Validate command data length
	if len(cmd.Data) != 4 {
		return fmt.Errorf("invalid draw character command length: expected 4, got %d", len(cmd.Data))
	}

	// Extract coordinates and color
	x, y := cmd.Data[0], cmd.Data[1]
	colorIndex := cmd.Data[2]
	char := rune(cmd.Data[3])

	// Get current screen
	screen := tr.GetCurrentScreen()
	if screen == nil {
		return fmt.Errorf("no screen initialized")
	}

	// Create cell with character and color
	cell := Cell{
		Char:    char,
		FgColor: colorIndex,
	}

	// Set the cell in the screen buffer
	return screen.SetCell(x, y, cell)
}

// handleDrawLineCommand handles the draw line command (0x3)
func (tr *TerminalRenderer) handleDrawLineCommand(cmd *Command) error {
	// Validate command data length
	if len(cmd.Data) != 6 {
		return fmt.Errorf("invalid draw line command length: expected 6, got %d", len(cmd.Data))
	}

	// Extract coordinates
	x1, y1 := cmd.Data[0], cmd.Data[1]
	x2, y2 := cmd.Data[2], cmd.Data[3]
	colorIndex := cmd.Data[4]
	lineChar := rune(cmd.Data[5])

	// Get current screen
	screen := tr.GetCurrentScreen()
	if screen == nil {
		return fmt.Errorf("no screen initialized")
	}

	// Implement Bresenham's line drawing algorithm
	return tr.drawLine(screen, x1, y1, x2, y2, colorIndex, lineChar)
}
