package renderer

import "sync"

// Screen represents the terminal screen buffer
type Screen struct {
	mutex     sync.Mutex
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
