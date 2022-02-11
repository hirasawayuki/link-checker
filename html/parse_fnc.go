package html

import (
	"strings"

	"golang.org/x/net/html"
)

func parseAnchorNode(node *html.Node) Node {
	var text string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text = c.Data
			break
		}
	}

	var href string
	for _, v := range node.Attr {
		if v.Key == "href" {
			href = strings.TrimSpace(v.Val)
			break
		}
	}

	return &AnchorNode{Text: text, Href: href}
}

func parseImgNode(node *html.Node) Node {
	var src, alt string
	for _, v := range node.Attr {
		if v.Key == "src" {
			src = strings.TrimSpace(v.Val)
		}
		if v.Key == "alt" {
			alt = v.Val
		}
		if src != "" && alt != "" {
			break
		}
	}

	return &ImgNode{Src: src, Alt: alt}
}
