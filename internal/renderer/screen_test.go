package renderer

import (
	"testing"
)

func TestScreenSetupCommand(t *testing.T) {
	// Test cases for screen setup
	testCases := []struct {
		name           string
		inputData      []byte
		expectedError  bool
		expectedWidth  uint8
		expectedHeight uint8
		expectedColor  ColorMode
	}{
		{
			name:           "Valid 80x24 16-color screen",
			inputData:      []byte{80, 24, 0x01},
			expectedError:  false,
			expectedWidth:  80,
			expectedHeight: 24,
			expectedColor:  ColorMode16,
		},
		{
			name:           "Valid Monochrome Small Screen",
			inputData:      []byte{40, 10, 0x00},
			expectedError:  false,
			expectedWidth:  40,
			expectedHeight: 10,
			expectedColor:  ColorModeMonochrome,
		},
		{
			name:          "Invalid Zero Width",
			inputData:     []byte{0, 24, 0x01},
			expectedError: true,
		},
		{
			name:          "Invalid Zero Height",
			inputData:     []byte{80, 0, 0x01},
			expectedError: true,
		},
		{
			name:          "Invalid Color Mode",
			inputData:     []byte{80, 24, 0x03},
			expectedError: true,
		},
		{
			name:          "Insufficient Data",
			inputData:     []byte{80},
			expectedError: true,
		},
	}

	// Create a screen manager for testing
	screenManager := NewScreenManager()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validate the command first
			err := ValidateScreenSetupCommand(tc.inputData)
			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected validation error: %v", err)
				return
			}

			// Attempt to set up the screen
			screen, err := screenManager.HandleScreenSetupCommand(tc.inputData, nil)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected an error, but screen setup succeeded")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error during screen setup: %v", err)
				return
			}

			// Verify screen properties
			if screen.width != tc.expectedWidth {
				t.Errorf("Incorrect screen width. Expected %d, got %d", tc.expectedWidth, screen.width)
			}

			if screen.height != tc.expectedHeight {
				t.Errorf("Incorrect screen height. Expected %d, got %d", tc.expectedHeight, screen.height)
			}

			if screen.colorMode != tc.expectedColor {
				t.Errorf("Incorrect color mode. Expected %d, got %d", tc.expectedColor, screen.colorMode)
			}

			// Verify screen initialization
			cell, err := screen.GetCell(0, 0)
			if err != nil {
				t.Errorf("Failed to get initial cell: %v", err)
			}

			if cell.Char != ' ' {
				t.Errorf("Initial cell not properly initialized. Expected space, got %c", cell.Char)
			}
		})
	}
}

func TestColorModeToString(t *testing.T) {
	testCases := []struct {
		mode     ColorMode
		expected string
	}{
		{ColorModeMonochrome, "Monochrome"},
		{ColorMode16, "16 Colors"},
		{ColorMode256, "256 Colors"},
		{ColorMode(0x04), "Unknown"},
	}

	for _, tc := range testCases {
		result := ColorModeToString(tc.mode)
		if result != tc.expected {
			t.Errorf("Unexpected color mode string. For mode %d, expected %s, got %s",
				tc.mode, tc.expected, result)
		}
	}
}
