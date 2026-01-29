package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"marko-backend/internal/models"
)

type Store struct {
	Dir string
	mu  sync.RWMutex
}

func NewStore(dir string) *Store {
	return &Store{Dir: dir}
}

func (s *Store) List() ([]models.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}

	// Initialize as empty slice so it marshals to [] instead of null
	notes := []models.Note{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			info, err := f.Info()
			if err != nil {
				continue
			}
			
			// We optimize list by not reading full content of every file if possible,
			// but to get Title we might need to read the header.
			// For simplicity and correctness with the requirement "If frontmatter is missing, derive title",
			// we will read the file. Modern SSDs can handle this for reasonable note counts.
			// For optimization we could limit reading to the first 500 bytes.
			
			content, err := os.ReadFile(filepath.Join(s.Dir, f.Name()))
			if err != nil {
				continue
			}
			
			note := ParseNoteContent(f.Name(), content, info.ModTime())
			note.Content = "" // Don't return full content in list
			notes = append(notes, note)
		}
	}
	return notes, nil
}

func (s *Store) Get(id string) (models.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.Dir, id)
	
	// If file doesn't exist, try appending .md
	// This helps when ID in URL is "note-123" but file is "note-123.md"
	if _, err := os.Stat(path); os.IsNotExist(err) && !strings.HasSuffix(id, ".md") {
		path = filepath.Join(s.Dir, id+".md")
	}

	// Security check to prevent directory traversal
	if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(s.Dir)) {
		return models.Note{}, fmt.Errorf("invalid path")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return models.Note{}, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return models.Note{}, err
	}

	return ParseNoteContent(id, content, info.ModTime()), nil
}

func (s *Store) Save(id string, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure directory exists
	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return err
	}
	
	if id == "" {
		return fmt.Errorf("id required")
	}

	// Use filename as ID. Ensure it has .md
	if !strings.HasSuffix(id, ".md") {
		id = id + ".md"
	}
	
	path := filepath.Join(s.Dir, id)
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.Dir, id)
	return os.Remove(path)
}
