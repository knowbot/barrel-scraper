package main

import (
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var baseURL string = "https://www.beverfood.com/directory-aziende-beverage/"

type category struct {
	Name          string
	SubCategories []subCategory
}

type subCategory struct {
	Name string
	URL  string
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

func getCategories() []category {
	var categories []category
	scraper := colly.NewCollector()
	extensions.RandomUserAgent(scraper)
	scraper.OnHTML("div.elenco-directory", func(e *colly.HTMLElement) {
		c := category{Name: cleanText(e.ChildText("h2"))}
		e.ForEach("a[href]", func(_ int, a *colly.HTMLElement) {
			c.SubCategories = append(
				c.SubCategories,
				subCategory{
					Name: cleanText(a.Text),
					URL:  a.Attr("href"),
				},
			)
		})
		categories = append(categories, c)
	})
	scraper.Visit(baseURL)
	return categories
}

func main() {
	barrelScraper := app.New()
	bsWindow := barrelScraper.NewWindow("Barrel Scraper")

	categories := getCategories()

	list := widget.NewList(
		func() int {
			return len(categories)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(categories[i].Name)
		})
	bsWindow.SetContent(list)
	bsWindow.Show()
	barrelScraper.Run()
}
