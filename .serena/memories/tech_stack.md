# Tech Stack

**Language:** Go 1.26.4

**Key dependencies:**
- `fyne.io/fyne/v2` — desktop GUI framework (cross-platform)
- `github.com/gocolly/colly` — web scraper with random user-agent + HTML parsing
- `modernc.org/sqlite` — pure-Go SQLite driver (no cgo)

**Database:** SQLite; schema in `db/transactions/01-schema.sql`, seed in `db/transactions/02-seed.sql`. DB file: `sql/barrel.db` (PRAGMA foreign_keys=1).

**Build:** Standard Go toolchain; no custom build script yet.

**Module:** `barrel-scraper` at C:\Users\knowbot\Documents\Projects\barrel-scraper