package service

// Scraper - what fetches the data

import (
	"barrel-scraper/internal/model"
	"barrel-scraper/internal/utils"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"

	_ "modernc.org/sqlite"
)

var baseURL string = "https://www.beverfood.com/directory-aziende-beverage/"

var brandWords = map[string]bool{
	"marca": true, "marchio": true, "marche": true, "marchi": true,
}

func createCollector(async bool) *colly.Collector {
	c := colly.NewCollector()
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*beverfood.*",
		Delay:       2000 * time.Millisecond,
		RandomDelay: 2000 * time.Millisecond,
	})
	c.CacheDir = "./.cache/colly"
	c.IgnoreRobotsTxt = true
	c.Async = async
	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("request %s failed with status %d: %v", r.Request.URL, r.StatusCode, e)
	})
	return c
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
	categories := make([]model.Category, 0)
	scraper := createCollector(false)
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
		categories = append(categories, c)
	})
	if err := scraper.Visit(baseURL); err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, fmt.Errorf("no categories at %s (selector matched nothing)", baseURL)
	}
	return categories, nil
}

func FetchCompanyURLs(durl string) ([]string, error) {
	company_urls := make([]string, 0)
	letter_index_urls, err := FetchLetterIndexURLs(durl)
	if err != nil {
		return nil, fmt.Errorf("letter index fetch failed: %w", err)
	}
	fmt.Println(len(letter_index_urls))
	if len(letter_index_urls) > 0 {
		for _, li_url := range letter_index_urls {
			res, _ := FetchDirectURLs(li_url)
			company_urls = append(company_urls, res...)
		}

	} else {
		company_urls, _ = FetchDirectURLs(durl)
	}
	company_urls = utils.RemoveDuplicates(company_urls)
	fmt.Println(len(company_urls))
	return company_urls, nil
}

func FetchLetterIndexURLs(directoryURL string) ([]string, error) {
	scraper := createCollector(false)
	letter_urls := make([]string, 0)
	scraper.OnHTML("div.lettere", func(e *colly.HTMLElement) {
		e.ForEach("a:not([class])", func(_ int, a *colly.HTMLElement) {
			if url := a.Attr("href"); url != "" {
				letter_urls = append(letter_urls, url)
			}
		})
	})
	scraper.Visit(directoryURL)
	return letter_urls, nil
}

func FetchDirectURLs(directoryURL string) ([]string, error) {
	// First read is sync
	// TODO: handle async error
	var mutex sync.Mutex
	var fetchErr error
	urls := make([]string, 0)
	pages := 1
	pageParam := "page"

	scanPages := func(e *colly.HTMLElement) {
		pageNums := e.DOM.Find("a.page-numbers")
		if pageNums.Length() > 0 {
			if query, exists := pageNums.Eq(-1).Attr("href"); exists {
				params, err := url.ParseQuery(strings.TrimPrefix(query, "?"))
				if err != nil {
					fetchErr = err
					return
				}
				for k, v := range params {
					pageParam = k
					pages, err = strconv.Atoi(v[0])
					if err != nil {
						fetchErr = err
						return
					}
					break
				}
			}
		}
	}
	collectURLs := func(e *colly.HTMLElement) {
		if url := e.ChildAttr("a", "href"); url != "" {
			mutex.Lock()
			urls = append(urls, url)
			mutex.Unlock()
		}
	}

	scraper := createCollector(false)
	scraper.OnResponse(func(r *colly.Response) {
		directoryURL = strings.Split(r.Request.URL.String(), "?")[0] // post-redirect, no query
	})
	scraper.OnHTML("nav.main-pagination", scanPages)
	scraper.OnHTML("h2.titolo-azz", collectURLs)
	if fetchErr = scraper.Visit(directoryURL); fetchErr != nil {
		return nil, fmt.Errorf("page 1 visit failed: %w", fetchErr)
	}

	if pages > 1 {
		scraper = createCollector(true)
		scraper.OnHTML("h2.titolo-azz", collectURLs)
		for i := 2; i <= pages; i++ {
			if fetchErr = scraper.Visit(directoryURL + fmt.Sprintf("?%s=%d", pageParam, i)); fetchErr != nil {
				return nil, fmt.Errorf("page %d visit failed: %w", i, fetchErr)
			}
		}
		scraper.Wait()
	}

	return urls, nil
}

// func scrapeCompany(url string) (*model.Company, error) {
// 	scraper := createCollector()
// }
