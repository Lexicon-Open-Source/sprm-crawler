package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sprm-crawler",
	Short: "An SPRM rasuah crawler",
	Long:  "sprm-crawler is a command-line web crawler application that can scrape websites and store results in a PostgreSQL database",
}

var crawlCmd = &cobra.Command{
	Use:   "crawl [url]",
	Short: "Crawl a website",
	Long:  `Crawl a website and store the results in the database. You can specify a single URL or multiple URLs to crawl.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Started crawl", args[0])
		// TODO: Crawl URL
	},
}

func init() {
	crawlCmd.Flags().IntP("depth", "d", 1, "crawl depth")
	crawlCmd.Flags().IntP("concurrent", "c", 5, "number of concurrent workers")
	crawlCmd.Flags().DurationP("delay", "t", time.Second, "delay between requests")

	rootCmd.AddCommand(crawlCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
