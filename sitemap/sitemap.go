package sitemap

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/url"
	"strings"

	models "github.com/connectwithub/go-html-parser/Models"
	htmlParser "github.com/connectwithub/go-html-parser/html-parser"
)

type loc struct {
	Link string `xml:"loc"`
}

func filterCurrentDomainLinks(links []models.Link, urlInput string) ([]models.Link, error) {
	filteredLinks := []models.Link{}
	urlParsed, err := url.Parse(urlInput)
	if err != nil {
		log.Fatalf("Error parsing loc: %v", err)
	}
	baseURL := &url.URL{
		Scheme: urlParsed.Scheme,
		Host:   urlParsed.Host,
	}
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			filteredLinks = append(filteredLinks, models.Link{Href: baseURL.String() + link.Href, Text: link.Text})
		case strings.HasPrefix(link.Href, "https://"+urlParsed.Host) || strings.HasPrefix(link.Href, "http://"+urlParsed.Host):
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks, nil
}

func Bfs(urlStr string, depth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i < depth; i++ {
		q, nq = nq, make(map[string]struct{})
		for nextURL := range q {
			if _, ok := seen[nextURL]; ok {
				continue
			}
			seen[nextURL] = struct{}{}
			links, _ := filterCurrentDomainLinks(htmlParser.ParseHTMLLinks(false, nextURL), nextURL)
			for _, link := range links {
				nq[link.Href] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for key := range seen {
		ret = append(ret, key)
	}
	return ret
}

func ConvertXML(links []string) string {
	xmlLinks := make([]loc, 0, len(links))
	for _, link := range links {
		xmlLinks = append(xmlLinks, loc{Link: link})
	}
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")

	enc.Encode(struct {
		XMLName xml.Name `xml:"urlset"`
		Xmlns   string   `xml:"xmlns,attr"`
		Links   []loc    `xml:"url"`
	}{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", Links: xmlLinks})
	return xml.Header + buf.String()
}
