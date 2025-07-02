/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit a note",
	Long: `edit the contents of a note:

Usage:
To overwrite the note:
simple-jot edit <note-id> -n "<note-content>"

To append to the note:
simple-jot edit <note-id> -a "<note-content>"`,
	Run: func(cmd *cobra.Command, args []string) {
		// parse args
		if len(args) < 1 {
			fmt.Println("Error: Note ID is required. Please provide a note ID.")
			os.Exit(1)
		}
		noteID := args[0]
		noteContent := ""
		appendContent := ""
		for i, arg := range args {
			if i == 0 {
				continue
			}
			if arg == "-n" {
				noteContent = args[i + 1]
			}
			if arg == "-a" {
				appendContent = args[i + 1]
			}
		}
		// handle invalid args
		if noteContent == "" && appendContent == "" {
			fmt.Println("Error: Note content cannot be empty. Please use the -n or --note flag to provide content.")
			os.Exit(1)
		} else if noteContent != "" && appendContent != "" {
			fmt.Println("Error: Cannot use both -n and -a flags. Please use only one.")
			os.Exit(1)
		}
		// fetch notes
		noteList, err := storage.GetNotes()
		if err != nil {
			fmt.Println("Error: Cannot fetch notes. See error: " + err.Error())
			os.Exit(1)
		}
		// update note
		for i, n := range noteList {
			if n.ID == noteID {
				if appendContent != "" {
					noteList[i].Content += appendContent
				} else {
					noteList[i].Content = noteContent
				}
				noteList[i].UpdatedAt = time.Now().Format(time.DateTime)
				break
			}
		}

		// save notes
		err = storage.SaveNotes(noteList)
		if err != nil {
			fmt.Println("Error: Cannot save notes. See error: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("Note updated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
