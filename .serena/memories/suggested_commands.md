# Project Commands

**Development:**
- `go run ./cmd/main.go` — start app (opens Fyne window)
- `go build -o barrel-scraper.exe ./cmd/main.go` — build binary (Windows .exe)
- `go test ./...` — run tests (none exist yet)
- `go mod tidy` — clean up dependencies
- `go vet ./...` — lint

**Database:**
- Direct SQLite: `sqlite3 sql/barrel.db` (if sqlite3 CLI available; else use Go SQL methods)

**Git:**
- Standard: `git status`, `git add`, `git commit`, `git log`
- Modified tracked: cmd/main.go, db/transactions/01-schema.sql, .gitignore
- Deleted tracked: db/barrel.db, db/seed.go (file-based seed replaced by SQL)

**Notes:**
- No Makefile yet; use `go` commands directly
- DB file (sql/barrel.db) should be in gitignore (not tracked, generated at runtime)