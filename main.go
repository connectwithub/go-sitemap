package main

import (
	"flag"
	"fmt"

	"github.com/connectwithub/go-sitemap/sitemap"
)

func main() {
	loc := flag.String("loc", "", "Website for generating sitemap for.")
	depth := flag.Int("depth", 3, "Depth to traverse to (Default is 3)")
	flag.Parse()
	links := sitemap.Bfs(*loc, *depth)
	xmlSitemap := sitemap.ConvertXML(links)
	fmt.Println(xmlSitemap)

}
