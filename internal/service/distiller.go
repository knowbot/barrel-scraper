package service

import (
	"barrel-scraper/internal/model"
	"barrel-scraper/internal/storage"
	"fmt"
	"time"
)

type Distiller struct {
	barrel *storage.Barrel
}

func NewDistiller() (*Distiller, error) {
	b, err := storage.NewBarrel()
	if err != nil {
		return nil, err
	}
	return &Distiller{barrel: b}, nil
}

func (d *Distiller) Close() error {
	return d.barrel.Close()
}

func (d *Distiller) GetCategories() ([]model.Category, error) {
	needFetch, err := d.HasRecentData("categories")
	if err != nil {
		return nil, err
	}
	if needFetch {
		categories, err := FetchCategories()
		if err != nil {
			return nil, err
		}
		err = d.barrel.UpsertCategories(categories)
		if err != nil {
			return nil, err
		}
		fmt.Println("No categories, had to fetch!")
		return categories, nil
	}
	categories, err := d.barrel.SelectAllCategories()
	if err != nil {
		return nil, err
	}
	fmt.Println("We gottem!")
	return categories, nil
}

func (d *Distiller) HasRecentData(key string) (bool, error) {
	t, err := d.barrel.SelectLastUpdated(key)
	if err != nil {
		return false, err
	}
	// If not found or older than 7 days
	if t == nil || (time.Since(*t).Hours()/24) > 7 {
		return false, nil
	}
	return true, nil
}
