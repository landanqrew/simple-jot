/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	storage "github.com/landanqrew/simple-jot/internal/storage"
	tags "github.com/landanqrew/simple-jot/internal/tags"
	tablewriter "github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "view tags",
	Long: `utility to view tags:

Usage:
simple-jot tag list <optional-prefix>`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var tagListCmd = &cobra.Command {
	Use: "list",
	Short: "list tags (optionally filter by prefix)",
	Long: `list tags (optionally filter by prefix)

Usage:
simple-jot tag list <optional-prefix>`,
	Run: func(cmd *cobra.Command, args []string) {
		prefix := ""
		if len(args) > 0 {
			prefix = args[0]
		}
		notes, err := storage.GetNotes()
		if err != nil {
			log.Fatalln("Failed to get notes. Exiting Program")
		}
		tagMap := tags.TagMap{}
		tagMap.BuildTagMap(notes)
		dataFrame := make([][]string, 0)
		dataFrame = append(dataFrame, []string{"tag", "notes", "noteCount"})

		for tag, val := range tagMap.TagMap {
			if prefix != "" && !strings.HasPrefix(tag, prefix) {
				continue
			}
			row := make([]string, 0)
			row = append(row, tag)
			row = append(row, "")
			row = append(row, strconv.Itoa(len(val)))
			for noteID, _ := range val {
				row[1] = row[1] + noteID + ", "
			}
			row[1] = row[1][:len(row[1])-2]
			dataFrame = append(dataFrame, row)
		}

		if len(dataFrame) == 1 {
			fmt.Println("No tags found matching your query (prefix: ", prefix, ")")
			return
		}
		
		table := tablewriter.NewWriter(os.Stdout)
		table.Header(dataFrame[0])
		table.Bulk(dataFrame[1:])
		err = table.Render()
		if err != nil {
			log.Fatalln("Failed to render table for your tag query. Exiting Program")
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.AddCommand(tagListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
