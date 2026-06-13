package main

import (
	"barrel-scraper/internal/model"
	"barrel-scraper/internal/service"
	"barrel-scraper/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Barrel Scraper")
	w.Resize(fyne.NewSize(800, 400))

	go func() {
		var ff *ui.Filter
		fyne.Do(func() {
			ff = ui.NewFilter()
			w.SetContent(ff.Container)
		})
		categories, err := service.FetchCategories()
		fyne.Do(func() {
			if err != nil || len(categories) == 0 {
				w.SetContent(widget.NewLabel("No categories found."))
				return
			}
			ff.Populate(categories, []model.Region{})
		})
	}()

	w.ShowAndRun()
}
