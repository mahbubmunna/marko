package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"marko-backend/internal/filesystem"
	"marko-backend/internal/handlers"
)

func main() {
	// Initialize Store
	// User said "Default directory: ./data/notes".
	// We check if ../data/notes exists (running from backend/) or use ./data/notes
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
