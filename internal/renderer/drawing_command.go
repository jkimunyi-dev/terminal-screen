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

// drawLine uses Bresenham's line algorithm to draw a line on the screen
func (tr *TerminalRenderer) drawLine(screen *Screen, x1, y1, x2, y2 uint8, colorIndex uint8, lineChar rune) error {
	steep := false
	if abs(int(x1)-int(x2)) < abs(int(y1)-int(y2)) {
		// Swap x and y coordinates to make the line more horizontal
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		steep = true
	}

	// Ensure line is drawn left to right
	if int(x1) > int(x2) {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	dx := int(x2) - int(x1)
	dy := int(y2) - int(y1)
	derror2 := abs(dy) * 2
	error2 := 0
	y := int(y1)

	cell := Cell{
		Char:    lineChar,
		FgColor: colorIndex,
	}

	for x := int(x1); x <= int(x2); x++ {
		var setX, setY uint8
		if steep {
			setX, setY = uint8(y), uint8(x)
		} else {
			setX, setY = uint8(x), uint8(y)
		}

		err := screen.SetCell(setX, setY, cell)
		if err != nil {
			return err
		}

		error2 += derror2
		if error2 > dx {
			y += sign(dy)
			error2 -= dx * 2
		}
	}

	return nil
}

// handleRenderTextCommand handles the render text command (0x4)
func (tr *TerminalRenderer) handleRenderTextCommand(cmd *Command) error {
	// Validate minimum command data length
	if len(cmd.Data) < 3 {
		return fmt.Errorf("invalid render text command length: at least 3 bytes required, got %d", len(cmd.Data))
	}

	// Extract coordinates and color
	x, y := cmd.Data[0], cmd.Data[1]
	colorIndex := cmd.Data[2]

	// Extract text (remaining bytes)
	text := cmd.Data[3:]

	// Get current screen
	screen := tr.GetCurrentScreen()
	if screen == nil {
		return fmt.Errorf("no screen initialized")
	}

	// Render text character by character
	for i, b := range text {
		char := rune(b)
		cellX := x + uint8(i)

		cell := Cell{
			Char:    char,
			FgColor: colorIndex,
		}

		err := screen.SetCell(cellX, y, cell)
		if err != nil {
			return fmt.Errorf("error rendering text at (%d, %d): %v", cellX, y, err)
		}
	}

	return nil
}
