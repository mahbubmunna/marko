package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"marko-backend/internal/filesystem"
	"marko-backend/internal/handlers"
	"marko-backend/internal/search"
)

func main() {
	seedPtr := flag.Int("seed", 0, "Number of dummy notes to generate")
	flag.Parse()

	// Initialize Store
	dataDir := "./data/notes"
	if _, err := os.Stat("../data/notes"); err == nil {
		dataDir = "../data/notes"
	}

	store := filesystem.NewStore(dataDir)

	// Initialize Search
	searchService, err := search.NewService(dataDir)
	if err != nil {
		log.Printf("Warning: Failed to initialize search service: %v", err)
	} else {
		defer searchService.Close()

		// Initial sync/reindex on startup (simple approach)
		// strictly speaking we should only do this if requested or if valid index check fails
		// but for V1 let's accept startup cost or just do it in background
		go func() {
			log.Println("Syncing search index...")
			notes, err := store.List()
			if err == nil {
				// We need full content to index, List() usually optimizes this out.
				// store.List() currently returns full content?
				// Checked store.go: "note.Content = "" // Don't return full content in list"
				// So we need to re-read all.
				for _, n := range notes {
					fullNote, err := store.Get(n.ID)
					if err == nil {
						searchService.Index(fullNote)
					}
				}
			}
			log.Println("Search index synced.")
		}()
	}

	// Handle Seeding
	if *seedPtr > 0 {
		fmt.Printf("Seeding %d note(s)...\n", *seedPtr)
		seedNotes(store, searchService, *seedPtr)
		return
	}

	noteHandler := handlers.NewNoteHandler(store, searchService)

	mux := http.NewServeMux()
	mux.Handle("/api/notes", noteHandler)
	mux.Handle("/api/notes/", noteHandler)

	// Explicit search route
	mux.HandleFunc("/api/search", noteHandler.Search)

	// Wrap with CORS
	handler := corsMiddleware(mux)

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Data directory: %s\n", dataDir)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func seedNotes(store *filesystem.Store, search *search.Service, count int) {
	fmt.Println("Clearing existing notes...")
	if existing, err := store.List(); err == nil {
		for _, n := range existing {
			_ = store.Delete(n.ID)
			if search != nil {
				_ = search.Delete(n.ID)
			}
		}
	}

	topics := []string{
		"Go Concurrency Patterns", "React Server Components", "Docker Optimization",
		"Postgres Indexing", "Microservices Architecture", "Kubernetes Deployment",
		"Redis Caching Strategies", "GraphQL Schema Design", "Typescript Generics",
		"Linux Kernel Tuning", "AWS Lambda Functions", "Next.js App Router",
		"Rust Memory Safety", "Distributed Systems 101", "Git Workflow",
	}

	categories := []string{"Tutorial", "Snippet", "Debug Log", "Meeting Notes", "Draft", "Reference"}

	authors := []string{"Alice", "Bob", "Charlie", "Dave", "Eve"}

	for i := 0; i < count; i++ {
		topic := topics[rand.Intn(len(topics))]
		category := categories[rand.Intn(len(categories))]
		author := authors[rand.Intn(len(authors))]

		title := fmt.Sprintf("%s - %s", topic, category)
		if i < len(topics) {
			// Ensure we have at least one clean title per topic
			title = topics[i]
		}

		id := fmt.Sprintf("seed-note-%d", i)

		var contentBuilder strings.Builder
		contentBuilder.WriteString(fmt.Sprintf("---\ntitle: %s\nauthor: %s\ntags: [%s, %s]\n---\n\n", title, author, "tech", category))
		contentBuilder.WriteString(fmt.Sprintf("# %s\n\n", title))
		contentBuilder.WriteString(fmt.Sprintf("## Overview\n\nThis is a note regarding **%s**. It is crucial for understanding current stack.\n\n", topic))

		// Simulate code block
		contentBuilder.WriteString("## Code Snippet\n\n```go\nfunc main() {\n\tfmt.Println(\"Hello World\")\n}\n```\n\n")

		contentBuilder.WriteString("## Key Takeaways\n\n")
		contentBuilder.WriteString("- Importance of clean code\n")
		contentBuilder.WriteString("- Performance matters\n")
		contentBuilder.WriteString("- Scalability is key\n\n")

		contentBuilder.WriteString(fmt.Sprintf("> Created by %s at %s\n", author, time.Now().Format(time.RFC3339)))

		content := contentBuilder.String()

		if err := store.Save(id, content); err != nil {
			log.Printf("Failed to save seed note %d: %v", i, err)
			continue
		}

		if search != nil {
			// Construct note model for indexing
			n := filesystem.ParseNoteContent(id+".md", []byte(content), time.Now())
			if err := search.Index(n); err != nil {
				log.Printf("Failed to index seed note %d: %v", i, err)
			}
		}
	}
	fmt.Println("Seeding complete with realistic data.")
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
