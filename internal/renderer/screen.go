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

type Cell struct {
	Char      rune
	FgColor   uint8
	BgColor   uint8
	Highlight bool
}
