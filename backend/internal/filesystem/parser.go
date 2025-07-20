package filesystem

import (
	"bufio"
	"bytes"
	"strings"
	"time"

	"marko-backend/internal/models"
)

// ParseNoteContent parses the raw file content into metadata and body.
// It implements a simple custom parser to avoid external dependencies.
func ParseNoteContent(id string, content []byte, fileModTime time.Time) models.Note {
	meta := models.NoteMetadata{}
	body := string(content)

	// Check for frontmatter
	if bytes.HasPrefix(content, []byte("---\n")) || bytes.HasPrefix(content, []byte("---\r\n")) {
		parts := strings.SplitN(string(content), "---", 3)
		if len(parts) == 3 {
			// part 0 is empty (before first ---)
			// part 1 is frontmatter
			// part 2 is body
			parseFrontmatter(parts[1], &meta)
			body = strings.TrimSpace(parts[2])
		}
	}

	// Fallback/Defaults
	if meta.Title == "" {
		// Use ID or filename as title if missing
		meta.Title = strings.TrimSuffix(id, ".md")
		meta.Title = strings.ReplaceAll(meta.Title, "-", " ")
		meta.Title = strings.Title(meta.Title)
	}

	created := fileModTime
	if meta.Created != "" {
		if t, err := time.Parse("2006-01-02", meta.Created); err == nil {
			created = t
		}
	}
	
	updated := fileModTime
	if meta.Updated != "" {
		if t, err := time.Parse("2006-01-02", meta.Updated); err == nil {
			updated = t
		}
	}

	return models.Note{
		ID:        id,
		Title:     meta.Title,
		Content:   body,
		CreatedAt: created,
		UpdatedAt: updated,
	}
}

func parseFrontmatter(raw string, meta *models.NoteMetadata) {
	scanner := bufio.NewScanner(strings.NewReader(raw))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			
			// Simple unstuffing
			val = strings.Trim(val, `"'`)

			switch key {
			case "title":
				meta.Title = val
			case "created":
				meta.Created = val
			case "updated":
				meta.Updated = val
			case "tags":
				// Very basic tag parsing [a, b]
				val = strings.Trim(val, "[]")
				tags := strings.Split(val, ",")
				for _, t := range tags {
					meta.Tags = append(meta.Tags, strings.TrimSpace(t))
				}
			}
		}
	}
}
