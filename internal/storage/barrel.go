package storage

// Barrel - where data is stored

import (
	"barrel-scraper/internal/model"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var base_path string = "./db"
var transaction_path string = base_path + "/transactions/"

type Barrel struct {
	db    *sql.DB
	ready bool
}

func NewBarrel() (*Barrel, error) {
	db, err := sql.Open("sqlite", base_path+"/barrel.db?_pragma=foreign_keys(1)")
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	schema, err := os.ReadFile(transaction_path + "/01-schema.sql")
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
	query, err := os.ReadFile(transaction_path + "/02-seed.sql")
	if err != nil {
		return fmt.Errorf("read seed file: %w", err)
	}
	_, err = tx.Exec(string(query))
	if err != nil {
		return fmt.Errorf("execute seed transaction: %w", err)
	}
	return tx.Commit()
}

func (b *Barrel) SelectLastUpdated(key string) (*time.Time, error) {
	var time *time.Time
	err := b.db.QueryRow(`SELECT value FROM meta WHERE meta.key = ?`, key).Scan(time)
	// Here we just return null pointer if can't find last update
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return time, nil
}

func (b *Barrel) UpdateLastUpdated(key string) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO metadata(key, last_updated) VALUES(?, datetime('now)) ON CONFLICT(key) DO UPDATE SET value=EXCLUDED.value`, key)
	if err != nil {
		return err
	}
	return nil
}

func (b *Barrel) UpsertCategories(categories []model.Category) error {
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
	if err = b.UpdateLastUpdated("categories"); err != nil {
		return err
	}
	if err = b.UpdateLastUpdated("subcategories"); err != nil {
		return err
	}
	return tx.Commit()
}

func (b *Barrel) SelectAllCategories() ([]model.Category, error) {
	cRows, err := b.db.Query(`SELECT name FROM categories`)
	if err != nil {
		return nil, err
	}
	categories := make([]model.Category, 0)
	for cRows.Next() {
		var c model.Category
		cRows.Scan(&c.Name)
		scRows, err := b.db.Query(`SELECT name, url FROM subcategories`)
		if err != nil {
			return nil, err
		}
		for scRows.Next() {
			var sc model.SubCategory
			scRows.Scan(&sc.Name, &sc.URL)
			c.SubCategories = append(c.SubCategories, sc)
		}
		categories = append(categories, c)
	}
	return categories, nil
}
