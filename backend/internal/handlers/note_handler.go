package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"marko-backend/internal/filesystem"
	"marko-backend/internal/models"
)

type NoteHandler struct {
	Store *filesystem.Store
}

func NewNoteHandler(store *filesystem.Store) *NoteHandler {
	return &NoteHandler{Store: store}
}

func (h *NoteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Simple router for now
	path := strings.TrimPrefix(r.URL.Path, "/api/notes")
	
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

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": req.ID})
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request, id string) {
	var req models.Note
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ID from URL takes precedence
	if err := h.Store.Save(id, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Remove non-alphanumeric (simplified)
	// For production, regex or loop is better
	return s
}
