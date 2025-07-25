package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/spf13/cobra"
)

func TestCreateCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		stdinContent   string
		expectedError  bool
		expectedOutput string
		mockStorage    *mockStorage
	}{
		{
			name:           "Create note with -n flag",
			args:           []string{"Test Note", "-n", "Test content"},
			stdinContent:   "",
			expectedError:  false,
			expectedOutput: "Creating note: Test Note\nNote content: Test content",
			mockStorage:    &mockStorage{notes: []notes.Note{}},
		},
		{
			name:           "Create note with stdin",
			args:           []string{"Test Note"},
			stdinContent:   "Test content from stdin",
			expectedError:  false,
			expectedOutput: "Creating note: Test Note\nNote content: Test content from stdin",
			mockStorage:    &mockStorage{notes: []notes.Note{}},
		},
		{
			name:           "Create note with both stdin and -n flag",
			args:           []string{"Test Note", "-n", "Flag content"},
			stdinContent:   "Stdin content",
			expectedError:  true,
			expectedOutput: "cannot provide note content via both -n flag and stdin",
			mockStorage:    &mockStorage{notes: []notes.Note{}},
		},
		{
			name:           "Create note without content",
			args:           []string{"Test Note"},
			stdinContent:   "",
			expectedError:  true,
			expectedOutput: "note content cannot be empty",
			mockStorage:    &mockStorage{notes: []notes.Note{}},
		},
		{
			name:           "Storage error",
			args:           []string{"Test Note", "-n", "Test content"},
			stdinContent:   "",
			expectedError:  true,
			expectedOutput: "failed to save notes: mock storage error",
			mockStorage:    &mockStorage{notes: []notes.Note{}, saveError: fmt.Errorf("mock storage error")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock storage
			storage.SetDefaultStorage(tt.mockStorage)

			// Save original stdin and restore it after the test
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// If we have stdin content, create a pipe and write to it
			if tt.stdinContent != "" {
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatalf("Failed to create pipe: %v", err)
				}
				os.Stdin = r
				go func() {
					defer w.Close()
					io.WriteString(w, tt.stdinContent)
				}()
			} else {
				// If no stdin content, use /dev/null
				devNull, err := os.Open(os.DevNull)
				if err != nil {
					t.Fatalf("Failed to open /dev/null: %v", err)
				}
				defer devNull.Close()
				os.Stdin = devNull
			}

			// Create a new command instance
			cmd := cobra.Command{
				Use: "create",
				RunE: func(cmd *cobra.Command, args []string) error {
					return createCmd.RunE(cmd, args)
				},
			}

			// Add flags
			cmd.Flags().StringP("note", "n", "", "Content of the note")
			cmd.Flags().BoolP("set", "s", false, "Set as active note")

			// Set up output capture
			output := new(bytes.Buffer)
			cmd.SetOut(output)
			cmd.SetErr(output)

			// Set args and execute
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			// Check error
			if tt.expectedError {
				if err == nil {
					t.Error("Expected an error but got none")
				} else if !strings.Contains(err.Error(), tt.expectedOutput) {
					t.Errorf("Expected error containing %q, got %q", tt.expectedOutput, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				got := output.String()
				if !strings.Contains(got, tt.expectedOutput) {
					t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, got)
				}
			}
		})
	}
}

func TestCreateCmdValidation(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name:          "No note name",
			args:          []string{},
			expectedError: "accepts 1 arg(s), received 0",
		},
		{
			name:          "Too many arguments",
			args:          []string{"Note1", "Note2"},
			expectedError: "accepts 1 arg(s), received 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := createCmd.Args(cmd, tt.args)
			if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing %q, got %v", tt.expectedError, err)
			}
		})
	}
}
