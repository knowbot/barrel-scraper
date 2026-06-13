# Barrel Scraper - Core

Beverage industry directory scraper + browser. Scrapes category/subcategory data from beverfood.com, stores in SQLite, provides GUI to browse by region/category.

**Architecture:** 
- `cmd/main.go` — Fyne GUI entrypoint; async init of Barrel storage + category fetch
- `internal/model/` — data structures (Category, SubCategory, Region, Province) + type-generic helpers (BuildMap, BuildList)
- `internal/service/scraper.go` — Colly-based web scraper; filters out "brand" categories (marca/marchio/marche)
- `internal/storage/barrel.go` — SQLite DB layer; manages schema+seed, CRUD for categories/companies
- `internal/ui/filter.go` — Fyne Select-based cascading filters (category→subcategory, region→province)
- `internal/utils/` — text utilities (punctuation cleanup)

**Database:** SQLite with auto-pragma foreign_keys; schema read from `db/transactions/01-schema.sql`, seed from `db/transactions/02-seed.sql`. Tables: categories, subcategories, provinces, companies, meta (for update tracking).

**Key invariants:**
- Filter.OnChanged callback tracks FilterSelection state (pointers to chosen items)
- Scraper isBrandCategory() filters Italian brand keywords from subcategories
- Barrel.Populate() called after async fetch to enable filters
- DB schema applied on NewBarrel() before seeding

Related: `mem:tech_stack`, `mem:conventions`, `mem:suggested_commands`, `mem:task_completion`