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

// GetWidth returns the screen width
func (s *Screen) GetWidth() uint8 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.width
}

// GetHeight returns the screen height
func (s *Screen) GetHeight() uint8 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.height
}

// GetColorMode returns the screen's color mode
func (s *Screen) GetColorMode() ColorMode {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.colorMode
}
