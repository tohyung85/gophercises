package linkparser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type htmlLinkParser struct{}

type Link struct {
	Href string
	Text string
}

func (lp *htmlLinkParser) ParseHTML(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := processNode(doc)

	return links, nil
}

func processNode(n *html.Node) []Link {
	links := make([]Link, 0)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, processNode(c)...)
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		linkVal, found := getHrefAttr(n.Attr)
		if found {
			link := Link{linkVal, strings.TrimSpace(getAllNodeText(n))}
			links = append(links, link)
		}
	}
	return links
}

func getAllNodeText(n *html.Node) string {
	result := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += getAllNodeText(c)
	}
	if n.Type == html.TextNode {
		result += n.Data
	}

	return result
}

func getHrefAttr(attrs []html.Attribute) (string, bool) {
	found := false
	foundHref := ""
	for _, attr := range attrs {
		if attr.Key == "href" {
			found = true
			foundHref = attr.Val
		}
	}

	return foundHref, found
}

func NewParser() *htmlLinkParser {
	return &htmlLinkParser{}
}
