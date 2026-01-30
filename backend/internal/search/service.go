package search

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"marko-backend/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	db *sql.DB
}

func NewService(dataDir string) (*Service, error) {
	dbPath := filepath.Join(dataDir, "index.db")

	// Ensure directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	s := &Service{db: db}
	if err := s.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

func (s *Service) initSchema() error {
	// Create FTS5 virtual table
	// We use contentless table if we didn't want to store data,
	// but we might want snippets, so standard FTS is fine.
	query := `
	CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts5(id, title, content);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *Service) Index(note models.Note) error {
	// Upsert: Delete then Insert (simplest for FTS)
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing if any
	_, err = tx.Exec("DELETE FROM notes_fts WHERE id = ?", note.ID)
	if err != nil {
		return err
	}

	// Insert new
	_, err = tx.Exec("INSERT INTO notes_fts (id, title, content) VALUES (?, ?, ?)", note.ID, note.Title, note.Content)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM notes_fts WHERE id = ?", id)
	return err
}

func (s *Service) Search(query string) ([]models.Note, error) {
	// Simple prefix search possibility, but standard match is fine
	// FTS5 syntax: MATCH 'query'
	// We'll wrap in wildcards for partial match convenience if user wants
	rows, err := s.db.Query(`
		SELECT id, title, snippet(notes_fts, 2, '<b>', '</b>', '...', 64) 
		FROM notes_fts 
		WHERE notes_fts MATCH ? 
		ORDER BY rank 
		LIMIT 20`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []models.Note{}
	for rows.Next() {
		var n models.Note
		var snippet string
		if err := rows.Scan(&n.ID, &n.Title, &snippet); err != nil {
			continue // Skip bad rows
		}
		// We smuggle the snippet into Content for display in search results
		n.Content = snippet
		results = append(results, n)
	}
	return results, nil
}

func (s *Service) Close() error {
	return s.db.Close()
}

func (s *Service) ReindexAll(notes []models.Note) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear all
	if _, err := tx.Exec("DELETE FROM notes_fts"); err != nil {
		return err
	}

	// Batch insert
	stmt, err := tx.Prepare("INSERT INTO notes_fts (id, title, content) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, n := range notes {
		if _, err := stmt.Exec(n.ID, n.Title, n.Content); err != nil {
			log.Printf("Failed to index note %s: %v", n.ID, err)
		}
	}

	return tx.Commit()
}
