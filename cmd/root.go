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

var cfgFile string = "./internal/config/config.go"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simple-jot",
	Short: "A simple cli tool for taking and managing notes",
	Long: `simple-jot is a simple cli tool for taking and managing notes within the context of a project.

## Usage Examples: ##

to get started, run:

	simple-jot init

to set up your configuration, run:

	simple-jot config set note <note-id> 

to get the value of a configuration, run:

	simple-jot config get note

to configure your Gemini API Key:

	simple-jot config set gemini-api-key <YOUR_API_KEY>

to view your configured Gemini API Key:

	simple-jot config get gemini-api-key

to create a new note, run:

	simple-jot create <note-name> -n "<note-content>" -s (optional - will set the note configuration to the new note)

to list all notes, run:

	simple-jot list

to search for a note, run:

	simple-jot search <query>

to run a semantic search with llm agent, run:

	simple-jot search --semantic <query>

to edit a note, run:

	simple-jot edit <note-id> -n "<note-content>"

to tag a note, run:

	simple-jot tag <note-id> (optional - will default to the current note) <tag>

to delete a note, run:

	simple-jot delete <note-id>`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.simple-jot.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".simple-jot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".simple-jot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
