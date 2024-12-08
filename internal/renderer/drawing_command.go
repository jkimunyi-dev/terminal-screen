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
