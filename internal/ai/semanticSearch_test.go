package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/landanqrew/simple-jot/internal/notes"
	"github.com/landanqrew/simple-jot/tabler"
)

func TestSemanticSearch(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Fatalf("GEMINI_API_KEY not set, please run \n`export GEMINI_API_KEY=<your-api-key>`\nand run the test again")
	}

	// Read notes from testNotes.json
	filePath := "../notes/testNotes.json"
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read testNotes.json: %v", err)
	}

	var data []notes.Note
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal notes from testNotes.json: %v", err)
	}
	fmt.Println("notes identified: ", len(data))

	query := "Rank notes by relevance to 'Go Programming'."

	var results []SearchResponse
	results, err = SemanticSearch(data, query, apiKey)
	if err != nil {
		t.Fatalf("SemanticSearch returned an error: %v", err)
	}
	if results == nil {
		t.Fatal("SemanticSearch returned nil results")
	}

	if len(results) == 0 {
		t.Error("Expected at least one search result, got 0")
	}

	preppedRows := make([]tabler.RowPrepper, len(results))
	for i, res := range results {
		preppedRows[i] = res
	}
	dataFrame := tabler.PrepTable(preppedRows, []string{"PrimaryKey", "Score"})
	tabler.RenderTable(dataFrame, []string{"PrimaryKey", "Score"})

	// Assuming the 6th note (index 5) in the generated data is the 'Go Programming Best Practices' note
	// This assertion relies on the specific content and ordering generated in testNotes.json
	expectedGoNoteID := data[5].ID
	foundGoNote := false
	for _, res := range results {
		if res.PrimaryKey == expectedGoNoteID && res.Score > 0.7 {
			foundGoNote = true
			break
		}
	}

	if !foundGoNote {
		t.Errorf("Expected 'Go Programming Best Practices' note (ID %s) to be highly ranked, but it wasn't or score was too low", expectedGoNoteID)
	}
}
