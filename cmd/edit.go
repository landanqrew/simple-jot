/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/osutils"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/landanqrew/simple-jot/tabler"
	"github.com/spf13/cobra"
)

var (
	// noteContent   string -- already defined in create.go
	appendContent string
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit <note-id>",
	Short: "edit a note",
	Long: `edit the contents of a note:

Usage:
To overwrite the note:
simple-jot edit <note-id> -n "<note-content>"
cat <my-file.txt> | simple-jot edit <note-id>

To append to the note:
simple-jot edit <note-id> -a "<note-content>"`,
	Run: func(cmd *cobra.Command, args []string) {
		// parse args
		if len(args) < 1 {
			fmt.Println("Error: Note ID is required. Please provide a note ID.")
			os.Exit(1)
		}
		noteID := args[0]

		// check if stdin is from a pipe or redirect
		stdinContent, err := osutils.ReadStdin()
		if err != nil {
			fmt.Println("Error: Cannot read from stdin. See error: " + err.Error())
			os.Exit(1)
		}

		if stdinContent != "" {
			// If content from stdin, prioritize it
			if noteContent != "" {
				log.Fatal("Error: Cannot provide note content via both -n flag and stdin. Please choose one.")
			}
			if appendContent != "" {
				log.Fatal("Error: Cannot use stdin to append to a note.")
			}
			noteContent = stdinContent
		}

		// handle invalid args
		if noteContent == "" && appendContent == "" {
			log.Fatal("Error: Note content cannot be empty. Please use the -n or --note flag to provide content, or pipe content to stdin.")
		} else if noteContent != "" && appendContent != "" {
			log.Fatal("Error: Cannot use both -n and -a flags. Please use only one.")
		}

		// fetch notes
		noteList, err := storage.GetNotes()
		if err != nil {
			fmt.Println("Error: Cannot fetch notes. See error: " + err.Error())
			os.Exit(1)
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
			log.Fatal("Error: Note with ID '" + noteID + "' not found.")
		}

		// save notes
		err = storage.SaveNotes(noteList)
		if err != nil {
			log.Fatal("Error: Cannot save notes. See error: " + err.Error())
		}
		fmt.Println("Note updated successfully.")

		tabler.RenderTable([][]string{currentNote.PrepRow()}, currentNote.GetHeaders())
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	editCmd.Flags().StringVarP(&noteContent, "note", "n", "", "Content of the note. If not provided, content will be read from stdin.")
	editCmd.Flags().StringVarP(&appendContent, "append", "a", "", "Append this content to the active configuration note")

	editCmd.Flags().SetAnnotation("note", cobra.BashCompFilenameExt, []string{"txt"})

	// editCmd.MarkFlagRequired("note") // Commented out to allow stdin input
}
