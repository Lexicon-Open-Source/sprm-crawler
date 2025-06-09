package main

import (
	"fmt"
	"os"
	"sprm-crawler/common"
	"sprm-crawler/crawler"
	"sprm-crawler/repository"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sprm-crawler",
	Short: "An SPRM rasuah crawler",
	Long:  "sprm-crawler is a command-line web crawler application that can scrape websites and store results in a PostgreSQL database",
}

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl a website",
	Long:  `Crawl a website and store the results in the database. You can specify a single URL or multiple URLs to crawl.`,
	Run: func(cmd *cobra.Command, args []string) {
		crawler.StartCrawlingUrl()
	},
}

func init() {
	crawlCmd.Flags().IntP("depth", "d", 1, "crawl depth")
	crawlCmd.Flags().IntP("concurrent", "c", 5, "number of concurrent workers")
	crawlCmd.Flags().DurationP("delay", "t", time.Second, "delay between requests")

	rootCmd.AddCommand(crawlCmd)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env: %v\n", err)
		os.Exit(1)
	}
}

func setupDatabase() {
	err := common.ConnectDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to set database: %v\n", err)
		os.Exit(1)
	}

	query := repository.New(common.Pool)
	err = common.SetQuery(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to set query")
		os.Exit(1)
	}
}

func main() {
	loadEnv()
	setupDatabase()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
