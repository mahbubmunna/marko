package models

import "time"

// Note represents a markdown note
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content,omitempty"` // Content is omitted in list view
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NoteMetadata is the frontmatter/metadata of a note
type NoteMetadata struct {
	Title     string    `yaml:"title"`
	Tags      []string  `yaml:"tags,omitempty"`
	Created   string    `yaml:"created,omitempty"`
	Updated   string    `yaml:"updated,omitempty"`
}
