package cmd

import (
	"github.com/landanqrew/simple-jot/internal/notes"
)

// mockStorage implements storage.NoteStorage interface for testing
type mockStorage struct {
	notes     []notes.Note
	getError  error
	saveError error
}

func (m *mockStorage) GetNotes() ([]notes.Note, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return m.notes, nil
}

func (m *mockStorage) SaveNotes(notes []notes.Note) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.notes = notes
	return nil
}
