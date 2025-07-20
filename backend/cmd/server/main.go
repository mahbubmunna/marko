package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"marko-backend/internal/filesystem"
	"marko-backend/internal/handlers"
)

func main() {
	// Initialize Store
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	
	// Assuming running from root or backend/
	// We want data/notes to be relative to the PROJECT root if possible,
	// but strictly we run backend from backend/ folder usually.
	// Let's resolve data/notes relative to where we run.
	// User said "Default directory: ./data/notes".
	// If we run `go run cmd/server/main.go` from `backend/`, then `./data/notes` would be `backend/data/notes`.
	// But project root has `data/notes`.
	// So we might need to go up one level if we are in backend.
	
	// Heuristic: check if ../data/notes exists, else use ./data/notes
	dataDir := "./data/notes"
	if _, err := os.Stat("../data/notes"); err == nil {
		dataDir = "../data/notes"
	}
	
	store := filesystem.NewStore(dataDir)
	noteHandler := handlers.NewNoteHandler(store)

	mux := http.NewServeMux()
	mux.Handle("/api/notes", noteHandler)
	mux.Handle("/api/notes/", noteHandler)

	// Wrap with CORS
	handler := corsMiddleware(mux)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Data directory: %s\n", dataDir)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all for local dev
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
