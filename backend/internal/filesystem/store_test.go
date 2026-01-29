package filesystem

import (
	"os"
	"testing"
	"time"
)

func TestStore_SaveAndGet(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "notes-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	store := NewStore(tmpDir)

	// Test Save
	id := "test-note"
	content := "# Test Note\n\nContent"
	if err := store.Save(id, content); err != nil {
		t.Errorf("Save failed: %v", err)
	}

	// Test Get
	// ID might be normalized to test-note.md
	note, err := store.Get(id + ".md")
	if err != nil {
		t.Errorf("Get failed with extension: %v", err)
	}
	
	// Test Get without extension (should also work now)
	noteNoExt, err := store.Get(id)
	if err != nil {
		t.Errorf("Get failed without extension: %v", err)
	}

	if note.Title != "Test Note" {
		t.Errorf("Expected title 'Test Note', got '%s'", note.Title)
	}
	if noteNoExt.Title != "Test Note" {
		t.Errorf("Expected title 'Test Note' (no ext), got '%s'", noteNoExt.Title)
	}
}

// TestParseNoteContent tests the independent parser logic
func TestParseNoteContent(t *testing.T) {
	raw := []byte("---\ntitle: Foo\n---\nBody content")
	note := ParseNoteContent("foo.md", raw, time.Now())
	
	if note.Title != "Foo" {
		t.Errorf("Expected Foo, got %s", note.Title)
	}
	if note.Content != "Body content" {
		t.Errorf("Expected Body content, got %s", note.Content)
	}
}
