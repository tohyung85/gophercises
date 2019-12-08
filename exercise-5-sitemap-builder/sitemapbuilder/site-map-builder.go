package sitemapbuilder

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tohyung85/gophercises/exercise-4-link-parser/linkparser"
)

type SiteMapBuilder struct {
	domain          string
	hostname        string
	excludeExternal bool
	maxDepth        int
}

func NewBuilder(domain string, excludeExt bool, maxDepth int) (*SiteMapBuilder, error) {
	domainUrl, err := url.Parse(domain)
	if err != nil {
		return nil, err
	}
	return &SiteMapBuilder{domain, domainUrl.Hostname(), excludeExt, maxDepth}, nil
}

func (smb *SiteMapBuilder) BuildSiteMap() (string, error) {
	domainLinks := make(map[string]bool)

	smb.populateMapWithSiteLinks(domainLinks, smb.domain, 1)

	fmt.Printf("DONE Populating, domain map has %d entries\n", len(domainLinks))

	//TODO: Generate XML
	result := ""
	for href, _ := range domainLinks {
		result += fmt.Sprintf("url: %s\n", href)
	}

	return result, nil
}

func (smb *SiteMapBuilder) populateMapWithSiteLinks(domainLinks map[string]bool, site string, currDepth int) {
	if currDepth > smb.maxDepth { // Avoid searching too deep, 3 levels max
		return
	}
	site = getFullSite(site, smb.domain)
	if processed, _ := domainLinks[site]; processed { // Site already checked
		return
	}
	siteUrl, err := url.Parse(site)
	if err != nil {
		fmt.Println("Site is not a url")
		return
	}
	domainLinks[site] = true                // set site has been processed
	if siteUrl.Hostname() != smb.hostname { // no need to get more links on other domains
		return
	}
	// fmt.Printf("Working on url: %s with hostname: %s\n", site, siteUrl.Hostname())
	siteLinks, err := getSiteLinks(site)
	if err != nil {
		return
	}

	for sUrl, _ := range siteLinks { // add new site links to domain map, set site processed status to false and recurse on new site
		sUrl = getFullSite(sUrl, smb.domain)
		urlObj, err := url.Parse(sUrl)
		if err != nil {
			fmt.Printf("url error: %s\n Error: %s\n", sUrl, err)
		}
		if smb.excludeExternal && urlObj.Hostname() != smb.hostname { // if wish to exclude external sites, check and move on. Else add them to list and process.
			continue
		}
		if _, inMap := domainLinks[sUrl]; !inMap {
			domainLinks[sUrl] = false
		}
		smb.populateMapWithSiteLinks(domainLinks, sUrl, currDepth+1)
	}
}

func getFullSite(site string, domain string) string {
	if site == "" {
		return site
	}
	if string(site[0]) == "/" {
		return strings.TrimSuffix(domain, "/") + site
	}

	return site
}

func getSiteLinks(site string) (map[string]string, error) {
	siteLinks := make(map[string]string)
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Issue getting domain")
		return siteLinks, err
	}
	defer resp.Body.Close()

	linkParser := linkparser.NewParser()
	links, err := linkParser.ParseHTML(resp.Body)
	if err != nil {
		return siteLinks, err
	}

	siteLinks = getUniquePaths(links)

	return siteLinks, err
}

func getUniquePaths(allLinks []linkparser.Link) map[string]string {
	siteLinks := make(map[string]string)
	for _, lk := range allLinks {
		if _, inMap := siteLinks[lk.Href]; !inMap {
			siteLinks[lk.Href] = lk.Text
		}
	}

	return siteLinks
}
