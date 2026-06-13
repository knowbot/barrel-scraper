package storage

// Barrel - where data is stored

import (
	"barrel-scraper/internal/model"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type Barrel struct {
	db    *sql.DB
	ready bool
}

func NewBarrel() (*Barrel, error) {
	db, err := sql.Open("sqlite", "sql/barrel.db?_pragma=foreign_keys(1)")
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	schema, err := os.ReadFile("./sql/schema.sql")
	if err != nil {
		return nil, fmt.Errorf("read schema file: %w", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		return nil, fmt.Errorf("apply schema transaction: %w", err)
	}
	b := &Barrel{
		db:    db,
		ready: false,
	}
	err = b.seed()
	if err != nil {
		return nil, fmt.Errorf("seeding: %w", err)
	}
	b.ready = true
	return b, nil
}

func (b *Barrel) Close() error {
	b.ready = false
	return b.db.Close()
}

func (b *Barrel) seed() error {
	tx, err := b.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("starting seed transaction: %w", err)
	}
	query, err := os.ReadFile("sql/seed.sql")
	if err != nil {
		return fmt.Errorf("read seed file: %w", err)
	}
	_, err = tx.Exec(string(query))
	if err != nil {
		return fmt.Errorf("execute seed transaction: %w", err)
	}
	return tx.Commit()
}

func (b *Barrel) LastUpdated(tableName string) error {
	row := b.db.QueryRow(`SELECT value FROM meta WHERE meta.key = ?`, tableName)
	err := row.Scan()
	if err != nil {
		return err
	}
	return nil
}

func (b *Barrel) LogUpdate(tableName string) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO meta(key, value) VALUES(?, datetime('now)) ON CONFLICT(key) DO UPDATE SET value=EXCLUDED.value`, tableName)
	if err != nil {
		return err
	}
	return nil
}

func (b *Barrel) InsertCategories(categories []model.Category) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, c := range categories {
		var cID int64
		tx.QueryRow(`
			INSERT INTO categories(name)
			VALUES(?)
			ON CONFLICT(name)
			DO UPDATE SET name=EXCLUDED.name
		`, c.Name).Scan(&cID)

		for _, sc := range c.SubCategories {
			_, err = tx.Exec(`
			INSERT INTO subcategories(name, url, category_id)
				VALUES(?, ?, ?)
				ON CONFLICT(name)
				DO UPDATE SET name=EXCLUDED.name
			`, sc.Name, sc.URL, cID)
			if err != nil {
				return err
			}
		}
	}
	if err = b.LogUpdate("companies"); err != nil {
		return err
	}
	return tx.Commit()
}
