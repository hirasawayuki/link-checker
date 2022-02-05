package html

import (
	"bytes"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var parseFns = make(map[string]func(node *html.Node) Node)

type Node interface {
	URL() (*url.URL, error)
	String() string
}

type AnchorNode struct {
	Href string
	Text string
}

func (n *AnchorNode) URL() (*url.URL, error) {
	return url.Parse(n.Href)
}

func (n *AnchorNode) String() string {
	if n.Text == "" {
		return "empty"
	}
	return strings.TrimSpace(n.Text)
}

type ImgNode struct {
	Src string
	Alt string
}

func (n *ImgNode) URL() (*url.URL, error) {
	return url.Parse(n.Src)
}

func (n *ImgNode) String() string {
	if n.Alt == "" {
		return "empty"
	}
	return strings.TrimSpace(n.Alt)
}
func Parse(r io.Reader) ([]Node, error) {
	parseFns[atom.A.String()] = parseAnchorNode
	parseFns[atom.Img.String()] = parseImgNode

	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	ns := make([]Node, 0, 100)
	parseNodes(&ns, node)

	return ns, nil
}

func parseNodes(ns *[]Node, node *html.Node) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if f, ok := parseFns[c.DataAtom.String()]; ok {
				*ns = append(*ns, f(c))
			}
			parseNodes(ns, c)
		}
	}
}

func parseAnchorNode(node *html.Node) Node {
	var buff bytes.Buffer
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			buff.WriteString(c.Data)
		}
	}

	var href string
	for _, v := range node.Attr {
		if v.Key == "href" {
			href = v.Val
			break
		}
	}

	return &AnchorNode{Text: buff.String(), Href: href}
}

func parseImgNode(node *html.Node) Node {
	var src, alt string
	for _, v := range node.Attr {
		if v.Key == "src" {
			src = v.Val
		}
		if v.Key == "alt" {
			alt = v.Val
		}
	}

	return &ImgNode{Src: src, Alt: alt}
}
