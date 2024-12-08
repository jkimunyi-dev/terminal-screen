package main

import (
	"fmt"
	"log"

	"github.com/jkimunyi-dev/termial-screen/internal/renderer"
)

func main() {
	// Create a new terminal renderer
	termRenderer := renderer.NewTerminalRenderer()

	// Prepare a screen setup command
	cmd := &renderer.Command{
		Type:   renderer.CommandScreenSetup,
		Length: 3,
		Data:   []byte{80, 24, 0x01}, // 80x24 16-color screen
	}

	// Handle the command
	err := termRenderer.HandleCommand(cmd)
	if err != nil {
		log.Fatalf("Screen setup failed: %v", err)
	}

	// Retrieve the current screen
	screen := termRenderer.GetCurrentScreen()
	if screen == nil {
		log.Fatal("No screen created")
	}

	// Print out screen details
	fmt.Printf("Screen Created:\n")
	fmt.Printf("Width: %d\n", screen.GetWidth())
	fmt.Printf("Height: %d\n", screen.GetHeight())
	fmt.Printf("Color Mode: %s\n", renderer.ColorModeToString(screen.GetColorMode()))

	// Example of setting and getting a cell
	err = screen.SetCell(10, 5, renderer.Cell{
		Char:    'H',
		FgColor: 1, // Red foreground
		BgColor: 0, // Black background
	})
	if err != nil {
		log.Fatalf("Error setting cell: %v", err)
	}

	// Retrieve and print the cell we just set
	cell, err := screen.GetCell(10, 5)
	if err != nil {
		log.Fatalf("Error getting cell: %v", err)
	}
	fmt.Printf("\nCell at (10,5):\n")
	fmt.Printf("Character: %c\n", cell.Char)
	fmt.Printf("Foreground Color: %d\n", cell.FgColor)
	fmt.Printf("Background Color: %d\n", cell.BgColor)
}
