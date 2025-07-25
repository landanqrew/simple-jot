package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/spf13/cobra"
)

func TestEditCmd(t *testing.T) {
	// Create a test note that exists in storage
	existingNote := notes.Note{
		ID:        "test-id",
		Title:     "Test Note",
		Content:   "Original content",
		CreatedAt: time.Now().Format(time.DateTime),
		UpdatedAt: time.Now().Format(time.DateTime),
		Tags:      []string{},
	}

	tests := []struct {
		name           string
		args           []string
		stdinContent   string
		expectedError  bool
		expectedOutput string
		mockStorage    *mockStorage
	}{
		{
			name:           "Edit note with -n flag",
			args:           []string{"test-id", "-n", "New content"},
			stdinContent:   "",
			expectedError:  false,
			expectedOutput: "Note updated successfully",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}},
		},
		{
			name:           "Edit note with stdin",
			args:           []string{"test-id"},
			stdinContent:   "New content from stdin",
			expectedError:  false,
			expectedOutput: "Note updated successfully",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}},
		},
		{
			name:           "Edit note with -a flag",
			args:           []string{"test-id", "-a", " Appended content"},
			stdinContent:   "",
			expectedError:  false,
			expectedOutput: "Note updated successfully",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}},
		},
		{
			name:           "Edit note with both stdin and -n flag",
			args:           []string{"test-id", "-n", "Flag content"},
			stdinContent:   "Stdin content",
			expectedError:  true,
			expectedOutput: "cannot provide note content via both -n flag and stdin",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}},
		},
		{
			name:           "Edit note with both -n and -a flags",
			args:           []string{"test-id", "-n", "New content", "-a", "Appended content"},
			stdinContent:   "",
			expectedError:  true,
			expectedOutput: "cannot use both -n and -a flags",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}},
		},
		{
			name:           "Edit note with stdin and -a flag",
			args:           []string{"test-id", "-a", "Appended content"},
			stdinContent:   "Stdin content",
			expectedError:  true,
			expectedOutput: "cannot use stdin to append to a note",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}},
		},
		{
			name:           "Edit note without content",
			args:           []string{"test-id"},
			stdinContent:   "",
			expectedError:  true,
			expectedOutput: "note content cannot be empty",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}},
		},
		{
			name:           "Edit non-existent note",
			args:           []string{"non-existent-id", "-n", "New content"},
			stdinContent:   "",
			expectedError:  true,
			expectedOutput: "note with ID 'non-existent-id' not found",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}},
		},
		{
			name:           "Storage error",
			args:           []string{"test-id", "-n", "New content"},
			stdinContent:   "",
			expectedError:  true,
			expectedOutput: "cannot save notes: mock storage error",
			mockStorage:    &mockStorage{notes: []notes.Note{existingNote}, saveError: fmt.Errorf("mock storage error")},
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
				Use: "edit",
				RunE: func(cmd *cobra.Command, args []string) error {
					return editCmd.RunE(cmd, args)
				},
			}

			// Add flags
			cmd.Flags().StringP("note", "n", "", "Content of the note")
			cmd.Flags().StringP("append", "a", "", "Append content to the note")

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

func TestEditCmdValidation(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name:          "No note ID",
			args:          []string{},
			expectedError: "accepts 1 arg(s), received 0",
		},
		{
			name:          "Too many arguments",
			args:          []string{"note1", "note2"},
			expectedError: "accepts 1 arg(s), received 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := editCmd.Args(cmd, tt.args)
			if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing %q, got %v", tt.expectedError, err)
			}
		})
	}
}
