package notes

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func (n *Note) AddTag(tag string) {
	if !slices.Contains(n.Tags, tag) {
		n.Tags = append(n.Tags, tag)
		n.UpdatedAt = time.Now().Format(time.DateTime)
	}
}

func (n *Note) RemoveTag(tag string) {
	newTags := make([]string, 0, len(n.Tags))
	for _, t := range n.Tags {
		if t != tag {
			newTags = append(newTags, t)
		}
	}
	n.Tags = newTags
	n.UpdatedAt = time.Now().Format(time.DateTime)
}

func (n *Note) GetTags() []string {
	return n.Tags
}

func (n *Note) CheckContentMatch(matchContent string) bool {
	return strings.Contains(strings.ToLower(n.Content), strings.ToLower(matchContent))
}

func (n *Note) UpdateContent(content string) {
	if content != "" {
		n.Content = content
		n.UpdatedAt = time.Now().Format(time.DateTime)
	}
}

func (n *Note) GetHeaders() []string {
	return []string{"ID", "Title", "Tags", "Content", "CreatedAt", "UpdatedAt"}
}

func (n *Note) PrepRow() []string {
	v := reflect.ValueOf(n).Elem()
	numFields := v.NumField()
	row := make([]string, numFields)
	// iterate over fields and add type specific logic to format field as string in row
	for i := 0; i < numFields; i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String { // Check if it's a slice of strings
				tags := make([]string, 0, field.Len())
				for j := 0; j < field.Len(); j++ {
					tags = append(tags, field.Index(j).String())
				}
				row[i] = strings.Join(tags, ", ")
			} else {
				row[i] = field.String()
			}
		case reflect.Int:
			row[i] = strconv.Itoa(int(field.Elem().Int()))
		default:
			row[i] = field.String()
		}
	}
	return row
}

// FilterNotesByContent filters a list of notes by content.
func FilterNotesByContent(notes []Note, content string) []Note {
	filtered := []Note{}
	for _, note := range notes {
		if note.CheckContentMatch(content) {
			filtered = append(filtered, note)
		}
	}
	return filtered
}
