package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FilterForm struct {
	container         fyne.CanvasObject
	categorySelect    *widget.Select
	subCategorySelect *widget.Select
	regionSelect      *widget.Select
	provinceSelect    *widget.Select
	categoryMap       map[string]Category
	regionMap         map[string]Region
}

func NewFilterForm() *FilterForm {
	// Declare selects
	f := &FilterForm{
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
	f.categorySelect.OnChanged = func(selected string) {
		activeCategory := f.categoryMap[selected]
		subCategoryOptions := make([]string, len(activeCategory.SubCategories))
		subCategoryMap := make(map[string]SubCategory, len(activeCategory.SubCategories))
		for i, sc := range activeCategory.SubCategories {
			subCategoryOptions[i] = sc.Name
			subCategoryMap[sc.Name] = sc
		}
		f.subCategorySelect.Options = subCategoryOptions
		f.subCategorySelect.ClearSelected()
		f.subCategorySelect.Enable()
		f.subCategorySelect.Refresh()
	}

	f.regionSelect.OnChanged = func(selected string) {
		activeRegion := f.regionMap[selected]
		provinceOptions := make([]string, len(activeRegion.Provinces))
		provinceMap := make(map[string]Province, len(activeRegion.Provinces))
		for i, p := range activeRegion.Provinces {
			provinceOptions[i] = p.Name
			provinceMap[p.Name] = p
		}
		f.provinceSelect.Options = provinceOptions
		f.provinceSelect.ClearSelected()
		f.provinceSelect.Enable()
		f.provinceSelect.Refresh()
	}

	f.container = container.NewVBox(f.categorySelect, f.subCategorySelect, f.regionSelect, f.provinceSelect)
	return f
}

func PopulateFilterForm(f *FilterForm, categories []Category) {
	categoryOptions := make([]string, len(categories))
	f.categoryMap = make(map[string]Category, len(categories))
	for i, c := range categories {
		categoryOptions[i] = c.Name
		f.categoryMap[c.Name] = c
	}
	f.categorySelect.Options = categoryOptions
	f.categorySelect.Enable()
	f.categorySelect.Refresh()
}
