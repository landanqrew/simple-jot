/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/landanqrew/simple-jot/internal/config"
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
  simple-jot config set gemini-api-key <api-key>
  simple-jot config get gemini-api-key
`,
	// configCmd itself will not have a direct action, it acts as a container for subcommands.
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// setCmd represents the set subcommand of config
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set application configuration values",
	Long:  `Allows you to set various configuration values for simple-jot.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// noteSetCmd represents the note subcommand of config set
var noteSetCmd = &cobra.Command{
	Use:   "note <note-id>",
	Short: "Set the current active note ID",
	Long:  `Sets the specified note ID as the active note in the configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteID := args[0]

		viper.Set("active_note", noteID)

		// Try to write the config, and if it fails because no config file exists, create one
		err := viper.WriteConfig()
		if err != nil {
			// If writing fails, try to write a new config file
			home, homeErr := os.UserHomeDir()
			if homeErr != nil {
				return fmt.Errorf("failed to get home directory: %w", homeErr)
			}

			configPath := fmt.Sprintf("%s/.simple-jot.yaml", home)
			err = viper.WriteConfigAs(configPath)
			if err != nil {
				return fmt.Errorf("error creating configuration file: %w", err)
			}
		}

		cmd.Printf("Active note set to: %s\n", noteID)
		return nil
	},
}

// geminiAPIKeySetCmd represents the gemini-api-key subcommand of config set
var geminiAPIKeySetCmd = &cobra.Command{
	Use:   "gemini-api-key <api-key>",
	Short: "Set the Gemini API key",
	Long:  `Sets the Gemini API key in the configuration for semantic search.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := args[0]

		viper.Set("gemini_api_key", apiKey)

		// Try to write the config, and if it fails because no config file exists, create one
		err := viper.WriteConfig()
		if err != nil {
			// If writing fails, try to write a new config file
			home, homeErr := os.UserHomeDir()
			if homeErr != nil {
				return fmt.Errorf("failed to get home directory: %w", homeErr)
			}

			configPath := fmt.Sprintf("%s/.simple-jot.yaml", home)
			err = viper.WriteConfigAs(configPath)
			if err != nil {
				return fmt.Errorf("error creating configuration file: %w", err)
			}
		}

		cmd.Printf("Gemini API Key set successfully.\n")
		return nil
	},
}

// noteGetCmd represents the note subcommand of config get
var noteGetCmd = &cobra.Command{
	Use:   "note",
	Short: "Get the current active note ID",
	Long:  `Retrieves the currently active note ID from the configuration.`,
	Args:  cobra.NoArgs, // Expects no arguments
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetConfig()
		if cfg.ActiveNote == "" {
			cmd.Println("No active note is currently set.")
		} else {
			cmd.Println(cfg.ActiveNote)
		}
		return nil
	},
}

// geminiAPIKeyGetCmd represents the gemini-api-key subcommand of config get
var geminiAPIKeyGetCmd = &cobra.Command{
	Use:   "gemini-api-key",
	Short: "Get the Gemini API key",
	Long:  `Retrieves the configured Gemini API key.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetConfig()
		if cfg.GeminiAPIKey == "" {
			cmd.Println("No Gemini API Key is currently set.")
		} else {
			cmd.Printf("Gemini API Key: %s\n", cfg.GeminiAPIKey)
		}
		return nil
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
