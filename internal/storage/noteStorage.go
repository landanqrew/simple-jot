package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/landanqrew/simple-jot/internal/notes"
)

// NoteStorage defines the interface for note storage operations
type NoteStorage interface {
	GetNotes() ([]notes.Note, error)
	SaveNotes([]notes.Note) error
}

// FileNoteStorage implements NoteStorage using the filesystem
type FileNoteStorage struct {
	filePath string
}

// NewFileNoteStorage creates a new FileNoteStorage instance
func NewFileNoteStorage(filePath string) *FileNoteStorage {
	return &FileNoteStorage{filePath: filePath}
}

// GetNotes retrieves all notes from storage
func (s *FileNoteStorage) GetNotes() ([]notes.Note, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []notes.Note{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read notes file: %v", err)
	}

	var notes []notes.Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return nil, fmt.Errorf("failed to parse notes: %v", err)
	}

	return notes, nil
}

// SaveNotes saves all notes to storage
func (s *FileNoteStorage) SaveNotes(notes []notes.Note) error {
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes: %v", err)
	}

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notes to file: %v", err)
	}

	return nil
}

// Default storage instance
var defaultStorage NoteStorage = NewFileNoteStorage("notes.json")

// GetNotes is a convenience function that uses the default storage
func GetNotes() ([]notes.Note, error) {
	return defaultStorage.GetNotes()
}

// SaveNotes is a convenience function that uses the default storage
func SaveNotes(notes []notes.Note) error {
	return defaultStorage.SaveNotes(notes)
}

// SetDefaultStorage allows changing the default storage implementation
func SetDefaultStorage(storage NoteStorage) {
	defaultStorage = storage
}
