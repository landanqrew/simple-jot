/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "command to delete a note",
	Long: `command that takes in a noteId and deletes the note from storage:

examples:
  simple-jot delete 1234567890
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("note id is required")
		}

		// get notes
		noteSlice, err := storage.GetNotes()
		if err != nil {
			log.Fatal("failed to get notes: " + err.Error())
		}

		// find note based on id
		noteId := args[0]
		newNotes := make([]notes.Note,len(noteSlice))
		found := false
		for _, note := range noteSlice {
			if note.ID != noteId {
				newNotes = append(newNotes, note)
			} else {
				found = true
			}
		}
		if !found {
			log.Fatal("note not found")
		}

		// save notes
		err = storage.SaveNotes(newNotes)
		if err != nil {
			log.Fatal("failed to save notes: " + err.Error())
		}
		fmt.Println("note deleted: " + noteId)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	
}
