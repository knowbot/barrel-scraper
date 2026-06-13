# Code Conventions

**Package structure:**
- Lowercase package names (model, service, storage, ui, utils)
- No public/internal naming suffix (internal/ dir enforces visibility)
- Domain-focused (storage = DB ops, service = scraping, ui = GUI)

**Types & generics:**
- `Named` interface for polymorphic BuildMap/BuildList (all domain types implement GetName())
- Type parameters [T Named] used for reusable collection helpers
- Pointers in FilterSelection struct to track selected model items

**Error handling:**
- fmt.Errorf with %w for wrapping; propagate up (no silent swallows)
- Deferred Rollback() before explicit Commit() in tx blocks
- Service methods return ([]T, error); storage methods return error

**GUI (Fyne):**
- Select widget OnChanged bound at init
- Disable/Enable/Refresh() on selects to control state
- CanvasObject interface for component composition

**Database:**
- ON CONFLICT(...) DO UPDATE for idempotent inserts
- Foreign keys enabled via PRAGMA; schema defines all constraints
- Meta table for tracking last updates (key=table name)

**File paths:**
- Relative to cwd: ./sql/schema.sql, sql/barrel.db, sql/seed.sql
- No hardcoded absolute paths