package storage

import (
	"encoding/json"
	"errors"
	"os"

	notes "github.com/landanqrew/simple-jot/internal/notes"
	"github.com/spf13/viper"
)

func GetNotes() ([]notes.Note, error) {
	// read notes.json from local file system based on the config file
	// if the file does not exist, return an empty slice
	// if the file exists, parse the json and return the notes
	var noteFilePath string = viper.GetString("data_dir") + "/notes.json"

	notes := make([]notes.Note, 0)
	if _, err := os.Stat(noteFilePath); os.IsNotExist(err) {
		// create the file:
		SaveNotes(notes)
		return notes, nil
	}

	bytes, err := os.ReadFile(noteFilePath)
	if err != nil {
		return nil, errors.New("failed to read notes file: " + err.Error())
	}

	err = json.Unmarshal(bytes, &notes)
	if err != nil {
		return nil, errors.New("failed to unmarshal notes: " + err.Error())
	}

	return notes, nil
}

func SaveNotes(notes []notes.Note) error {
	var noteFilePath string = viper.GetString("data_dir") + "/notes.json"
	bytes, err := json.Marshal(notes)
	if err != nil {
		return errors.New("failed to marshal notes: " + err.Error())
	}
	err = os.WriteFile(noteFilePath, bytes, 0644)
	if err != nil {
		return errors.New("failed to write notes to file: " + err.Error())
	}
	return nil
}
