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

func sitemapWorker(sitemap string, sitemapResult *SitemapResult, wg *sync.WaitGroup) {
	defer wg.Done()
	urls, err := crawler.FetchSitemap(sitemap)

	if err != nil {
		fmt.Println("Error fetching sitemap:", err)
		return
	}

	sitemapResult.mutex.Lock()
	sitemapResult.URLs = append(sitemapResult.URLs, urls.URLs...)
	sitemapResult.mutex.Unlock()

}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a sitemap for broken links",
	Long:  `Scan a sitemap for broken links`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			cmd.Help()
			os.Exit(1)
			return
		}

		entrypoint := args[0]

		sitemaps, err := crawler.FetchSitemaps(entrypoint)

		if err != nil {
			cmd.Println("Error fetching sitemaps:", err)
			os.Exit(1)
		}

		var wg sync.WaitGroup

		sitemapResult := SitemapResult{}

		for _, sitemap := range sitemaps {
			wg.Add(1)
			go sitemapWorker(sitemap, &sitemapResult, &wg)
		}
		wg.Wait()

		cmd.Printf("Fetched %v URLs from %v sitemaps\n", len(sitemapResult.URLs), len(sitemaps))

		if len(sitemapResult.URLs) == 0 {
			cmd.Println("No URLs found in the sitemaps.")
			os.Exit(1)
		}
		result := crawler.Crawl(sitemapResult.URLs, 12)

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

	rootCmd.AddCommand(scanCmd)
}
