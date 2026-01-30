package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"marko-backend/internal/filesystem"
	"marko-backend/internal/models"
	"marko-backend/internal/search"
)

type NoteHandler struct {
	Store         *filesystem.Store
	SearchService *search.Service
}

func NewNoteHandler(store *filesystem.Store, search *search.Service) *NoteHandler {
	return &NoteHandler{Store: store, SearchService: search}
}

func (h *NoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Simple router for now
	path := strings.TrimPrefix(r.URL.Path, "/api/notes")

	// Search endpoint: /api/search?q=... (mapped in main, but let's handle here if logical or separate)
	// Actually strict REST for /api/notes doesn't include search usually.
	// The user requested GET /api/search. We should handle that in a separate handler method or main router.
	// For simplicity, let's add a separate SearchHandler method and register it in main.
	// But sticking to NoteHandler for now.

	switch r.Method {
	case http.MethodGet:
		if path == "" || path == "/" {
			h.ListNotes(w, r)
		} else {
			id := strings.TrimPrefix(path, "/")
			h.GetNote(w, r, id)
		}
	case http.MethodPost:
		if path == "" || path == "/" {
			h.CreateNote(w, r)
		}
	case http.MethodPut:
		id := strings.TrimPrefix(path, "/")
		h.UpdateNote(w, r, id)
	case http.MethodDelete:
		id := strings.TrimPrefix(path, "/")
		h.DeleteNote(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *NoteHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.Store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request, id string) {
	note, err := h.Store.Get(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req models.Note
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If ID is not provided, try to derive from content or title
	if req.ID == "" {
		// First try to parse content to find title
		parsed := filesystem.ParseNoteContent("temp", []byte(req.Content), time.Now())
		if parsed.Title != "" && parsed.Title != "Temp" {
			req.Title = parsed.Title
		}

		if req.Title != "" {
			req.ID = slugify(req.Title)
		} else {
			req.ID = fmt.Sprintf("note-%d", time.Now().Unix())
		}
	}

	if err := h.Store.Save(req.ID, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Index asynchronously or synchronously
	if h.SearchService != nil {
		// Need full note for indexing (title etc). Parse again or read back.
		// Since we didn't parse full metadata from req.Content except for ID generation,
		// let's parse it properly or Read back.
		// Reading back is safer to stay in sync with what's on disk.
		go func() {
			if savedNote, err := h.Store.Get(req.ID); err == nil {
				h.SearchService.Index(savedNote)
			}
		}()
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": req.ID})
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request, id string) {
	var req models.Note
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare note object for indexing (need ID and Title if possible)
	// Currently req might only have Content. We need to parse full note to get Title for index.
	// Store.Save writes raw bytes.
	// Let's rely on Store to parse? Store.Save takes Content string.
	// We can parse it here to get metadata for Indexing.
	// OR we can read back from Store. reading back is safer.

	if err := h.Store.Save(id, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.SearchService != nil {
		go func() {
			// Read back to get full parsed note
			if savedNote, err := h.Store.Get(id); err == nil {
				h.SearchService.Index(savedNote)
			}
		}()
	}

	w.WriteHeader(http.StatusOK)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.SearchService != nil {
		go h.SearchService.Delete(id)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *NoteHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query required", http.StatusBadRequest)
		return
	}

	// Check if search service is available
	if h.SearchService == nil {
		http.Error(w, "Search service invalid/unavailable (check -tags fts5)", http.StatusServiceUnavailable)
		return
	}

	results, err := h.SearchService.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove non-alphanumeric (simplified)
	// For production, regex or loop is better
	return s
}
