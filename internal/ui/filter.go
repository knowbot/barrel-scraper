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
	resetButton       *widget.Button
	exportButton      *widget.Button

	categories []model.Category
	regions    []model.Region

	Selected  FilterSelection
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
		exportButton:      widget.NewButton("Esporta", nil),
	}

	f.resetButton = widget.NewButton("Reset", func() {
		f.categorySelect.ClearSelected()
		f.subCategorySelect.ClearSelected()
		f.regionSelect.ClearSelected()
		f.provinceSelect.ClearSelected()
	})

	// Init with placeholder
	f.categorySelect.PlaceHolder = "Seleziona..."
	f.subCategorySelect.PlaceHolder = "Seleziona..."
	f.regionSelect.PlaceHolder = "Seleziona..."
	f.provinceSelect.PlaceHolder = "Seleziona..."

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
		ca := f.categories[i]
		f.Selected.Category = &ca
		f.Selected.SubCategory = nil
		populateSelect(f.subCategorySelect, ca.SubCategories)
		if f.OnChanged != nil {
			f.OnChanged(f.Selected)
		}

	}

	f.subCategorySelect.OnChanged = func(selected string) {
		i := f.subCategorySelect.SelectedIndex()
		if i < 0 {
			return
		}
		f.Selected.SubCategory = &f.Selected.Category.SubCategories[i]
		if f.OnChanged != nil {
			f.OnChanged(f.Selected)
		}
	}

	f.regionSelect.OnChanged = func(selected string) {
		i := f.regionSelect.SelectedIndex()
		if i < 0 {
			return
		}
		re := f.regions[i]
		f.Selected.Region = &re
		f.Selected.Province = nil
		populateSelect(f.provinceSelect, re.Provinces)
		if f.OnChanged != nil {
			f.OnChanged(f.Selected)
		}
	}

	f.provinceSelect.OnChanged = func(selected string) {
		i := f.provinceSelect.SelectedIndex()
		if i < 0 {
			return
		}
		f.Selected.Province = &f.Selected.Region.Provinces[i]
		if f.OnChanged != nil {
			f.OnChanged(f.Selected)
		}
	}

	industry := widget.NewForm(
		widget.NewFormItem("Categoria *", f.categorySelect),
		widget.NewFormItem("Sottocategoria *", f.subCategorySelect),
	)

	location := widget.NewForm(
		widget.NewFormItem("Regione", f.regionSelect),
		widget.NewFormItem("Provincia", f.provinceSelect),
	)

	forms := container.NewGridWrap(fyne.NewSize(420, 680),
		container.NewVBox(
			industry,
			location,
		),
	)

	f.Container = container.NewVBox(
		forms,
		container.NewHBox(f.resetButton, f.exportButton),
	)
	return f
}

func (f *Filter) Selection() FilterSelection { return f.Selected }

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
