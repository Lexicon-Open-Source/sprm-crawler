package scraper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sprm-crawler/crawler/models"
	"sprm-crawler/crawler/services"
	"sprm-crawler/repository"
	scraperModels "sprm-crawler/scraper/models"
	scraperService "sprm-crawler/scraper/services"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/samber/lo"
)

func StartScraping() error {
	ctx := context.Background()

	unscraped_urls, err := services.GetUnscrappedUrlFrontiers(ctx, 1000)
	if err != nil {
		return err
	}

	fmt.Printf("[DEBUG] Found %d unscraped URLs\n", len(unscraped_urls))

	if len(unscraped_urls) == 0 {
		fmt.Println("[INFO] No URLs to scrape")
		return nil
	}

	extractions, err := scrapeUrls(unscraped_urls)
	if err != nil {
		return fmt.Errorf("error scraping URLs: %w", err)
	}

	err = services.UpdateFrontierStatuses(ctx, lo.Map(unscraped_urls, func(urlFrontier repository.UrlFrontier, _ int) lo.Tuple2[string, int16] {
		return lo.Tuple2[string, int16]{A: urlFrontier.ID, B: models.URL_FRONTIER_STATUS_CRAWLED}
	}))

	if err != nil {
		return err
	}

	err = scraperService.UpsertExtraction(ctx, extractions)
	if err != nil {
		return err
	}

	return nil
}

func notContains(list []int, str int) bool {
	for _, v := range list {
		if v == str {
			return false
		}
	}
	return true
}

func scrapeUrls(urlFrontiers []repository.UrlFrontier) ([]repository.Extraction, error) {
	var extractions []repository.Extraction

	// Get the latest page number from extractions
	var pages []int
	for _, frontier := range urlFrontiers {
		if strings.Contains(frontier.Url, "page=") {
			parts := strings.Split(frontier.Url, "page=")
			if len(parts) > 1 {
				pageStr := strings.Split(parts[1], "&")[0]
				if page, err := strconv.Atoi(pageStr); err == nil && notContains(pages, page) {
					pages = append(pages, page)
				}
			}
		}
	}

	fmt.Printf("[INFO] Grouped %d URLs into %d unique pages to scrape\n", len(urlFrontiers), len(pages))
	fmt.Println("[INFO] Pages to scrape", pages)

	for _, page := range pages {
		var baseUrl string
		var ids []string

		for _, frontier := range urlFrontiers {
			if strings.Contains(frontier.Url, fmt.Sprintf("page=%d&", page)) {
				urls := strings.Split(frontier.Url, "#")
				baseUrl = urls[0]
				ids = append(ids, urls[1])
			}
		}

		fmt.Printf("[INFO] Scraping page: %s %d (contains %d items)\n", baseUrl, page, len(ids))

		c := colly.NewCollector()
		c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: 2,
			Delay:       2 * time.Second,
		})

		var siteContent string

		c.OnHTML("html", func(h *colly.HTMLElement) {
			siteContent = h.DOM.Text()
		})

		for idx, id := range ids {
			c.OnHTML(fmt.Sprintf("div.col-md-3.div-pesalah[data-key='%s']", id), func(e *colly.HTMLElement) {
				var metadata scraperModels.Metadata

				id := e.Attr("data-key")
				if id == "" {
					fmt.Println("Warning: Found item without data-key")
					return
				}

				fmt.Println("[INFO] Data Key", page, id, baseUrl, id)

				name := strings.TrimSpace(e.ChildText("div a"))
				if name == "" {
					fmt.Println("Warning: Found item without name, data-key:", id)
					return
				}

				tables := e.DOM.Find("table.table-custom")
				if tables.Length() >= 2 {
					e.DOM.Find("table.table-custom").Eq(0).Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
						key := strings.TrimSpace(s.Find("td").Eq(0).Text())
						value := strings.TrimSpace(s.Find("td").Eq(1).Text())

						switch key {
						case "Tertuduh":
							metadata.Accused = value
						case "No Pengenalan Diri":
							metadata.IDNumber = value
						case "Jantina":
							metadata.Gender = value
						case "Warganegara":
							metadata.Nationality = value
						case "Negeri":
							metadata.State = value
						case "Kategory":
							metadata.Category = value
						case "Majikan":
							metadata.Employer = value
						case "Jawatan":
							metadata.Position = value
						case "Mahkamah":
							metadata.Court = value
						case "Hakim":
							metadata.Judge = value
						case "Timbalan Pendakwa Raya / Pegawai Pendakwa":
							metadata.Officer = value
						case "Peguam Bela":
							metadata.DefenseAttorney = value
						case "Sabitan Lampau":
							metadata.PastConvictions = value
						case "Tarikh Jatuh Hukuman":
							metadata.SentencingDate = value
						case "Rayuan":
							metadata.Appeal = value
						}
					})

					e.DOM.Find("table.table-custom").Eq(1).Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
						key := strings.TrimSpace(s.Find("td").Eq(0).Text())
						value := strings.TrimSpace(s.Find("td").Eq(1).Text())

						switch key {
						case "Tertuduh":
							metadata.Accused = value
						case "No Pengenalan Diri":
							metadata.IDNumber = value
						case "Jantina":
							metadata.Gender = value
						case "Warganegara":
							metadata.Nationality = value
						case "Negeri":
							metadata.State = value
						case "Kategory":
							metadata.Category = value
						case "Majikan":
							metadata.Employer = value
						case "Jawatan":
							metadata.Position = value
						case "Mahkamah":
							metadata.Court = value
						case "Hakim":
							metadata.Judge = value
						case "Timbalan Pendakwa Raya / Pegawai Pendakwa":
							metadata.Officer = value
						case "Peguam Bela":
							metadata.DefenseAttorney = value
						case "Sabitan Lampau":
							metadata.PastConvictions = value
						case "Tarikh Jatuh Hukuman":
							metadata.SentencingDate = value
						case "Rayuan":
							metadata.Appeal = value
						}
					})
				}

				var details []scraperModels.ProcurementDetail
				e.DOM.Find("table.table-bordered").Eq(0).Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
					var detail scraperModels.ProcurementDetail
					s.Find("td").Each(func(i int, s *goquery.Selection) {
						value := strings.TrimSpace(s.Text())
						switch i {
						case 1:
							detail.Number = value
						case 2:
							detail.Summary = value
						case 3:
							detail.Offenses = value
						case 4:
							detail.Punishments = strings.Replace(value, "\n", " ", 100)
						}
					})
					details = append(details, detail)
				})

				metadata.ProcurementDetails = details

				extractionID := sha256.Sum256(fmt.Appendf(nil, "%s-%d-%s", id, page, time.Now().String()))
				extraction := repository.Extraction{
					ID:            hex.EncodeToString(extractionID[:]),
					UrlFrontierID: urlFrontiers[idx].ID,
					SiteContent:   &siteContent,
					RawPageLink:   &baseUrl,
					Metadata:      metadata,
					Language:      "ms",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}

				extractions = append(extractions, extraction)
			})
		}

		c.OnRequest(func(r *colly.Request) {
			fmt.Printf("[DEBUG] Requesting: %s\n", r.URL.String())
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Printf("[ERROR] Error scraping %s: %v\n", r.Request.URL.String(), err)
		})

		err := c.Visit(baseUrl)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("[DEBUG] Result", len(extractions))

	return extractions, nil
}
