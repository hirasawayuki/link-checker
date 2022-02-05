package httprequest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hirasawayuki/link-checker/html"
)

type CheckResult struct {
	AnchorResults []string
	ImgResults    []string
}

func (cr *CheckResult) append(node html.Node, result string) {
	switch v := node.(type) {
	case *html.AnchorNode:
		cr.AnchorResults = append(cr.AnchorResults, result)
	case *html.ImgNode:
		cr.ImgResults = append(cr.ImgResults, result)
	default:
		fmt.Println(v)
	}
}

func CheckPage(pageURL string) (*CheckResult, error) {
	check := &CheckResult{}
	u, err := url.Parse(pageURL)
	if err != nil {
		fmt.Printf("[ERROR] Parse URL failed. Plese check page url. (url=%s)\n", pageURL)
		return nil, fmt.Errorf("[ERROR] Parse URL failed. Plese check page url. (url=%s)\n", pageURL)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[ERROR] Request failed. url=%s, HTTP Status=%d\n", u, resp.StatusCode)
	}

	ns, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Parse HTML failed. (err=%w)\n", err)
	}

	for _, n := range ns {
		url, err := n.URL()
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Unexpected URL. (err=%w)\n", err)
		}
		if url.Host == "" {
			url.Host = u.Host
		}
		if url.Scheme == "" {
			url.Scheme = u.Scheme
		}
		resp, err := http.Get(url.String())
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Request failed. url=%s, HTTP Status=%d\n", u, resp.StatusCode)
		}
		if resp.StatusCode == http.StatusOK {
			check.append(n, fmt.Sprintf("✓ HTTP Status: %d URL: %s Text(alt): %s", resp.StatusCode, url, n))
			continue
		}
		check.append(n, fmt.Sprintf("X HTTP Status: %d URL: %s Text(alt): %s", resp.StatusCode, url, n))
	}

	return check, nil
}
