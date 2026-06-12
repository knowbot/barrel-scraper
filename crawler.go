package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"

	_ "modernc.org/sqlite"
)

var baseURL string = "https://www.beverfood.com/directory-aziende-beverage/"

var brandWords = map[string]bool{
	"marca": true, "marchio": true, "marche": true, "marchi": true,
}

func isBrandCategory(s string) bool {
	s = strings.ToLower(s)
	words := strings.SplitSeq(s, " ")
	for w := range words {
		_, found := brandWords[w]
		if found {
			return true
		}
	}
	return false
}

func fetchCategories() ([]Category, error) {
	var fetchRes []Category
	var fetchErr error
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
					URL:  a.Attr("href"),
				},
			)
		})
		fetchRes = append(fetchRes, c)
	})
	scraper.OnError(func(r *colly.Response, e error) {
		fetchErr = fmt.Errorf("Request %s failed with status %d: %w", r.Request.URL, r.StatusCode, e)
	})
	scraper.Visit(baseURL)
	if fetchErr != nil {
		return nil, fetchErr
	}
	if len(fetchRes) == 0 {
		return nil, fmt.Errorf("No categories at %s (selector matched nothing)", baseURL)
	}
	return fetchRes, nil
}
