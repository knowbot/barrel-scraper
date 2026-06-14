package storage

// Barrel - where data is stored

import (
	"barrel-scraper/internal/model"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var base_path string = "./db"
var migrations_path string = base_path + "/migrations/"

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
	schema, err := os.ReadFile(migrations_path + "/01-schema.sql")
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
	err = b.Migrate()
	if err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}
	b.ready = true
	return b, nil
}

func (b *Barrel) Close() error {
	b.ready = false
	return b.db.Close()
}

func (b *Barrel) Migrate() error {
	files, err := os.ReadDir(migrations_path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			tx, err := b.db.Begin()
			defer tx.Rollback()
			if err != nil {
				return err
			}
			query, err := os.ReadFile(filepath.Join(migrations_path, file.Name()))
			if err != nil {
				return err
			}
			_, err = tx.Exec(string(query))
			if err != nil {
				return err
			}
			err = tx.Commit()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Barrel) SelectLastUpdated(key string) (*time.Time, error) {
	var ts string
	err := b.db.QueryRow(`SELECT last_updated FROM metadata WHERE metadata.key = ?`, key).Scan(&ts)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, ts)
	// Here we just return null pointer if can't find last update
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (b *Barrel) UpdateLastUpdated(key string) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`INSERT INTO metadata(key, last_updated) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET last_updated=EXCLUDED.last_updated`, key, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (b *Barrel) SelectAllCategories() ([]model.Category, error) {
	caRows, err := b.db.Query(`SELECT id, name FROM categories`)
	if err != nil {
		return nil, err
	}
	categories := make([]model.Category, 0)
	for caRows.Next() {
		var ca model.Category
		caRows.Scan(&ca.ID, &ca.Name)
		scRows, err := b.db.Query(`SELECT id, name, url FROM subcategories WHERE category_id = ?`, ca.ID)
		if err != nil {
			return nil, err
		}
		for scRows.Next() {
			var sc model.SubCategory
			scRows.Scan(&sc.ID, &sc.Name, &sc.URL)
			ca.SubCategories = append(ca.SubCategories, sc)
		}
		categories = append(categories, ca)
	}
	return categories, nil
}

func (b *Barrel) SelectAllRegions() ([]model.Region, error) {
	reRows, err := b.db.Query(`SELECT id, name FROM regions`)
	if err != nil {
		return nil, err
	}
	regions := make([]model.Region, 0)
	for reRows.Next() {
		var re model.Region
		reRows.Scan(&re.ID, &re.Name)
		prRows, err := b.db.Query(`SELECT id, name, code FROM provinces WHERE region_id = ?`, re.ID)
		if err != nil {
			return nil, err
		}
		for prRows.Next() {
			var pr model.Province
			prRows.Scan(&pr.ID, &pr.Name, &pr.Code)
			pr.Name = pr.Name + " - " + pr.Code
			re.Provinces = append(re.Provinces, pr)
		}
		regions = append(regions, re)
	}
	return regions, nil
}

func (b *Barrel) UpsertCategories(categories []model.Category) error {
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, ca := range categories {
		var caID int64
		err = tx.QueryRow(`
			INSERT INTO categories(name)
			VALUES(?)
			ON CONFLICT(name)
			DO UPDATE SET name=EXCLUDED.name
			RETURNING id
		`, ca.Name).Scan(&caID)
		if err != nil {
			return err
		}
		for _, sc := range ca.SubCategories {
			_, err = tx.Exec(`
			INSERT INTO subcategories(name, url, category_id)
				VALUES(?, ?, ?)
				ON CONFLICT(name)
				DO UPDATE SET name=EXCLUDED.name
			`, sc.Name, sc.URL, caID)
			if err != nil {
				return err
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	if err = b.UpdateLastUpdated("categories"); err != nil {
		return err
	}
	if err = b.UpdateLastUpdated("subcategories"); err != nil {
		return err
	}
	return nil
}

func (b *Barrel) SelectCompanies(ca *model.Category, sc *model.SubCategory, re *model.Region, pr *model.Province) ([]model.Company, error) {
	query := `
		SELECT id, name, street_addr, cap, city, phone, fax, website, pr.name, sc.name
		FROM companies co
		INNER JOIN subcategories sc ON sc.id = co.subcategory_id 
		INNER JOIN categories ca ON ca.id = sc.category_id 
		INNER JOIN provinces pr ON pr.id = co.province_id
		INNER JOIN regions re ON re.id = pr.region_id
		WHERE 1=1
	`
	args := make([]any, 0)
	if ca != nil {
		query += " AND ca.id = ?"
		args = append(args, ca.ID)
	}
	if sc != nil {
		query += " AND sc.id = ?"
		args = append(args, sc.ID)
	}
	if re != nil {
		query += " AND re.id = ?"
		args = append(args, re.ID)
	}
	if pr != nil {
		query += " AND pr.id = ?"
		args = append(args, pr.ID)
	}
	companies := make([]model.Company, 0)
	rows, err := b.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var co model.Company
		err = rows.Scan(&co.ID, &co.Name, &co.StreetAddress, &co.CAP, &co.City, &co.Phone, &co.Fax, &co.Website, &co.Province, &co.Sector)
		if err != nil {
			return nil, err
		}
		companies = append(companies, co)
	}
	return companies, nil
}
