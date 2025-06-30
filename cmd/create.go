/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os" // Import os for exiting on error

	"github.com/spf13/cobra"
)

var (
	noteContent string
	setNote     bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <note-name>",
	Short: "Creates a new note with specified content and optional configuration.",
	Long: `The create command allows you to generate a new note by providing a name and its content.
Optionally, you can set this new note as your active configuration note using the -s flag.

Examples:
  simple-jot create my-first-note -n "This is the content of my first note."
  simple-jot create daily-log -n "Today's entry." -s
`,
	Args: cobra.ExactArgs(1), // Ensures exactly one positional argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		noteName := args[0]

		// Basic validation for note content
		if noteContent == "" {
			fmt.Println("Error: Note content cannot be empty. Please use the -n or --note flag to provide content.")
			os.Exit(1)
		}

		fmt.Printf("Creating note: %s\n", noteName)
		fmt.Printf("Note content: %s\n", noteContent)
		fmt.Printf("Set as current note: %t\n", setNote)

		// TODO: Add logic here to save the note and handle the 'setNote' flag

		fmt.Println("Note created successfully (placeholder).")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Define flags for the create command
	createCmd.Flags().StringVarP(&noteContent, "note", "n", "", "Content of the note")
	createCmd.Flags().BoolVarP(&setNote, "set", "s", false, "Set this note as the active configuration note")

	// Mark the -n flag as required, meaning it must be provided.
	createCmd.MarkFlagRequired("note")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
