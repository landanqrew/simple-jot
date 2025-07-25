package notes

import (
	"fmt"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/landanqrew/simple-jot/internal/osutils"
)

// functions to test:
/**
AddTag
RemoveTag
GetTags
CheckContentMatch
UpdateContent
GetHeaders
PrepRow
FilterNotesByContent
*/

func TestAddTag(t *testing.T) {
	note := Note{
		ID:        "1",
		Title:     "Test Note",
		Tags:      []string{"tag1", "tag2"},
		Content:   "This is a test note",
		CreatedAt: "2021-01-01 01:02:03",
		UpdatedAt: "2021-01-01 01:02:03",
	}

	note.AddTag("tag3")

	json, err := osutils.ToJsonString([]Note{note})
	if err != nil {
		t.Fatalf("Failed to convert note to JSON: %v", err)
	}

	fmt.Println("note:\n", json)

	if len(note.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(note.Tags))
	}

	if !slices.Contains(note.Tags, "tag3") {
		t.Errorf("Expected tag3 to be in tags")
	}

	if note.UpdatedAt != time.Now().Format(time.DateTime) {
		t.Errorf("Expected updated_at to be the current time, got %s", note.UpdatedAt)
	}
}

func TestRemoveTag(t *testing.T) {
	note := Note{
		ID:        "1",
		Title:     "Test Note",
		Tags:      []string{"tag1", "tag2", "tag3"},
		Content:   "This is a test note",
		CreatedAt: "2021-01-01 01:02:03",
		UpdatedAt: "2021-01-01 01:02:03",
	}

	note.RemoveTag("tag2")

	if len(note.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(note.Tags))
	}

	if slices.Contains(note.Tags, "tag2") {
		t.Errorf("Expected tag2 to be removed")
	}

	if note.UpdatedAt != time.Now().Format(time.DateTime) {
		t.Errorf("Expected updated_at to be the current time, got %s", note.UpdatedAt)
	}
}

func TestGetTags(t *testing.T) {
	note := Note{
		ID:   "1",
		Tags: []string{"tag1", "tag2"},
	}

	tags := note.GetTags()

	if len(tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(tags))
	}

	if !slices.Contains(tags, "tag1") || !slices.Contains(tags, "tag2") {
		t.Errorf("Expected tags to contain tag1 and tag2")
	}
}

func TestCheckContentMatch(t *testing.T) {
	note := Note{
		ID:      "1",
		Content: "This is a test note about golang.",
	}

	if !note.CheckContentMatch("golang") {
		t.Errorf("Expected content to match 'golang'")
	}

	if note.CheckContentMatch("python") {
		t.Errorf("Expected content not to match 'python'")
	}

	if !note.CheckContentMatch("GoLaNg") {
		t.Errorf("Expected case-insensitive match for 'GoLaNg'")
	}
}

func TestUpdateContent(t *testing.T) {
	note := Note{
		ID:        "1",
		Content:   "Old content",
		UpdatedAt: "2021-01-01 01:02:03",
	}

	newContent := "New content"
	note.UpdateContent(newContent)

	if note.Content != newContent {
		t.Errorf("Expected content to be updated to '%s', got '%s'", newContent, note.Content)
	}

	if note.UpdatedAt != time.Now().Format(time.DateTime) {
		t.Errorf("Expected updated_at to be the current time, got %s", note.UpdatedAt)
	}

	// Test with empty content (should not update)
	oldUpdatedAt := note.UpdatedAt
	note.UpdateContent("")
	if note.Content != newContent || note.UpdatedAt != oldUpdatedAt {
		t.Errorf("Expected content and updated_at not to change with empty content")
	}
}

func TestGetHeaders(t *testing.T) {
	note := Note{}
	headers := note.GetHeaders()

	expectedHeaders := []string{"ID", "Title", "Tags", "Content", "CreatedAt", "UpdatedAt"}
	if !slices.Equal(headers, expectedHeaders) {
		t.Errorf("Expected headers %v, got %v", expectedHeaders, headers)
	}
}

func TestPrepRow(t *testing.T) {
	// Test case with all fields populated
	note1 := Note{
		ID:        "id1",
		Title:     "Title1",
		Tags:      []string{"tagA", "tagB"},
		Content:   "Content1",
		CreatedAt: "2023-01-01 10:00:00",
		UpdatedAt: "2023-01-01 11:00:00",
	}
	expectedRow1 := []string{"id1", "Title1", "tagA, tagB", "Content1", "2023-01-01 10:00:00", "2023-01-01 11:00:00"}
	actualRow1 := note1.PrepRow()

	if !slices.Equal(actualRow1, expectedRow1) {
		t.Errorf("PrepRow (Note 1) failed. Expected %v, Got %v", expectedRow1, actualRow1)
	}

	// Test case with empty tags
	note2 := Note{
		ID:        "id2",
		Title:     "Title2",
		Tags:      []string{},
		Content:   "Content2",
		CreatedAt: "2023-02-01 12:00:00",
		UpdatedAt: "2023-02-01 13:00:00",
	}
	expectedRow2 := []string{"id2", "Title2", "", "Content2", "2023-02-01 12:00:00", "2023-02-01 13:00:00"}
	actualRow2 := note2.PrepRow()

	if !slices.Equal(actualRow2, expectedRow2) {
		t.Errorf("PrepRow (Note 2) failed. Expected %v, Got %v", expectedRow2, actualRow2)
	}

	// Test case with a single tag
	note3 := Note{
		ID:        "id3",
		Title:     "Title3",
		Tags:      []string{"single_tag"},
		Content:   "Content3",
		CreatedAt: "2023-03-01 14:00:00",
		UpdatedAt: "2023-03-01 15:00:00",
	}
	expectedRow3 := []string{"id3", "Title3", "single_tag", "Content3", "2023-03-01 14:00:00", "2023-03-01 15:00:00"}
	actualRow3 := note3.PrepRow()

	if !slices.Equal(actualRow3, expectedRow3) {
		t.Errorf("PrepRow (Note 3) failed. Expected %v, Got %v", expectedRow3, actualRow3)
	}
}

func TestFilterNotesByContent(t *testing.T) {
	notes := []Note{
		{ID: "1", Content: "This is about Go programming."},
		{ID: "2", Content: "Python for data science."},
		{ID: "3", Content: "Learn Go and Rust."}, // Should match "go"
		{ID: "4", Content: "Java development"},
	}

	// Test case: case-insensitive match
	filteredGo := FilterNotesByContent(notes, "go")
	if len(filteredGo) != 2 {
		t.Errorf("Expected 2 notes for 'go', got %d", len(filteredGo))
	}
	if filteredGo[0].ID != "1" || filteredGo[1].ID != "3" {
		t.Errorf("Expected notes with ID 1 and 3 for 'go', got %v", filteredGo)
	}

	// Test case: no match
	filteredRust := FilterNotesByContent(notes, "rust")
	if len(filteredRust) != 1 || filteredRust[0].ID != "3" {
		t.Errorf("Expected 1 note for 'rust', got %d", len(filteredRust))
	}

	// Test case: empty query
	filteredEmpty := FilterNotesByContent(notes, "")
	if len(filteredEmpty) != len(notes) {
		t.Errorf("Expected all notes for empty query, got %d", len(filteredEmpty))
	}
}

func TestFilterNotesByDate(t *testing.T) {
	path := "./internal/notes/testNotes.json"
	rawBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read testNotes.json: %v", err)
	}

	testNotes, err := osutils.ReadJson[Note](rawBytes)
	if err != nil {
		t.Fatalf("Failed to read testNotes.json: %v", err)
	}

	filteredNotes := FilterNotesByDate(testNotes, "2021-01-01", "2021-03-01")
	expectedNoteCount := 2
	if len(filteredNotes) != expectedNoteCount {
		t.Errorf("Expected %d notes for '2021-01-01', got %d", expectedNoteCount, len(filteredNotes))
	}

	expectedNoteIds := []string{"6d3c8b6d-6988-4d1c-9d67-f51e9cea20da", "f165afab-dd12-465b-9e8d-c39d6fb605f2"}
	for _, note := range filteredNotes {
		if !slices.Contains(expectedNoteIds, note.ID) {
			t.Errorf("Expected note ID %s to be in filtered notes", note.ID)
		}
	}
}


