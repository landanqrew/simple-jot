package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	// Save original environment and config state
	originalGlobalConfig := globalConfig

	// Restore original state after all tests
	defer func() {
		viper.Reset()
		globalConfig = originalGlobalConfig
	}()

	t.Run("InitConfig with defaults", func(t *testing.T) {
		viper.Reset()
		globalConfig = nil

		err := InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() failed: %v", err)
		}

		config := GetConfig()
		if config == nil {
			t.Fatal("Expected config to be non-nil")
		}

		// Check that data directory is set to default
		home, _ := os.UserHomeDir()
		expectedDataDir := filepath.Join(home, ".note-cli", "data")
		if config.DataDir != expectedDataDir {
			t.Errorf("Expected DataDir %q, got %q", expectedDataDir, config.DataDir)
		}

		// Check that data directory was created
		if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
			t.Errorf("Expected data directory %q to be created", config.DataDir)
		}

		// Clean up created directory
		os.RemoveAll(filepath.Join(home, ".note-cli"))
	})

	t.Run("InitConfig with config file", func(t *testing.T) {
		viper.Reset()
		globalConfig = nil

		// Create a temporary directory and config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, ".note-cli.yaml")
		testDataDir := filepath.Join(tempDir, "test-notes")

		configContent := `
data_dir: ` + testDataDir + `
editor: code
active_note: test-note-id
gemini_api_key: file-api-key
notes_directory: ` + filepath.Join(tempDir, "notes") + `
note_id: active-note-123
`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// Change to temp directory so config file is found
		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		err = InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() failed: %v", err)
		}

		config := GetConfig()

		if config.DataDir != testDataDir {
			t.Errorf("Expected DataDir %q, got %q", testDataDir, config.DataDir)
		}

		if config.Editor != "code" {
			t.Errorf("Expected Editor %q, got %q", "code", config.Editor)
		}

		if config.ActiveNote != "test-note-id" {
			t.Errorf("Expected ActiveNote %q, got %q", "test-note-id", config.ActiveNote)
		}

		if config.GeminiAPIKey != "file-api-key" {
			t.Errorf("Expected GeminiAPIKey %q, got %q", "file-api-key", config.GeminiAPIKey)
		}

		expectedNotesDir := filepath.Join(tempDir, "notes")
		if config.NotesDirectory != expectedNotesDir {
			t.Errorf("Expected NotesDirectory %q, got %q", expectedNotesDir, config.NotesDirectory)
		}

		if config.NoteID != "active-note-123" {
			t.Errorf("Expected NoteID %q, got %q", "active-note-123", config.NoteID)
		}
	})

	t.Run("InitConfig with invalid config file", func(t *testing.T) {
		viper.Reset()
		globalConfig = nil

		// Create a temporary directory and invalid config file
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, ".note-cli.yaml")

		// Invalid YAML content
		invalidConfigContent := `
invalid yaml content:
  - unclosed bracket [
  malformed: content
`

		err := os.WriteFile(configFile, []byte(invalidConfigContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// Change to temp directory so config file is found
		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		err = InitConfig()
		if err == nil {
			t.Fatal("Expected InitConfig() to fail with invalid config file")
		}

		// Should contain "failed to read config file" in error message
		if err.Error()[:24] != "failed to read config fi" {
			t.Errorf("Expected error about config file, got: %v", err)
		}
	})

	t.Run("InitConfig with permission denied for data directory", func(t *testing.T) {
		viper.Reset()
		globalConfig = nil

		// Try to create data directory in a location without permissions
		if os.Getuid() == 0 {
			t.Skip("Skipping permission test when running as root")
		}

		// Create a temp config file with an invalid data directory
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, ".note-cli.yaml")

		configContent := `data_dir: /root/test-data`

		err := os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		// Change to temp directory so config file is found
		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		err = InitConfig()
		if err == nil {
			t.Fatal("Expected InitConfig() to fail when data directory can't be created")
		}

		// Should contain "failed to create data directory" in error message
		if err.Error()[:30] != "failed to create data director" {
			t.Errorf("Expected error about data directory creation, got: %v", err)
		}
	})

	t.Run("GetConfig without initialization", func(t *testing.T) {
		viper.Reset()
		globalConfig = nil

		// This should panic since InitConfig hasn't been called
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected GetConfig() to panic when not initialized")
			} else {
				expectedPanic := "Config not initialized. Call InitConfig() first."
				if r != expectedPanic {
					t.Errorf("Expected panic message %q, got %q", expectedPanic, r)
				}
			}
		}()

		GetConfig()
	})

	t.Run("Multiple InitConfig calls", func(t *testing.T) {
		viper.Reset()
		globalConfig = nil

		// Use temp directory to avoid conflicts
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, ".note-cli.yaml")

		// First initialization with one config
		configContent1 := `
data_dir: ` + filepath.Join(tempDir, "data1") + `
editor: vim
`
		err := os.WriteFile(configFile, []byte(configContent1), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		err = InitConfig()
		if err != nil {
			t.Fatalf("First InitConfig() failed: %v", err)
		}

		config1 := GetConfig()

		// Second initialization with different config
		configContent2 := `
data_dir: ` + filepath.Join(tempDir, "data2") + `
editor: emacs
`
		err = os.WriteFile(configFile, []byte(configContent2), 0644)
		if err != nil {
			t.Fatalf("Failed to update test config file: %v", err)
		}

		err = InitConfig()
		if err != nil {
			t.Fatalf("Second InitConfig() failed: %v", err)
		}

		config2 := GetConfig()

		// Should be different instances
		if config1 == config2 {
			t.Error("Expected different config instances after re-initialization")
		}

		if config2.Editor != "emacs" {
			t.Errorf("Expected Editor %q after re-init, got %q", "emacs", config2.Editor)
		}

		expectedDataDir2 := filepath.Join(tempDir, "data2")
		if config2.DataDir != expectedDataDir2 {
			t.Errorf("Expected DataDir %q after re-init, got %q", expectedDataDir2, config2.DataDir)
		}
	})
}

func TestConfigEdgeCases(t *testing.T) {
	// Save original state
	originalGlobalConfig := globalConfig

	// Restore original state after tests
	defer func() {
		viper.Reset()
		globalConfig = originalGlobalConfig
	}()

	t.Run("Data directory already exists", func(t *testing.T) {
		viper.Reset()
		globalConfig = nil

		// Create a temporary directory that already exists
		tempDir := t.TempDir()
		testDataDir := filepath.Join(tempDir, "existing-data")

		// Pre-create the directory
		err := os.MkdirAll(testDataDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create config file that points to existing directory
		configFile := filepath.Join(tempDir, ".note-cli.yaml")
		configContent := `data_dir: ` + testDataDir
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		err = InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() failed: %v", err)
		}

		config := GetConfig()
		if config.DataDir != testDataDir {
			t.Errorf("Expected DataDir %q, got %q", testDataDir, config.DataDir)
		}

		// Directory should still exist
		if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
			t.Errorf("Expected data directory %q to exist", config.DataDir)
		}
	})

	t.Run("Config with EDITOR environment variable default", func(t *testing.T) {
		viper.Reset()
		globalConfig = nil

		// Set EDITOR environment variable
		originalEditor := os.Getenv("EDITOR")
		os.Setenv("EDITOR", "nano")
		defer func() {
			if originalEditor != "" {
				os.Setenv("EDITOR", originalEditor)
			} else {
				os.Unsetenv("EDITOR")
			}
		}()

		// Use a temporary directory for the config
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, ".note-cli.yaml")
		configContent := `data_dir: ` + filepath.Join(tempDir, "data")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}

		originalWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(originalWd)

		err = InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() failed: %v", err)
		}

		config := GetConfig()

		// Should use EDITOR env var as default
		if config.Editor != "nano" {
			t.Errorf("Expected Editor %q from EDITOR env var, got %q", "nano", config.Editor)
		}
	})
}

// TestConfigDirectManipulation tests Viper configuration directly
func TestConfigDirectManipulation(t *testing.T) {
	defer viper.Reset()

	t.Run("Viper environment variable mapping", func(t *testing.T) {
		viper.Reset()

		// Test direct Viper configuration for environment variables
		viper.SetEnvPrefix("NOTECLI")
		viper.AutomaticEnv()

		// Test setting a value and getting it back
		viper.Set("test_key", "test_value")
		if viper.GetString("test_key") != "test_value" {
			t.Errorf("Expected test_key to be 'test_value', got %q", viper.GetString("test_key"))
		}

		// Test default values
		viper.SetDefault("default_key", "default_value")
		if viper.GetString("default_key") != "default_value" {
			t.Errorf("Expected default_key to be 'default_value', got %q", viper.GetString("default_key"))
		}
	})

	t.Run("Config struct unmarshaling", func(t *testing.T) {
		viper.Reset()

		// Set some values in Viper
		viper.Set("data_dir", "/test/path")
		viper.Set("editor", "vim")
		viper.Set("active_note", "note-123")

		// Unmarshal into config struct
		var testConfig Config
		err := viper.Unmarshal(&testConfig)
		if err != nil {
			t.Fatalf("Failed to unmarshal config: %v", err)
		}

		if testConfig.DataDir != "/test/path" {
			t.Errorf("Expected DataDir '/test/path', got %q", testConfig.DataDir)
		}

		if testConfig.Editor != "vim" {
			t.Errorf("Expected Editor 'vim', got %q", testConfig.Editor)
		}

		if testConfig.ActiveNote != "note-123" {
			t.Errorf("Expected ActiveNote 'note-123', got %q", testConfig.ActiveNote)
		}
	})
}
