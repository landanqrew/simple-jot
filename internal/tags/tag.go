package tags

import (
	"strings"

	notes "github.com/landanqrew/simple-jot/internal/notes"
)

type TagMap struct {
	TagMap map[string]map[string]struct{} `json:"tag_map"`
}

func (t *TagMap) BuildTagMap(notes []notes.Note) {
	for _, note := range notes {
		for _, tag := range note.Tags {
			if _, ok := t.TagMap[tag]; !ok {
				t.TagMap[tag] = make(map[string]struct{})
			}
			t.TagMap[tag][note.ID] = struct{}{}
		}
	}
}

func (t *TagMap) AddTag(noteID, tag string) {
	if _, ok := t.TagMap[tag]; !ok {
		t.TagMap[tag] = make(map[string]struct{})
	}
	t.TagMap[tag][noteID] = struct{}{}
}

func (t *TagMap) RemoveTag(noteID, tag string) {
	if _, ok := t.TagMap[tag]; !ok {
		return
	}
	delete(t.TagMap[tag], noteID)
	if len(t.TagMap[tag]) == 0 {
		delete(t.TagMap, tag)
	}
}

func (t *TagMap) GetTagsForNote(noteID string) []string {
	tags := make([]string, 0, len(t.TagMap))
	for tag := range t.TagMap {
		if _, ok := t.TagMap[tag][noteID]; ok {
			tags = append(tags, tag)
		}
	}
	return tags
}

func (t *TagMap) GetNotesForTag(tag string) []string {
	notes := make([]string, 0, len(t.TagMap[tag]))
	for noteID := range t.TagMap[tag] {
		notes = append(notes, noteID)
	}
	return notes
}

func (t *TagMap) GetAllTags(prefix string) []string {
	tags := make([]string, 0, len(t.TagMap))
	for tag := range t.TagMap {
		if prefix != "" && strings.HasPrefix(tag, prefix) { // TODO: maybe base this on a trie for faster lookups?
			tags = append(tags, tag)
		} else if prefix == "" {
			tags = append(tags, tag)
		}
	}
	return tags
}