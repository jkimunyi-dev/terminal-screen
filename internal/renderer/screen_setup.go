package renderer

import (
	"fmt"
	"sync"
)

// ScreenSetupOptions provides additional configuration for screen initialization
type ScreenSetupOptions struct {
	// BackgroundColor sets the initial background color for the entire screen
	BackgroundColor uint8
	// InitialFillCharacter sets the character used to fill the initial screen
	InitialFillCharacter rune
}

// DefaultScreenSetupOptions provides sensible defaults for screen initialization
func DefaultScreenSetupOptions() *ScreenSetupOptions {
	return &ScreenSetupOptions{
		BackgroundColor:      0,   // Black background
		InitialFillCharacter: ' ', // Space character
	}
}

// ScreenManager manages screen creation and lifecycle
type ScreenManager struct {
	currentScreen *Screen
	mutex         sync.Mutex
}

// NewScreenManager creates a new screen manager
func NewScreenManager() *ScreenManager {
	return &ScreenManager{}
}

// HandleScreenSetupCommand processes the screen setup command
func (sm *ScreenManager) HandleScreenSetupCommand(data []byte, options *ScreenSetupOptions) (*Screen, error) {
	// Validate input data
	if len(data) < 3 {
		return nil, fmt.Errorf("insufficient data for screen setup: need at least 3 bytes, got %d", len(data))
	}

	// Extract screen dimensions and color mode
	width := data[0]
	height := data[1]
	colorMode := ColorMode(data[2])

	// Validate dimensions
	if width == 0 || height == 0 {
		return nil, fmt.Errorf("invalid screen dimensions: width=%d, height=%d", width, height)
	}

	// Validate color mode
	switch colorMode {
	case ColorModeMonochrome, ColorMode16, ColorMode256:
		// Valid color modes
	default:
		return nil, fmt.Errorf("unsupported color mode: %d", colorMode)
	}

	// Apply default options if not provided
	if options == nil {
		options = DefaultScreenSetupOptions()
	}

	// Create new screen
	screen := NewScreen(width, height, colorMode)

	// Initialize screen with background and fill character
	for y := uint8(0); y < height; y++ {
		for x := uint8(0); x < width; x++ {
			err := screen.SetCell(x, y, Cell{
				Char:      options.InitialFillCharacter,
				FgColor:   7, // Default light gray foreground
				BgColor:   options.BackgroundColor,
				Highlight: false,
			})
			if err != nil {
				return nil, fmt.Errorf("error initializing screen cell at (%d, %d): %v", x, y, err)
			}
		}
	}

	// Safely update the current screen
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.currentScreen = screen

	return screen, nil
}

// GetCurrentScreen returns the currently active screen
func (sm *ScreenManager) GetCurrentScreen() *Screen {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return sm.currentScreen
}

// ValidateScreenSetupCommand provides additional validation for the screen setup command
func ValidateScreenSetupCommand(data []byte) error {
	if len(data) < 3 {
		return fmt.Errorf("screen setup command requires at least 3 bytes")
	}

	width := data[0]
	height := data[1]
	colorMode := data[2]

	// Validate width and height
	if width == 0 || width > 255 {
		return fmt.Errorf("invalid screen width: %d (must be between 1 and 255)", width)
	}

	if height == 0 || height > 255 {
		return fmt.Errorf("invalid screen height: %d (must be between 1 and 255)", height)
	}

	// Validate color mode
	switch colorMode {
	case 0x00, 0x01, 0x02:
		return nil
	default:
		return fmt.Errorf("invalid color mode: %d (must be 0x00, 0x01, or 0x02)", colorMode)
	}
}

// ColorModeToString provides a human-readable representation of color modes
func ColorModeToString(mode ColorMode) string {
	switch mode {
	case ColorModeMonochrome:
		return "Monochrome"
	case ColorMode16:
		return "16 Colors"
	case ColorMode256:
		return "256 Colors"
	default:
		return "Unknown"
	}
}
