package crawler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sprm-crawler/common"
	"sprm-crawler/crawler/models"
	"sprm-crawler/crawler/services"
	"sprm-crawler/repository"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func StartCrawlingUrl() {
	var (
		totalPage int
		err       error
		c         = colly.NewCollector()
	)

	c.OnHTML("li.last a", func(h *colly.HTMLElement) {
		totalPage, err = strconv.Atoi(h.Attr("data-page"))
		if err != nil {
			fmt.Println("Cannot parse last page invalid format", err)
		}
		totalPage += 1
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Started to get last page", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Successfully get last page", r.Request.URL)
		fmt.Println("Successfully get total page", totalPage)
		startCrawlingItems(totalPage)
	})

	c.Visit("https://www.sprm.gov.my/index.php?id=21&page_id=96&page=1&per-page=8")
}

func startCrawlingItems(totalPages int) {
	var frontiers []repository.UrlFrontier

	// Create a new collector for crawling items
	c := colly.NewCollector()

	// Setup rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 3,
		Delay:       1 * time.Second,
	})

	c.OnHTML("div.col-md-3.div-pesalah", func(e *colly.HTMLElement) {
		var metadata models.UrlFrontierMetadata

		id := e.Attr("data-key")
		if id == "" {
			fmt.Println("Warning: Found item without data-key")
			return
		}

		name := strings.TrimSpace(e.ChildText("div a"))
		if name == "" {
			fmt.Println("Warning: Found item without name, data-key:", id)
			return
		}

		number := ""
		divs := e.ChildTexts("div")
		if len(divs) >= 2 {
			number = strings.TrimSpace(divs[len(divs)-2])
		}

		tables := e.DOM.Find("table.table-custom")
		if tables.Length() >= 2 {
			e.DOM.Find("table.table-custom").Eq(1).Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
				key := strings.TrimSpace(s.Find("td").Eq(0).Text())
				value := strings.TrimSpace(s.Find("td").Eq(1).Text())
				if key == "Kategori" {
				}

				if key == "Tarikh Jatuh Hukuman" {
					metadata.StartDate = value
				}
			})
		}

		e.DOM.Find("table.table-bordered").Eq(0).Find("tbody > tr").Eq(0).Find("td").Each(func(i int, s *goquery.Selection) {
			value := strings.TrimSpace(s.Text())
			switch i {
			case 1:
				metadata.PackageName = value
			case 2:
				metadata.Scenario = value
			case 3:
			case 4:
				metadata.Status = value
			}
		})

		metadata.Title = name
		metadata.PackageNumber = number
		metadata.Status = "new"

		frontierID := sha256.Sum256([]byte(fmt.Sprintf("%s-%s", id, name)))
		baseURL := fmt.Sprintf("%s#%s", e.Request.URL.String(), id)
		frontier := repository.UrlFrontier{
			ID:        hex.EncodeToString(frontierID[:]),
			Domain:    common.CRAWLER_DOMAIN,
			Url:       baseURL,
			Crawler:   common.CRAWLER_NAME,
			Status:    int16(models.URL_FRONTIER_STATUS_NEW),
			Metadata:  metadata,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		frontiers = append(frontiers, frontier)
		fmt.Println("Scraping", baseURL)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Crawling page: %s\n", r.URL.String())
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Printf("Finished crawling page: %s\n", r.Request.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error crawling page %s: %v\n", r.Request.URL.String(), err)
	})

	fmt.Print("Starting to crawl all data", totalPages, "\n\n")

	// Crawl all pages
	for page := 1; page <= totalPages; page++ {
		url := fmt.Sprintf("https://www.sprm.gov.my/index.php?id=21&page_id=96&page=%d&per-page=8", page)
		err := c.Visit(url)
		if err != nil {
			fmt.Printf("Error visiting page %d: %v\n", page, err)
		}
	}

	c.Wait()

	// Save to database
	if len(frontiers) > 0 {
		fmt.Printf("Saving %d URL frontiers to database...\n", len(frontiers))
		ctx := context.Background()
		err := services.UpsertUrl(ctx, frontiers)
		if err != nil {
			fmt.Printf("Error saving frontiers to database: %v\n", err)
		} else {
			fmt.Printf("Successfully saved %d URL frontiers\n", len(frontiers))
		}
	} else {
		fmt.Println("No items found to save")
	}
}
