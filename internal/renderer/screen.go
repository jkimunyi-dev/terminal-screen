package renderer

import (
	"fmt"
	"sync"
)

// Screen represents the terminal screen buffer
type Screen struct {
	mutex     sync.RWMutex // Changed from sync.Mutex to sync.RWMutex
	width     uint8
	height    uint8
	colorMode ColorMode
	buffer    [][]Cell
}

// Cell represents a single cell in the terminal screen
type Cell struct {
	Char      rune
	FgColor   uint8
	BgColor   uint8
	Highlight bool
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

// NewScreen initializes a new screen with given dimensions and color mode
func NewScreen(width, height uint8, mode ColorMode) *Screen {
	screen := &Screen{
		width:     width,
		height:    height,
		colorMode: mode,
		buffer:    make([][]Cell, height),
	}

	for i := range screen.buffer {
		screen.buffer[i] = make([]Cell, width)
	}

	return screen
}

// SetCell updates a specific cell in the screen buffer
func (s *Screen) SetCell(x, y uint8, cell Cell) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if x >= s.width || y >= s.height {
		return fmt.Errorf("coordinates out of bounds: (%d, %d)", x, y)
	}

	s.buffer[y][x] = cell
	return nil
}

// GetCell retrieves a cell from the screen buffer
func (s *Screen) GetCell(x, y uint8) (Cell, error) {
	s.mutex.RLock()         // Now using RLock()
	defer s.mutex.RUnlock() // And RUnlock()

	if x >= s.width || y >= s.height {
		return Cell{}, fmt.Errorf("coordinates out of bounds: (%d, %d)", x, y)
	}

	return s.buffer[y][x], nil
}

// Clear resets the entire screen buffer
func (s *Screen) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for y := range s.buffer {
		for x := range s.buffer[y] {
			s.buffer[y][x] = Cell{}
		}
	}
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
