package service

// Scraper - what fetches the data

import (
	"barrel-scraper/internal/model"
	"barrel-scraper/internal/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"

	_ "modernc.org/sqlite"
)

func createCollector() *colly.Collector {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "beverfood.com",
		Delay:       250 * time.Millisecond,
		RandomDelay: 100 * time.Millisecond,
		Parallelism: 10,
	})
	c.Async = true
	return c
}

var baseURL string = "https://www.beverfood.com/directory-aziende-beverage/"

var brandWords = map[string]bool{
	"marca": true, "marchio": true, "marche": true, "marchi": true,
}

// Legacy code: use when not relying on a sitemap

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
	scraper := createCollector()
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

// func FetchCompanyUrls(durl string) ([]model.Company, error) {
// 	scraper := createCollector()
// 	extensions.RandomUserAgent(scraper)
// 	company_urls := make([]string, 0)
// 	letter_index_urls, err := FetchLetterIndexURLs(durl)
// 	if err != nil {
// 		return nil, fmt.Errorf("letter index fetch failed: %w", err)
// 	}
// 	fmt.Println(len(letter_index_urls))
// 	if len(letter_index_urls) > 0 {
// 		for _, li_url := range letter_index_urls {
// 			res, err := FetchCompanyURLs(li_url)
// 			if err != nil {
// 				return nil, fmt.Errorf("direct fetch failed: %w", err)
// 			}
// 			company_urls = append(company_urls, res...)
// 		}

// 	} else {
// 		company_urls, err = FetchCompanyURLs(durl)
// 		if err != nil {
// 			return nil, fmt.Errorf("direct fetch failed: %w", err)
// 		}
// 	}
// 	fmt.Println(len(company_urls))
// 	return nil, nil
// }

// func FetchLetterIndexURLs(directoryURL string) ([]string, error) {
// 	scraper := createCollector()
// 	letter_urls := make([]string, 0)
// 	scraper.OnHTML("div.lettere", func(e *colly.HTMLElement) {
// 		e.ForEach("a:not([class])", func(_ int, a *colly.HTMLElement) {
// 			if url := a.Attr("href"); url != "" {
// 				letter_urls = append(letter_urls, url)
// 			}
// 		})
// 	})
// 	scraper.Visit(directoryURL)
// 	scraper.Wait()
// 	return letter_urls, nil
// }

// func FetchCompanyURLs(directoryURL string) ([]string, error) {
// 	scraper := createCollector()
// 	urls := make([]string, 0)
// 	pages := 1
// 	scraper.OnHTML("nav.main-pagination", func(e *colly.HTMLElement) {
// 		pageNums := e.DOM.Find("a.page-numbers")
// 		if pageNums.Length() > 0 {
// 			n, err := strconv.Atoi(pageNums.Eq(-1).Text())
// 			if err != nil {
// 				return
// 			}
// 			pages = n
// 		}
// 	})
// 	scraper.OnHTML("h2.titolo-azz", func(e *colly.HTMLElement) {
// 		if url := e.ChildAttr("a", "href"); url != "" {
// 			urls = append(urls, url)
// 		}
// 	})
// 	err := scraper.Visit(directoryURL)
// 	for i := 2; i <= pages; i++ {
// 		err = scraper.Visit(directoryURL + fmt.Sprintf("?page=%d", i))
// 	}
// 	scraper.Wait()
// 	return urls, err
// }
