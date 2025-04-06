package cmd

import (
	"fmt"
	"github.com/aashan10/grawl/internal/crawler"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

type SitemapResult struct {
	URLs  []crawler.URL
	mutex sync.Mutex
}

var concurrency int
var skipCertcheck bool

func sitemapWorker(sitemap string, sitemapResult *SitemapResult, wg *sync.WaitGroup, skipCertcheck bool) {
	defer wg.Done()
	urls, err := crawler.FetchSitemap(sitemap, skipCertcheck)

	if err != nil {
		fmt.Println("Error fetching sitemap:", err)
		return
	}

	sitemapResult.mutex.Lock()
	sitemapResult.URLs = append(sitemapResult.URLs, urls.URLs...)
	sitemapResult.mutex.Unlock()

}

var scanCmd = &cobra.Command{
	Use:   "scan <url>",
	Short: "Scan a sitemap for broken links",
	Long:  `Scan a sitemap for broken links`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			cmd.Help()
			os.Exit(1)
			return
		}

		entrypoint := args[0]

		sitemaps, err := crawler.FetchSitemaps(entrypoint, skipCertcheck)

		if err != nil {
			cmd.Println("Error fetching sitemaps:", err)
			os.Exit(1)
		}

		var wg sync.WaitGroup

		sitemapResult := SitemapResult{}

		for _, sitemap := range sitemaps {
			wg.Add(1)
			go sitemapWorker(sitemap, &sitemapResult, &wg, skipCertcheck)
		}
		wg.Wait()

		cmd.Printf("Fetched %v URLs from %v sitemaps\n", len(sitemapResult.URLs), len(sitemaps))

		if len(sitemapResult.URLs) == 0 {
			cmd.Println("No URLs found in the sitemaps.")
			os.Exit(1)
		}
		result := crawler.Crawl(sitemapResult.URLs, concurrency, skipCertcheck)

		for _, page := range result {
			if page.Error != nil {
				cmd.Printf("Error crawling %s: %v\n", page.URL, page.Error)
			} else {
				cmd.Printf("Crawled %s: %d\n", page.URL, page.StatusCode)
			}
		}

	},
}

func init() {

	scanCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 12, "Number of concurrent requests")
	scanCmd.Flags().BoolVarP(&skipCertcheck, "skip-cert-check", "k", false, "Skip SSL certificate verification")
	rootCmd.AddCommand(scanCmd)
}
