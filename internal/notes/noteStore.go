package notes

import "errors"

type NoteStore struct {
	NoteMap map[string]Note `json:"note_map"`
}

func (n *NoteStore) BuildNoteMap(notes []Note) {
	for _, note := range notes {
		n.NoteMap[note.ID] = note
	}
}

func (n *NoteStore) AddNote(note Note) error {
	if _, ok := n.NoteMap[note.ID]; !ok {
		return errors.New("note with ID (" + note.ID + ") already exists")
	}
	n.NoteMap[note.ID] = note
	return nil
}

func (n *NoteStore) GetNoteByID(id string) (Note, error) {

	if note, ok := n.NoteMap[id]; ok {
		return note, nil
	}
	return Note{}, errors.New("note not found")
}