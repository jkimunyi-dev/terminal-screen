package renderer

import (
	"fmt"
)

// DrawingCommandHandler encapsulates the drawing command logic
type DrawingCommandHandler struct {
	renderer *TerminalRenderer
}

// NewDrawingCommandHandler creates a new drawing command handler
func NewDrawingCommandHandler(renderer *TerminalRenderer) *DrawingCommandHandler {
	return &DrawingCommandHandler{
		renderer: renderer,
	}
}

// HandleDrawCharacterCommand handles the draw character command (0x2)
func (h *DrawingCommandHandler) HandleDrawCharacterCommand(cmd *Command) error {
	// Validate command data length
	if err := h.validateDrawCharacterCommand(cmd); err != nil {
		return err
	}

	// Extract command parameters
	x, y, colorIndex, char := h.extractDrawCharacterParams(cmd)

	// Get current screen
	screen := h.renderer.GetCurrentScreen()
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

// validateDrawCharacterCommand ensures the command data is valid
func (h *DrawingCommandHandler) validateDrawCharacterCommand(cmd *Command) error {
	if len(cmd.Data) != 4 {
		return fmt.Errorf("invalid draw character command length: expected 4, got %d", len(cmd.Data))
	}
	return nil
}

// extractDrawCharacterParams extracts parameters from the draw character command
func (h *DrawingCommandHandler) extractDrawCharacterParams(cmd *Command) (x, y, colorIndex uint8, char rune) {
	x = cmd.Data[0]
	y = cmd.Data[1]
	colorIndex = cmd.Data[2]
	char = rune(cmd.Data[3])
	return
}

// HandleDrawLineCommand handles the draw line command (0x3)
func (h *DrawingCommandHandler) HandleDrawLineCommand(cmd *Command) error {
	// Validate command data length
	if err := h.validateDrawLineCommand(cmd); err != nil {
		return err
	}

	// Extract command parameters
	x1, y1, x2, y2, colorIndex, lineChar := h.extractDrawLineParams(cmd)

	// Get current screen
	screen := h.renderer.GetCurrentScreen()
	if screen == nil {
		return fmt.Errorf("no screen initialized")
	}

	// Draw line using Bresenham's algorithm
	return h.drawLine(screen, x1, y1, x2, y2, colorIndex, lineChar)
}

// validateDrawLineCommand ensures the command data is valid
func (h *DrawingCommandHandler) validateDrawLineCommand(cmd *Command) error {
	if len(cmd.Data) != 6 {
		return fmt.Errorf("invalid draw line command length: expected 6, got %d", len(cmd.Data))
	}
	return nil
}

// extractDrawLineParams extracts parameters from the draw line command
func (h *DrawingCommandHandler) extractDrawLineParams(cmd *Command) (x1, y1, x2, y2, colorIndex uint8, lineChar rune) {
	x1 = cmd.Data[0]
	y1 = cmd.Data[1]
	x2 = cmd.Data[2]
	y2 = cmd.Data[3]
	colorIndex = cmd.Data[4]
	lineChar = rune(cmd.Data[5])
	return
}

// HandleRenderTextCommand handles the render text command (0x4)
func (h *DrawingCommandHandler) HandleRenderTextCommand(cmd *Command) error {
	// Validate command data length
	if err := h.validateRenderTextCommand(cmd); err != nil {
		return err
	}

	// Extract command parameters
	x, y, colorIndex, text := h.extractRenderTextParams(cmd)

	// Get current screen
	screen := h.renderer.GetCurrentScreen()
	if screen == nil {
		return fmt.Errorf("no screen initialized")
	}

	// Render text character by character
	return h.renderText(screen, x, y, colorIndex, text)
}

// validateRenderTextCommand ensures the command data is valid
func (h *DrawingCommandHandler) validateRenderTextCommand(cmd *Command) error {
	if len(cmd.Data) < 3 {
		return fmt.Errorf("invalid render text command length: at least 3 bytes required, got %d", len(cmd.Data))
	}
	return nil
}

// extractRenderTextParams extracts parameters from the render text command
func (h *DrawingCommandHandler) extractRenderTextParams(cmd *Command) (x, y, colorIndex uint8, text []byte) {
	x = cmd.Data[0]
	y = cmd.Data[1]
	colorIndex = cmd.Data[2]
	text = cmd.Data[3:]
	return
}

// renderText renders text on the screen
func (h *DrawingCommandHandler) renderText(screen *Screen, x, y, colorIndex uint8, text []byte) error {
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

// drawLine uses Bresenham's line algorithm to draw a line on the screen
func (h *DrawingCommandHandler) drawLine(screen *Screen, x1, y1, x2, y2 uint8, colorIndex uint8, lineChar rune) error {
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

// Utility helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}
