package main

import (
	"fmt"
	"log"

	"github.com/jkimunyi-dev/termial-screen/internal/renderer"
)

func main() {
	// Example of creating a screen
	screen := renderer.NewScreen(80, 24, renderer.ColorMode16)

	// Example of setting a cell
	err := screen.SetCell(10, 5, renderer.Cell{
		Char:    'H',
		FgColor: 1, // Red foreground
		BgColor: 0, // Black background
	})

	if err != nil {
		log.Fatalf("Error setting cell: %v", err)
	}

	fmt.Println("Terminal Screen Renderer initialized successfully!")

}
