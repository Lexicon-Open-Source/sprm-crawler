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
	"time"

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
		fmt.Println("Successfully get total page", totalPage)
		ctx := context.Background()
		var frontiers []repository.UrlFrontier
		for i := 0; i < totalPage; i++ {
			url := fmt.Sprintf("https://www.sprm.gov.my/index.php?id=21&page_id=96&page=%d&per-page=8", i+1)
			frontierID := sha256.Sum256([]byte(url))
			baseURL := url
			frontier := repository.UrlFrontier{
				ID:        hex.EncodeToString(frontierID[:]),
				Domain:    common.CRAWLER_DOMAIN,
				Url:       baseURL,
				Crawler:   common.CRAWLER_NAME,
				Status:    int16(models.URL_FRONTIER_STATUS_NEW),
				Metadata:  nil,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			frontiers = append(frontiers, frontier)
		}
		services.UpsertUrl(ctx, frontiers)
		fmt.Printf("Successfully saved %d url frontiers\n", totalPage)
	})

	c.Visit("https://www.sprm.gov.my/index.php?id=21&page_id=96&page=1&per-page=8")
}
