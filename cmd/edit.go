/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/osutils"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/landanqrew/simple-jot/tabler"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit <note-id>",
	Short: "edit a note",
	Long: `edit the contents of a note:

Usage:
To overwrite the note:
simple-jot edit <note-id> -n '<note-content>'
cat <my-file.txt> | simple-jot edit <note-id>

To append to the note:
simple-jot edit <note-id> -a '<note-content>'`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteID := args[0]

		// Get flag values
		noteContent, _ := cmd.Flags().GetString("note")
		appendContent, _ := cmd.Flags().GetString("append")

		// check if stdin is from a pipe or redirect
		stdinContent, err := osutils.ReadStdin()
		if err != nil {
			return fmt.Errorf("cannot read from stdin: %v", err)
		}

		if stdinContent != "" {
			// If content from stdin, prioritize it
			if noteContent != "" {
				return fmt.Errorf("cannot provide note content via both -n flag and stdin. Please choose one")
			}
			if appendContent != "" {
				return fmt.Errorf("cannot use stdin to append to a note")
			}
			noteContent = stdinContent
		}

		// handle invalid args
		if noteContent == "" && appendContent == "" {
			return fmt.Errorf("note content cannot be empty. Please use the -n or --note flag to provide content, or pipe content to stdin")
		} else if noteContent != "" && appendContent != "" {
			return fmt.Errorf("cannot use both -n and -a flags. Please use only one")
		}

		// fetch notes
		noteList, err := storage.GetNotes()
		if err != nil {
			return fmt.Errorf("cannot fetch notes: %v", err)
		}

		// update note
		found := false
		currentNote := notes.Note{}
		for i, n := range noteList {
			if n.ID == noteID {
				if appendContent != "" {
					noteList[i].Content += appendContent
				} else {
					noteList[i].Content = noteContent
				}
				noteList[i].UpdatedAt = time.Now().Format(time.DateTime)
				found = true
				currentNote = noteList[i]
				break
			}
		}

		if !found {
			return fmt.Errorf("note with ID '%s' not found", noteID)
		}

		// save notes
		err = storage.SaveNotes(noteList)
		if err != nil {
			return fmt.Errorf("cannot save notes: %v", err)
		}
		cmd.Println("Note updated successfully.")

		err = tabler.RenderTable([][]string{currentNote.PrepRow()}, currentNote.GetHeaders())
		if err != nil {
			return fmt.Errorf("failed to render table: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringP("note", "n", "", "Content of the note. If not provided, content will be read from stdin.")
	editCmd.Flags().StringP("append", "a", "", "Append this content to the active configuration note")
}
