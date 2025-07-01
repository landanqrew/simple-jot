/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/internal/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all notes",
	Long: `List all notes:

Usage:
simple-jot list`,
	Run: func(cmd *cobra.Command, args []string) {
		noteList, err := storage.GetNotes()
		if err != nil {
			log.Fatal("cannot fetch notes. See error:", err.Error())
		}
		exNote := notes.Note{}
		headers := exNote.GetHeaders()
		dataFrame := make([][]string, len(noteList))
		for i, n := range noteList {
			dataFrame[i] = n.PrepRow()
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.Header(headers)
		table.Bulk(dataFrame)
		err = table.Render()
		if err != nil {
			log.Fatalln("Failed to render table for your notes query. Exiting Program")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
