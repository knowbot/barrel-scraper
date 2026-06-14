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
	needFetch, err := d.HasStaleData("categories")
	fmt.Println(needFetch)
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
		return categories, nil
	}
	categories, err := d.barrel.SelectAllCategories()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (d *Distiller) GetRegions() ([]model.Region, error) {
	regions, err := d.barrel.SelectAllRegions()
	if err != nil {
		return nil, err
	}
	return regions, nil
}

func (d *Distiller) Extract(category model.Category, subCategory model.SubCategory, region model.Region, province model.Province) error {
	// Run a query on the company table based on whether
	return nil
}

func (d *Distiller) HasStaleData(key string) (bool, error) {
	t, err := d.barrel.SelectLastUpdated(key)
	if err != nil {
		return true, err
	}
	// If not found or older than 7 days
	if t == nil || (time.Since(*t).Hours()/24) > 7 {
		return true, nil
	}
	return false, nil
}
