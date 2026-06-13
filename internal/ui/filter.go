package ui

import (
	"barrel-scraper/internal/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Filter struct {
	Container fyne.CanvasObject

	categorySelect    *widget.Select
	subCategorySelect *widget.Select
	regionSelect      *widget.Select
	provinceSelect    *widget.Select

	categories []model.Category
	regions    []model.Region

	selected  FilterSelection
	OnChanged func(FilterSelection)
}

type FilterSelection struct {
	Category    *model.Category
	SubCategory *model.SubCategory
	Region      *model.Region
	Province    *model.Province
}

func NewFilter() *Filter {
	// Declare selects
	f := &Filter{
		categorySelect:    widget.NewSelect(nil, nil),
		subCategorySelect: widget.NewSelect(nil, nil),
		regionSelect:      widget.NewSelect(nil, nil),
		provinceSelect:    widget.NewSelect(nil, nil),
	}

	// Init with placeholder
	f.categorySelect.PlaceHolder = "Categoria"
	f.subCategorySelect.PlaceHolder = "Sottocategoria"
	f.regionSelect.PlaceHolder = "Regione"
	f.provinceSelect.PlaceHolder = "Provincia"

	f.categorySelect.Disable()
	f.subCategorySelect.Disable()
	f.regionSelect.Disable()
	f.provinceSelect.Disable()

	// Wire in OnChanged logic
	f.categorySelect.OnChanged = func(string) {
		i := f.categorySelect.SelectedIndex()
		if i < 0 {
			return
		}
		c := f.categories[i]
		f.selected.Category = &c
		f.selected.SubCategory = nil
		populateSelect(f.subCategorySelect, c.SubCategories)
		f.OnChanged(f.selected)
	}

	f.subCategorySelect.OnChanged = func(selected string) {
		i := f.subCategorySelect.SelectedIndex()
		if i < 0 {
			return
		}
		f.selected.SubCategory = &f.selected.Category.SubCategories[i]
		f.OnChanged(f.selected)
	}

	f.regionSelect.OnChanged = func(selected string) {
		i := f.regionSelect.SelectedIndex()
		if i < 0 {
			return
		}
		r := f.regions[i]
		f.selected.Region = &r
		f.selected.Province = nil
		populateSelect(f.provinceSelect, r.Provinces)
		f.OnChanged(f.selected)
	}

	f.provinceSelect.OnChanged = func(selected string) {
		i := f.provinceSelect.SelectedIndex()
		if i < 0 {
			return
		}
		f.selected.Province = &f.selected.Region.Provinces[i]
		f.OnChanged(f.selected)
	}

	catLabel := widget.NewLabel("Per categoria:")
	locLabel := widget.NewLabel("Per luogo:")

	f.Container = container.NewGridWithColumns(2, container.NewVBox(catLabel, f.categorySelect, f.subCategorySelect), container.NewVBox(locLabel, f.regionSelect, f.provinceSelect))
	return f
}

func (f *Filter) Selection() FilterSelection { return f.selected }

func (f *Filter) Populate(categories []model.Category, regions []model.Region) {
	f.categories = categories
	f.regions = regions
	populateSelect(f.categorySelect, categories)
	populateSelect(f.regionSelect, regions)
}

func populateSelect[T model.Named](sel *widget.Select, items []T) {
	sel.Options = model.BuildList(items)
	sel.ClearSelected()
	sel.Enable()
	sel.Refresh()
}
