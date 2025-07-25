package notes

import (
	"os"
	"strings"
	"testing"

	"github.com/landanqrew/simple-jot/internal/osutils"
)

func TestNoteStore(t *testing.T) {
	jsonData, err := os.ReadFile("testNotes.json")
	if err != nil {
		t.Fatal("cannot read test notes: " + err.Error())
	}
	notes, err := osutils.ReadJson[Note](jsonData)
	if err != nil {
		t.Fatal("cannot serialize test notes: " + err.Error())
	}
	noteStore := NoteStore{
		NoteMap: make(map[string]Note),
	}
	noteStore.BuildNoteMap(notes)

	if len(noteStore.NoteMap) != len(notes) {
		t.Fatal("note map length does not match notes length")
	}

	for _, note := range notes {
		storeNote := noteStore.NoteMap[note.ID]
		noteJson, _ := osutils.ToJsonString(note)
		storeNoteJson, _ := osutils.ToJsonString(storeNote)
		unitFail := false
		if storeNote.ID != note.ID {
			t.Error("note map does not match notes: " + note.ID)
			unitFail = true
		}
		if storeNote.Title != note.Title {
			t.Error("note map does not match notes: " + note.ID)
			unitFail = true
		}
		if storeNote.Content != note.Content {
			t.Error("note map does not match notes: " + note.ID)
			unitFail = true
		}
		if storeNote.CreatedAt != note.CreatedAt {
			t.Fatal("note map does not match notes: " + note.ID)
			unitFail = true
		}
		if storeNote.UpdatedAt != note.UpdatedAt {
			t.Fatal("note map does not match notes: " + note.ID)
			unitFail = true
		}
		if strings.Join(storeNote.Tags, ",") != strings.Join(note.Tags, ",") {
			t.Error("note map does not match notes: " + note.ID)
			unitFail = true
		}
		if unitFail {
			t.Log("note: " + noteJson)
			t.Log("store note: " + storeNoteJson)
		}
	}

	
}