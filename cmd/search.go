/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/landanqrew/simple-jot/internal/tags"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for notes",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Usage:
to search for a note, run:
	simple-jot search -ds 2025-01-01 -de 2025-02-01
	simple-jot search -t "tag1, tag2"
to run a semantic search with llm agent, run:
	simple-jot search --semantic <query>`,
	Run: func(cmd *cobra.Command, args []string) {
		noteList, err := storage.GetNotes()
		if err != nil {
			log.Fatal("cannot fetch notes", err.Error())
		}
		filteredNotes := noteList
		if slices.Contains(args, "-t") {
			tagStr := ""
			for i, arg := range args {
				if i == 0 {
					continue
				}
				if args[i - 1] == "-t" {
					tagStr = arg
				}
			}
			if tagStr != "" {
				tagList := strings.Split(tagStr, ",")
				tagMap := tags.TagMap{}
				tagMap.BuildTagMap(noteList)
				noteSet := make(map[string]notes.Note)
				store := notes.NoteStore{}
				store.BuildNoteMap(noteList)
				for _, tagName := range tagList {
					relatedNoteIDs := tagMap.GetNotesForTag(strings.Trim(tagName, " "))
					if len(relatedNoteIDs) == 0 {
						fmt.Println("Error identified. No notes found for tag (" + tagName + ")")
						continue
					}
					for _, noteID := range relatedNoteIDs {
						n, err := store.GetNoteByID(noteID)
						if err != nil {
							fmt.Println("cannot find note for id (" + noteID + "). See error: " + err.Error())
						}
						noteSet[noteID] = n
					}
				}
				filteredNotes = make([]notes.Note, 0)
				for _, note := range noteSet {
					filteredNotes = append(filteredNotes, note)
				}
			}
		}
		dsStr, deStr := "", ""
		if slices.Contains(args, "-ds") {
			for i, arg := range args {
				if i == 0 {
					continue
				}
				if args[i - 1] == "-ds" {
					dsStr = arg
				}
			}
		}
		if slices.Contains(args, "-de") {
			for i, arg := range args {
				if i == 0 {
					continue
				}
				if args[i - 1] == "-de" {
					deStr = arg
				}
			}
		}
		if dsStr != "" || deStr != "" {
			filteredNotes = notes.FilterNotesByDate(filteredNotes, dsStr, deStr)
		}
		dataFrame := make([][]string, len(filteredNotes) + 1)
		dataFrame[0] = filteredNotes[0].GetHeaders()
		for i, n := range filteredNotes {
			row := n.PrepRow()
			dataFrame[i + 1] = row
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header(dataFrame[0])
		table.Bulk(dataFrame[1:])
		err = table.Render()
		if err != nil {
			log.Fatalln("Failed to render table for your notes query. Exiting Program")
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
