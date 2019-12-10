package main

import (
	"flag"
	"fmt"
	"github.com/tohyung85/gophercises/exercise-5-sitemap-builder/sitemapbuilder"
)

func main() {
	sitePtr := flag.String("site", "https://gophercises.com", "Site used to build the site map")
	domainSitesOnlyPtr := flag.Bool("internal-only", true, "Site map to contain internal sites only")
	searchDepthPtr := flag.Int("max-depth", 3, "Max search depth from domain site")
	flag.Parse()

	smBuilder, err := sitemapbuilder.NewBuilder(*sitePtr, *domainSitesOnlyPtr, *searchDepthPtr)
	if err != nil {
		panic(err)
	}
	xmlOutput, err := smBuilder.BuildSiteMap()
	if err != nil {
		fmt.Printf("Error generating xmloutput.\n %s", err)
	}

	fmt.Printf("Result:\n%s", xmlOutput)
}
