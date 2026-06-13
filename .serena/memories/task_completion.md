# Task Completion Checklist

When a coding task is done:

1. **Lint & vet:** `go vet ./...`
2. **Tests:** `go test ./...` (if tests added; none currently)
3. **Build:** `go build -o barrel-scraper.exe ./cmd/main.go` — must succeed with no errors
4. **Manual test:** `go run ./cmd/main.go` — GUI must open without panics; verify feature in UI
5. **Git status:** `git status` — all relevant files staged/committed (no uncommitted changes)
6. **Mod clean:** `go mod tidy` — no extraneous dependencies

**For DB changes:**
- Verify schema.sql syntax with `go run ./cmd/main.go` (applies schema on startup)
- Seed.sql must run without errors in Barrel.seed() transaction
- Check foreign key constraints defined (PRAGMA foreign_keys=1)