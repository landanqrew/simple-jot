package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application's configuration settings.
type Config struct {
	NotesDirectory string `mapstructure:"notes_directory"`
	ActiveNote     string `mapstructure:"active_note"`
	DataDir        string `mapstructure:"data_dir"` // Directory where notes data will be stored
	Editor         string `mapstructure:"editor"`   // Preferred text editor for editing notes (e.g., "vim", "nano", "code")
	NoteID         string `mapstructure:"note_id"`  // The ID of the active note
	// Add other configuration fields as your application grows
}

// globalConfig stores the loaded configuration.
var globalConfig *Config

// InitConfig initializes the Viper configuration system and loads settings.
// It searches for a config file, sets defaults, and binds environment variables.
func InitConfig() error {
	// Set the name of the config file (without extension)
	viper.SetConfigName(".note-cli")
	// Set the type of the config file
	viper.SetConfigType("yaml")

	// Add search paths for the config file
	// 1. Current working directory
	viper.AddConfigPath(".")
	// 2. User's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	viper.AddConfigPath(filepath.Join(home, ".config", "note-cli")) // Standard XDG config dir
	viper.AddConfigPath(home)                                       // Fallback to home directory directly

	// Set default values for configuration options
	defaultDataDir := filepath.Join(home, ".note-cli", "data")
	viper.SetDefault("data_dir", defaultDataDir)
	viper.SetDefault("editor", os.Getenv("EDITOR")) // Use EDITOR env var as default for editor

	// Read environment variables (e.g., NOTECLI_DATA_DIR, NOTECLI_EDITOR)
	viper.SetEnvPrefix("NOTECLI") // Prefix for environment variables (e.g., NOTECLI_DATA_DIR)
	viper.AutomaticEnv()          // Read matching environment variables

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist, we'll use defaults or env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, using defaults and environment variables.")
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal the configuration into the Config struct
	globalConfig = &Config{}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Ensure the data directory exists
	if _, err := os.Stat(globalConfig.DataDir); os.IsNotExist(err) {
		fmt.Printf("Data directory does not exist, creating: %s\n", globalConfig.DataDir)
		if err := os.MkdirAll(globalConfig.DataDir, 0755); err != nil {
			return fmt.Errorf("failed to create data directory %s: %w", globalConfig.DataDir, err)
		}
	}

	return nil
}

// GetConfig returns the loaded application configuration.
// It panics if InitConfig has not been called successfully.
func GetConfig() *Config {
	if globalConfig == nil {
		panic("Config not initialized. Call InitConfig() first.")
	}
	return globalConfig
}

// Example usage within your root.go or main.go after initialization:
/*
import "your_module_name/internal/config"

func main() {
    if err := config.InitConfig(); err != nil {
        fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
        os.Exit(1)
    }
    appConfig := config.GetConfig()
    fmt.Printf("Notes will be stored in: %s\n", appConfig.DataDir)
    fmt.Printf("Default editor: %s\n", appConfig.Editor)
    cmd.Execute()
}
*/
