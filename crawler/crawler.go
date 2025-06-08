package crawler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sprm-crawler/repository"
	"strconv"

	"github.com/gocolly/colly/v2"
)

func StartCrawlingUrl() {
	c := colly.NewCollector()

	c.OnHTML("li.last a", func(h *colly.HTMLElement) {
		page, err := strconv.Atoi(h.Attr("data-page"))
		if err != nil {
			fmt.Println("Cannot parse last page", err)
		}
		page += 1

		var frontiers []repository.UrlFrontier

		for i := range page {
			url := fmt.Sprintf("https://www.sprm.gov.my/index.php?id=21&page_id=96&page=%d&per-page=8", i+1)
			id := sha256.Sum256([]byte(url))
			frontiers = append(frontiers, repository.UrlFrontier{
				ID:  hex.EncodeToString(id[:]),
				Url: url,
			})
			fmt.Println(url)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Successfully scraping", r.Request.URL)
	})

	c.Visit("https://www.sprm.gov.my/index.php?id=21&page_id=96&page=1&per-page=8")
}
