package service

// Scraper - what fetches the data

import (
	"barrel-scraper/internal/model"
	"barrel-scraper/internal/utils"
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

func FetchCategories() ([]model.Category, error) {
	var fetchRes []model.Category
	var fetchErr error
	scraper := colly.NewCollector()
	extensions.RandomUserAgent(scraper)
	scraper.OnHTML("div.elenco-directory", func(e *colly.HTMLElement) {
		c := model.Category{Name: utils.CleanText(e.ChildText("h2"))}
		e.ForEach("a[href]", func(_ int, a *colly.HTMLElement) {
			text := utils.CleanText(a.Text)
			if isBrandCategory(text) {
				return
			}
			c.SubCategories = append(
				c.SubCategories,
				model.SubCategory{
					Name: text,
					URL:  a.Attr("href"),
				},
			)
		})
		fetchRes = append(fetchRes, c)
	})
	scraper.OnError(func(r *colly.Response, e error) {
		fetchErr = fmt.Errorf("request %s failed with status %d: %w", r.Request.URL, r.StatusCode, e)
	})
	scraper.Visit(baseURL)
	if fetchErr != nil {
		return nil, fetchErr
	}
	if len(fetchRes) == 0 {
		return nil, fmt.Errorf("no categories at %s (selector matched nothing)", baseURL)
	}
	return fetchRes, nil
}
