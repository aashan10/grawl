package crawler

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

type Sitemap struct {
	Loc string `xml:"loc"`
}

type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

func FetchSitemaps(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch sitemap: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var sitemapIndex SitemapIndex
	err = xml.Unmarshal(body, &sitemapIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sitemap: %v", err)
	}

	var urls []string

	for _, index := range sitemapIndex.Sitemaps {
		urls = append(urls, index.Loc)
	}

	return urls, nil
}

func FetchSitemap(url string) (URLSet, error) {
	resp, err := http.Get(url)

	if err != nil {
		return URLSet{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return URLSet{}, fmt.Errorf("failed to fetch sitemap: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return URLSet{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var urlSet URLSet
	err = xml.Unmarshal(body, &urlSet)
	if err != nil {
		return URLSet{}, fmt.Errorf("failed to parse sitemap: %v", err)
	}
	return urlSet, nil
}
