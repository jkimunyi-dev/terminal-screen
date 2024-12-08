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

	// Create a new terminal renderer
	termRenderer := renderer.NewTerminalRenderer()

	// Prepare a screen setup command
	cmd := &renderer.Command{
		Type:   renderer.CommandScreenSetup,
		Length: 3,
		Data:   []byte{80, 24, 0x01}, // 80x24 16-color screen
	}

	// Handle the command
	err = termRenderer.HandleCommand(cmd)
	if err != nil {
		log.Fatalf("Screen setup failed: %v", err)
	}
}
