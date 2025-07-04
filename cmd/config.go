/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
	Long: `Manage configuration settings for the application.

Usage:
  simple-jot config set note <note-id>
  simple-jot config get note
`,
	// configCmd itself will not have a direct action, it acts as a container for subcommands.
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// setCmd represents the set subcommand of config
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set application configuration values",
	Long:  `Allows you to set various configuration values for simple-jot.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// getCmd represents the get subcommand of config
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get application configuration values",
	Long: `Allows you to retrieve various configuration values from simple-jot
	
	Usage:
		simple-jot config get note
		simple-jot config get gemini-api-key
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// noteSetCmd represents the note subcommand of config set
var noteSetCmd = &cobra.Command{
	Use:   "note <note-id>",
	Short: "Set the current active note ID",
	Long:  `Sets the specified note ID as the active note in the configuration.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		noteID := args[0]
		viper.Set("active_note", noteID)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error saving configuration:", err)
			os.Exit(1)
		}
		fmt.Printf("Active note set to: %s\n", noteID)
	},
}

// geminiAPIKeySetCmd represents the gemini-api-key subcommand of config set
var geminiAPIKeySetCmd = &cobra.Command{
	Use:   "gemini-api-key <api-key>",
	Short: "Set the Gemini API key",
	Long:  `Sets the Gemini API key in the configuration for semantic search.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := args[0]
		viper.Set("gemini_api_key", apiKey)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error saving configuration:", err)
			os.Exit(1)
		}
		fmt.Printf("Gemini API Key set successfully.\n")
	},
}

// noteGetCmd represents the note subcommand of config get
var noteGetCmd = &cobra.Command{
	Use:   "note",
	Short: "Get the current active note ID",
	Long:  `Retrieves the currently active note ID from the configuration.`,
	Args:  cobra.NoArgs, // Expects no arguments
	Run: func(cmd *cobra.Command, args []string) {
		activeNote := viper.GetString("active_note")
		if activeNote == "" {
			fmt.Println("No active note is currently set.")
		} else {
			fmt.Printf("Current active note: %s\n", activeNote)
		}
	},
}

// geminiAPIKeyGetCmd represents the gemini-api-key subcommand of config get
var geminiAPIKeyGetCmd = &cobra.Command{
	Use:   "gemini-api-key",
	Short: "Get the Gemini API key",
	Long:  `Retrieves the configured Gemini API key.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("gemini_api_key")
		if apiKey == "" {
			fmt.Println("No Gemini API Key is currently set.")
		} else {
			fmt.Printf("Gemini API Key: %s\n", apiKey)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(getCmd)

	setCmd.AddCommand(noteSetCmd)
	setCmd.AddCommand(geminiAPIKeySetCmd)
	getCmd.AddCommand(noteGetCmd)
	getCmd.AddCommand(geminiAPIKeyGetCmd)

	// No flags directly on configCmd anymore, they are on subcommands if needed.
}
