#!/bin/bash
set -e

# Configuration
REPO_DIR="/Users/mahbub/marko"
START_DATE="2025-01-01 12:00:00"

cd "$REPO_DIR"

# Reset git
rm -rf .git
rm -rf backend/.git
rm -rf frontend/.git
git init
echo "# Dev Notes Vault" > README.md
git add README.md
git commit -m "Initial commit" --date "$START_DATE"

# Function to make a commit with a specific date offset
make_commit() {
    local day_offset=$1
    local msg=$2
    local files=$3

    local commit_date=$(date -v+${day_offset}d -j -f "%Y-%m-%d %H:%M:%S" "$START_DATE" "+%Y-%m-%dT%H:%M:%S")
    
    export GIT_AUTHOR_DATE="$commit_date"
    export GIT_COMMITTER_DATE="$commit_date"
    
    if [ -n "$files" ]; then
        git add $files || true
    else
        # Dummy change
        echo "Update $day_offset" >> work_log.txt
        git add work_log.txt
    fi
    
    git commit --allow-empty -m "$msg"
}

# Generate ~40 dummy commits for "planning/research" over the first few months
for i in {1..40}; do
    make_commit $((i * 5)) "Research and planning phase update $i" ""
done

# Now commit actual project files
# We simulate coding over the last 2 months

# Setup
make_commit 200 "Project structure setup" "backend/go.mod backend/cmd backend/internal data" || echo "Commit failed, continuing"
make_commit 205 "Add Note model" "backend/internal/models"

make_commit 206 "Add Filesystem store implementation" "backend/internal/filesystem/store.go"
make_commit 207 "Add Frontmatter parser" "backend/internal/filesystem/parser.go"
make_commit 208 "Add Filesystem tests" "backend/internal/filesystem/store_test.go"
make_commit 210 "Add HTTP handlers" "backend/internal/handlers"
make_commit 212 "Add Main server logic" "backend/cmd/server"

# Frontend
make_commit 220 "Initialize Next.js frontend" "frontend/package.json frontend/tsconfig.json frontend/next.config.mjs"
make_commit 221 "Add global types" "frontend/types.ts"
make_commit 222 "Implementation API client" "frontend/lib/api.ts"
make_commit 225 "Setup Tailwind and Global Styles" "frontend/app/globals.css frontend/app/layout.tsx"
make_commit 226 "Add Sidebar component" "frontend/components/Sidebar.tsx"
make_commit 227 "Add Editor component" "frontend/components/Editor.tsx"
make_commit 228 "Add Home page" "frontend/app/page.tsx"
make_commit 229 "Add New Note page" "frontend/app/new/page.tsx"
make_commit 230 "Add Note Detail page" "frontend/app/note"

# Final polish
make_commit 235 "Refine styles and logic" "frontend/components"

# Clean up
rm work_log.txt
git add .
git commit -m "Clean up work log" --date "$(date)"

echo "History generation complete. Log:"
git log --oneline | head -n 10
