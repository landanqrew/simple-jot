package osutils

import (
	"io"
	"os"
	"testing"
)

func TestReadStdin(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{
			name:        "No Stdin",
			input:       "",
			expectError: false,
			expected:    "",
		},
		{
			name:        "With Stdin",
			input:       "Hello, world!\nThis is a test.\n",
			expectError: false,
			expected:    "Hello, world!\nThis is a test.\n",
		},
		{
			name:        "With Stdin - Single Line",
			input:       "Single line content\n",
			expectError: false,
			expected:    "Single line content\n",
		},
		{
			name:        "With Stdin - Empty Input",
			input:       "\n", // Represents an empty line piped
			expectError: false,
			expected:    "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original stdin
			oldStdin := os.Stdin

			// Restore original stdin after the test run
			defer func() {
				os.Stdin = oldStdin
			}()

			if tt.input != "" {
				// Create a pipe for stdin simulation
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatalf("Failed to create pipe: %v", err)
				}
				os.Stdin = r // Redirect stdin to the read end of the pipe

				// Write test input to the write file of the pipe
				_, err = io.WriteString(w, tt.input)
				if err != nil {
					t.Fatalf("Failed to write to pipe: %v", err)
				}
				w.Close() // Close the write file
			}

			// Call the function under test
			content, err := ReadStdin()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				if content != tt.expected {
					t.Errorf("Expected content %q, got %q", tt.expected, content)
				}
			}
		})
	}
}
