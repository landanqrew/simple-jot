/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/osutils"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/landanqrew/simple-jot/tabler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <note-name>",
	Short: "Creates a new note with specified content and optional configuration.",
	Long: `The create command allows you to generate a new note by providing a name and its content.
Optionally, you can set this new note as your active configuration note using the -s flag.

Examples:
  simple-jot create my-first-note -n 'This is the content of my first note.'
  simple-jot create daily-log -n '2006-01-01 entry' -s
  cat some_file.txt | simple-jot create my-piped-note
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteName := args[0]

		// Get flag values
		noteContent, _ := cmd.Flags().GetString("note")
		setNote, _ := cmd.Flags().GetBool("set")

		// Check if content is provided via stdin
		if stdinContent, err := osutils.ReadStdin(); err != nil {
			return fmt.Errorf("error reading from stdin: %v", err)
		} else if stdinContent != "" {
			// If content from stdin, prioritize it over flag content if both are present
			if noteContent == "" {
				noteContent = stdinContent
			} else {
				return fmt.Errorf("cannot provide note content via both -n flag and stdin. Please choose one")
			}
		}

		// Basic validation for note content
		if noteContent == "" {
			return fmt.Errorf("note content cannot be empty. Please use the -n or --note flag to provide content, or pipe content to stdin")
		}

		cmd.Printf("Creating note: %s\n", noteName)
		cmd.Printf("Note content: %s\n", noteContent)
		cmd.Printf("Set as current note: %t\n", setNote)

		noteSlice, err := storage.GetNotes()
		if err != nil {
			return fmt.Errorf("failed to get notes: %v", err)
		}
		newNote := notes.Note{
			ID:        uuid.New().String(),
			Title:     noteName,
			Content:   noteContent,
			CreatedAt: time.Now().Format(time.DateTime),
			UpdatedAt: time.Now().Format(time.DateTime),
			Tags:      []string{},
		}
		noteSlice = append(noteSlice, newNote)
		err = storage.SaveNotes(noteSlice)
		if err != nil {
			return fmt.Errorf("failed to save notes: %v", err)
		}

		if setNote {
			viper.Set("active_note", newNote.ID)
			err = viper.WriteConfig()
			if err != nil {
				return fmt.Errorf("failed to save config: %v", err)
			}
			cmd.Printf("Active note set to: %s\n", newNote.ID)
		} else {
			cmd.Printf("Note created successfully. Use 'simple-jot config set note %s' to set this note as the active note.\n", newNote.ID)
		}

		err = tabler.RenderTable([][]string{newNote.PrepRow()}, newNote.GetHeaders())
		if err != nil {
			return fmt.Errorf("failed to render table: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Define flags for the create command
	createCmd.Flags().StringP("note", "n", "", "Content of the note. If not provided, content will be read from stdin.")
	createCmd.Flags().BoolP("set", "s", false, "Set this note as the active configuration note")
}
