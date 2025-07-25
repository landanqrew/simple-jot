/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/landanqrew/simple-jot/internal/ai"
	"github.com/landanqrew/simple-jot/internal/config"
	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/landanqrew/simple-jot/internal/tags"
	"github.com/landanqrew/simple-jot/tabler"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for notes",
	Long: `Search for notes using various criteria. You can search by content, tags, date range, or perform semantic search.

Examples:
  # Search by date range
  simple-jot search --date-start 2025-01-01 --date-end 2025-02-01
  
  # Search by tag
  simple-jot search --tag 'tag1,tag2'
  
  # Search by content (case-insensitive)
  simple-jot search --content 'your search term'
  
  # Semantic search with AI
  simple-jot search --semantic 'programming concepts'
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		noteList, err := storage.GetNotes()
		if err != nil {
			return fmt.Errorf("cannot fetch notes: %w", err)
		}

		// Get flag values
		semanticSearch, _ := cmd.Flags().GetString("semantic")
		contentSearch, _ := cmd.Flags().GetString("content")
		tagStr, _ := cmd.Flags().GetString("tag")
		dsStr, _ := cmd.Flags().GetString("date-start")
		deStr, _ := cmd.Flags().GetString("date-end")

		filteredNotes := noteList

		if semanticSearch != "" {
			cfg := config.GetConfig()
			geminiAPIKey := cfg.GeminiAPIKey
			if geminiAPIKey == "" {
				cmd.PrintErr("Error: Gemini API Key is not set. Please set it using 'simple-jot config set gemini-api-key <YOUR_API_KEY>'\n")
				return fmt.Errorf("gemini API key not configured")
			}
			cmd.Println("Performing semantic search with Gemini API...")
			cmd.Printf("Query: %s\n", semanticSearch)
			// TODO: Implement actual Gemini API call here
			searchResults, err := ai.SemanticSearch(noteList, semanticSearch, geminiAPIKey)
			if err != nil {
				return fmt.Errorf("failed to perform semantic search: %w", err)
			}
			searchDF := make([][]string, len(searchResults))
			headers := filteredNotes[0].GetHeaders()
			for i, result := range searchResults {
				searchDF[i] = result.PrepRow()
			}
			err = tabler.RenderTable(searchDF, headers)
			if err != nil {
				return fmt.Errorf("failed to render table: %w", err)
			}
			return nil
		}

		if contentSearch != "" {
			filteredNotes = notes.FilterNotesByContent(filteredNotes, contentSearch)
		}

		if tagStr != "" {
			tagList := strings.Split(tagStr, ",")
			tagMap := tags.TagMap{}
			tagMap.BuildTagMap(noteList)
			noteSet := make(map[string]notes.Note)
			store := notes.NoteStore{}
			store.BuildNoteMap(noteList)

			for _, tagName := range tagList {
				trimmedTag := strings.TrimSpace(tagName)
				relatedNoteIDs := tagMap.GetNotesForTag(trimmedTag)
				if len(relatedNoteIDs) == 0 {
					cmd.PrintErrf("Error identified. No notes found for tag (%s)\n", trimmedTag)
					continue
				}
				for _, noteID := range relatedNoteIDs {
					n, err := store.GetNoteByID(noteID)
					if err != nil {
						cmd.PrintErrf("cannot find note for id (%s). See error: %s\n", noteID, err.Error())
						continue
					}
					noteSet[noteID] = n
				}
			}
			filteredNotes = make([]notes.Note, 0, len(noteSet))
			for _, note := range noteSet {
				filteredNotes = append(filteredNotes, note)
			}
		}

		if dsStr != "" || deStr != "" {
			filteredNotes = notes.FilterNotesByDate(filteredNotes, dsStr, deStr)
		}

		// Prepare table data
		if len(filteredNotes) == 0 {
			cmd.Println("No notes found matching the search criteria.")
			return nil
		}

		dataFrame := make([][]string, len(filteredNotes))
		headers := filteredNotes[0].GetHeaders()

		for i, n := range filteredNotes {
			dataFrame[i] = n.PrepRow()
		}

		err = tabler.RenderTable(dataFrame, headers)
		if err != nil {
			return fmt.Errorf("failed to render table: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Define flags
	searchCmd.Flags().StringP("semantic", "s", "", "Perform a semantic search using Gemini")
	searchCmd.Flags().StringP("content", "c", "", "Search notes by content")
	searchCmd.Flags().StringP("tag", "t", "", "Search notes by tag (comma-separated for multiple tags)")
	searchCmd.Flags().StringP("date-start", "f", "", "Search notes by date start (format: YYYY-MM-DD)")
	searchCmd.Flags().StringP("date-end", "u", "", "Search notes by date end (format: YYYY-MM-DD)")
}
