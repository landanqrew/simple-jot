package notes

import (
	"slices"
	"strings"
	"time"
)

type Note struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Project   string `json:"project"`
	Tags      []string `json:"tags"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
	return strings.Contains(n.Content, matchContent)
}

func (n *Note) UpdateContent(content string) {
	if content != "" {
		n.Content = content
		n.UpdatedAt = time.Now().Format(time.DateTime)
	}
}