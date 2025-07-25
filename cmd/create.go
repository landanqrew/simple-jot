/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os" // Import os for exiting on error

	"time"

	"github.com/google/uuid"
	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/osutils"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/landanqrew/simple-jot/tabler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
  cat some_file.txt | simple-jot create
`,
	Args: cobra.ExactArgs(1), // Ensures exactly one positional argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("note name/title is required")
		}
		noteName := args[0]

		// Check if content is provided via stdin
		if stdinContent, err := osutils.ReadStdin(); err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		} else if stdinContent != "" {
			// If content from stdin, prioritize it over flag content if both are present
			if noteContent == "" {
				noteContent = stdinContent
			} else {
				// Or, decide how to handle the conflict, e.g., error out or concatenate
				log.Fatal("Error: Cannot provide note content via both -n flag and stdin. Please choose one.")
			}
		}

		// Basic validation for note content
		if noteContent == "" {
			log.Fatal("Error: Note content cannot be empty. Please use the -n or --note flag to provide content, or pipe content to stdin.")
		}

		fmt.Printf("Creating note: %s\n", noteName)
		fmt.Printf("Note content: %s\n", noteContent)
		fmt.Printf("Set as current note: %t\n", setNote)

		noteSlice, err := storage.GetNotes()
		if err != nil {
			log.Fatal("failed to get notes: " + err.Error())
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
			log.Fatal("failed to save notes: " + err.Error())
		}

		if setNote {
			viper.Set("active_note", newNote.ID)
			err := viper.WriteConfig()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error saving configuration:", err)
				activeNote := viper.GetString("active_note")
				fmt.Println("Active note is currently set to: " + activeNote)
				os.Exit(1)
			}
			fmt.Printf("Active note set to: %s\n", newNote.ID)
		} else {
			fmt.Printf("Note created successfully. Use 'simple-jot config set note %s' to set this note as the active note.\n", newNote.ID)
		}

		err = tabler.RenderTable([][]string{newNote.PrepRow()}, newNote.GetHeaders())
		if err != nil {
			log.Fatal("failed to render table: " + err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Define flags for the create command
	createCmd.Flags().StringVarP(&noteContent, "note", "n", "", "Content of the note. If not provided, content will be read from stdin.")
	createCmd.Flags().BoolVarP(&setNote, "set", "s", false, "Set this note as the active configuration note")

	// The `note` flag is no longer strictly required if content can be piped.
	// Instead, the check is moved into the Run function.
	// createCmd.MarkFlagRequired("note") // This line is commented out to make the flag optional

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
