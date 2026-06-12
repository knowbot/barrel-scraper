package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Barrel Scraper")
	w.Resize(fyne.NewSize(800, 400))

	go func() {
		var ff *FilterForm
		fyne.Do(func() {
			ff = NewFilterForm()
			w.SetContent(ff.container)
		})
		categories, err := fetchCategories()
		fyne.Do(func() {
			if err != nil || len(categories) == 0 {
				w.SetContent(widget.NewLabel("No categories found."))
				return
			}
			PopulateFilterForm(ff, categories)
		})
	}()

	w.Show()
	a.Run()

}
