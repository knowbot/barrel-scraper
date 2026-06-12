package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func openDB(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		fmt.Println("Error opening database")
		return nil
	}
	db.Exec("PRAGMA foreign_keys = ON;")
	return db
}

func initDB(db *sql.DB) {
	schema, err := os.ReadFile("./sql/schema.sql")
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		panic(err)
	}
}

func saveCategories(db *sql.DB, categories []Category) error {
	tx, err := db.Begin()
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
			tx.Exec(`
			INSERT INTO subcategories(name, url, category_id)
				VALUES(?, ?, ?)
				ON CONFLICT(name)
				DO UPDATE SET name=EXCLUDED.name
			`, sc.Name, sc.URL, cID)
		}
	}
	return tx.Commit()
}
