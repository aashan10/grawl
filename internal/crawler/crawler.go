package crawler

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type PageResult struct {
	URL        string
	StatusCode int
	Error      error
	Response   string
}

type CrawlResult struct {
	mutex sync.Mutex
	pages []PageResult
}

func (cr *CrawlResult) AddPageResult(result PageResult) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	cr.pages = append(cr.pages, result)
}

func (cr *CrawlResult) GetResults() []PageResult {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	return cr.pages
}

type SitemapResult struct {
	results []PageResult
	mutex   sync.Mutex
	sitemap string
}

func (sr *SitemapResult) AddPageResult(result PageResult) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	sr.results = append(sr.results, result)
}

func (sr *SitemapResult) GetResults() []PageResult {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	return sr.results
}

func Crawl(urls []URL, concurrency int, skipCertCheck bool) []PageResult {

	var wg sync.WaitGroup

	jobs := make(chan URL)
	results := make(chan PageResult)

	for range concurrency {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for url := range jobs {
				client := GetClient(skipCertCheck)
				response, err := client.Get(url.Loc)

				result := PageResult{
					URL: url.Loc,
				}
				if err != nil {
					result.Error = err
				} else {
					result.StatusCode = response.StatusCode
					body, err := io.ReadAll(response.Body)
					if err != nil {
						result.Error = err
					} else {
						result.Response = string(body)
					}
					response.Body.Close()
				}

				results <- result

				if result.StatusCode != http.StatusOK {
					fmt.Printf("Error: %s returned status code %d\n", url.Loc, result.StatusCode)
				} else {
					fmt.Printf("Success: %s returned status code %d\n", url.Loc, result.StatusCode)
				}
			}

		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	var finalResult []PageResult

	for result := range results {
		finalResult = append(finalResult, result)
	}

	return finalResult
}
