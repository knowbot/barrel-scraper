package main

import (
	"fmt"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var baseURL string = "https://www.beverfood.com/directory-aziende-beverage/"
var brandWords = map[string]bool{
	"marca": true, "marchio": true, "marche": true, "marchi": true,
}

type ScrapedData struct {
	Categories []Category
}

type Category struct {
	Name          string
	SubCategories []SubCategory
}

type SubCategory struct {
	Name string
	url  string
}

func cleanText(s string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsPunct(r) {
				return -1
			}
			return r
		},
		s,
	)
}

func isBrandCategory(s string) bool {
	s = strings.ToLower(s)
	words := strings.SplitSeq(s, " ")
	for w := range words {
		fmt.Println(w)
		_, found := brandWords[w]
		if found {
			return true
		}
	}
	return false

}

func fetchCategories(categories chan<- []Category) {
	var res []Category
	scraper := colly.NewCollector()
	extensions.RandomUserAgent(scraper)
	scraper.OnHTML("div.elenco-directory", func(e *colly.HTMLElement) {
		c := Category{Name: cleanText(e.ChildText("h2"))}
		e.ForEach("a[href]", func(_ int, a *colly.HTMLElement) {
			text := cleanText(a.Text)
			if isBrandCategory(text) {
				return
			}
			c.SubCategories = append(
				c.SubCategories,
				SubCategory{
					Name: text,
					url:  a.Attr("href"),
				},
			)
		})
		res = append(res, c)
	})
	scraper.Visit(baseURL)
	categories <- res
	close(categories)
}

func fetchRegions() {

}

func fetchCities() {

}

func fetchCompanies() {

}

func main() {
	barrelScraper := app.New()
	bsWindow := barrelScraper.NewWindow("Barrel Scraper")
	bsWindow.Resize(fyne.NewSize(800, 400))

	var data ScrapedData

	categories := make(chan []Category)
	go fetchCategories(categories)

	go func() {
		var ok bool
		data.Categories, ok = <-categories
		if !ok || len(data.Categories) == 0 {
			bsWindow.SetContent(widget.NewLabel("No categories found."))
			return
		}
		categoryOptions := make([]string, len(data.Categories))
		categoryMap := make(map[string]Category, len(data.Categories))
		for i, c := range data.Categories {
			categoryOptions[i] = c.Name
			categoryMap[c.Name] = c
		}

		categorySelect := widget.NewSelect(categoryOptions, nil)
		subCategorySelect := widget.NewSelect(nil, nil)
		// regionSelect := widget.NewSelect(nil, nil)
		// citySelect := widget.NewSelect(nil, nil)

		categorySelect.PlaceHolder = "Categoria"
		subCategorySelect.PlaceHolder = "Sottocategoria"

		categorySelect.OnChanged = func(selected string) {
			activeCategory := categoryMap[selected]
			subCategoryOptions := make([]string, len(activeCategory.SubCategories))
			subCategoryMap := make(map[string]SubCategory, len(activeCategory.SubCategories))
			for i, sc := range activeCategory.SubCategories {
				subCategoryOptions[i] = sc.Name
				subCategoryMap[sc.Name] = sc
			}
			subCategorySelect.Options = subCategoryOptions
			subCategorySelect.ClearSelected()
			subCategorySelect.Refresh()
		}
		bsWindow.SetContent(container.NewVBox(categorySelect, subCategorySelect))
	}()

	bsWindow.Show()
	barrelScraper.Run()

}
