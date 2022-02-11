package html

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const defaultNodeCap = 300

func Parse(r io.Reader) ([]Node, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	ns := make([]Node, 0, defaultNodeCap)

	p := &parser{parseFncs: make(map[string]parseFnc)}
	p.registerFnc(atom.A.String(), parseAnchorNode)
	p.registerFnc(atom.Img.String(), parseImgNode)
	p.parse(&ns, node)
	return ns, nil
}

type parser struct {
	parseFncs map[string]parseFnc
}

type parseFnc func(node *html.Node) Node

func (p *parser) registerFnc(key string, fnc parseFnc) {
	p.parseFncs[key] = fnc
}

func (p *parser) parse(ns *[]Node, node *html.Node) {
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			if fnc, ok := p.parseFncs[c.DataAtom.String()]; ok {
				*ns = append(*ns, fnc(c))
			}
			p.parse(ns, c)
		}
	}
}

type Node interface {
	URL() (string, error)
}

type AnchorNode struct {
	Href string
	Text string
}

func (n *AnchorNode) URL() (string, error) {
	u, err := url.Parse(n.Href)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

type ImgNode struct {
	Src string
	Alt string
}

func (n *ImgNode) URL() (string, error) {
	u, err := url.Parse(n.Src)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
