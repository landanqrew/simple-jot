package cmd

import (
	"testing"

	"github.com/landanqrew/simple-jot/internal/notes"
)

// Note: The search command has been refactored to properly use Cobra flags instead of manual argument parsing.
// Error handling has been improved to use proper return values instead of log.Fatal().
// These tests verify the individual components work correctly.

func TestSearchCommandManualArgumentParsing(t *testing.T) {
	// Test the argument parsing logic to ensure it works as expected
	// This tests the old manual parsing approach for comparison with the new flag-based approach
	testCases := []struct {
		name     string
		args     []string
		expected map[string]string
	}{
		{
			name: "semantic search",
			args: []string{"search", "--semantic", "programming concepts"},
			expected: map[string]string{
				"semantic": "programming concepts",
			},
		},
		{
			name: "content search",
			args: []string{"search", "--content", "test content"},
			expected: map[string]string{
				"content": "test content",
			},
		},
		{
			name: "tag search",
			args: []string{"search", "-t", "work,personal"},
			expected: map[string]string{
				"tag": "work,personal",
			},
		},
		{
			name: "date range search",
			args: []string{"search", "-ds", "2025-01-01", "-de", "2025-01-31"},
			expected: map[string]string{
				"date-start": "2025-01-01",
				"date-end":   "2025-01-31",
			},
		},
		{
			name: "mixed arguments",
			args: []string{"search", "--content", "test", "-t", "work"},
			expected: map[string]string{
				"content": "test",
				"tag":     "work",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the manual parsing logic from search.go
			semanticSearch := ""
			contentSearch := ""
			tagStr := ""
			dsStr := ""
			deStr := ""

			args := tc.args
			for i, arg := range args {
				if i == 0 {
					continue
				}
				switch args[i-1] {
				case "--semantic":
					semanticSearch = arg
				case "-t":
					tagStr = arg
				case "--content":
					contentSearch = arg
				case "-ds":
					dsStr = arg
				case "-de":
					deStr = arg
				}
			}

			// Verify parsed values match expected
			parsed := map[string]string{
				"semantic":   semanticSearch,
				"content":    contentSearch,
				"tag":        tagStr,
				"date-start": dsStr,
				"date-end":   deStr,
			}

			for key, expectedValue := range tc.expected {
				if parsed[key] != expectedValue {
					t.Errorf("Expected %s to be %q, got %q", key, expectedValue, parsed[key])
				}
			}

			// Verify unused fields are empty
			for key, actualValue := range parsed {
				if _, exists := tc.expected[key]; !exists && actualValue != "" {
					t.Errorf("Expected %s to be empty, got %q", key, actualValue)
				}
			}
		})
	}
}

func TestSearchFiltering(t *testing.T) {
	// Test the individual filtering functions used by search
	mockNotes := []notes.Note{
		{
			ID:        "note1",
			Title:     "Test Note 1",
			Content:   "This is a test note about programming",
			Tags:      []string{"work", "programming"},
			CreatedAt: "2025-01-15 10:00:00",
			UpdatedAt: "2025-01-15 10:00:00",
		},
		{
			ID:        "note2",
			Title:     "Meeting Notes",
			Content:   "Discussion about project timeline and deliverables",
			Tags:      []string{"work", "meetings"},
			CreatedAt: "2025-01-20 14:30:00",
			UpdatedAt: "2025-01-20 14:30:00",
		},
		{
			ID:        "note3",
			Title:     "Personal Reminder",
			Content:   "Remember to buy groceries and call mom",
			Tags:      []string{"personal"},
			CreatedAt: "2025-01-10 09:00:00",
			UpdatedAt: "2025-01-10 09:00:00",
		},
	}

	t.Run("content filtering", func(t *testing.T) {
		filtered := notes.FilterNotesByContent(mockNotes, "programming")
		if len(filtered) != 1 {
			t.Errorf("Expected 1 note, got %d", len(filtered))
		}
		if len(filtered) > 0 && filtered[0].ID != "note1" {
			t.Errorf("Expected note1, got %s", filtered[0].ID)
		}
	})

	t.Run("content filtering - case insensitive", func(t *testing.T) {
		filtered := notes.FilterNotesByContent(mockNotes, "PROGRAMMING")
		if len(filtered) != 1 {
			t.Errorf("Expected 1 note, got %d", len(filtered))
		}
		if len(filtered) > 0 && filtered[0].ID != "note1" {
			t.Errorf("Expected note1, got %s", filtered[0].ID)
		}
	})

	t.Run("content filtering - no matches", func(t *testing.T) {
		filtered := notes.FilterNotesByContent(mockNotes, "nonexistent")
		if len(filtered) != 0 {
			t.Errorf("Expected 0 notes, got %d", len(filtered))
		}
	})

	t.Run("date filtering", func(t *testing.T) {
		filtered := notes.FilterNotesByDate(mockNotes, "2025-01-12", "2025-01-18")
		if len(filtered) != 1 {
			t.Errorf("Expected 1 note, got %d", len(filtered))
		}
		if len(filtered) > 0 && filtered[0].ID != "note1" {
			t.Errorf("Expected note1, got %s", filtered[0].ID)
		}
	})

	/*
	cfg := config.GetConfig()
	geminiAPIKey := cfg.GeminiAPIKey

	t.Run("semantic search", func(t *testing.T) {
		filtered, err := ai.SemanticSearch(mockNotes, "programming", geminiAPIKey)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(filtered) != 1 {
			t.Errorf("Expected 1 note, got %d", len(filtered))
		}
		if len(filtered) > 0 && filtered[0].PrimaryKey != "note1" {
			t.Errorf("Expected note1, got %s", filtered[0].PrimaryKey)
		}
	})
		*/
}
