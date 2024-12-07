package renderer

// CommandType represents the different types of terminal rendering commands
type CommandType uint8

// ColorMode defines the available color modes for the terminal
type ColorMode uint8

const (
	// Screen Setup Command
	CommandScreenSetup CommandType = 0x1

	// Drawing Commands
	CommandDrawCharacter CommandType = 0x2
	CommandDrawLine      CommandType = 0x3
	CommandRenderText    CommandType = 0x4

	// Cursor and Rendering Control
	CommandMoveCursor   CommandType = 0x5
	CommandDrawAtCursor CommandType = 0x6
	CommandClearScreen  CommandType = 0x7

	// End of Stream
	CommandEndOfStream CommandType = 0xFF
)
