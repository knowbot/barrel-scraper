package main

import (
	"barrel-scraper/internal/model"
	"barrel-scraper/internal/service"
	"barrel-scraper/internal/ui"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Barrel Scraper")
	w.Resize(fyne.NewSize(800, 400))
	var distiller *service.Distiller
	var filter *ui.Filter
	go func() {
		var err error
		fyne.Do(func() {
			filter = ui.NewFilter()
			w.SetContent(filter.Container)
		})
		distiller, err = service.NewDistiller()
		if err != nil {
			log.Fatal(err)
		}
		defer distiller.Close()
		categories, err := distiller.GetCategories()
		if err != nil {
			fmt.Print(err)
		}

		fyne.Do(func() {
			filter.Populate(categories, []model.Region{})
		})
	}()

	w.ShowAndRun()
}
